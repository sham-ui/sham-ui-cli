package member

import (
	"context"
	"errors"
	"{{ shortName }}/internal/model"
	"{{ shortName }}/pkg/asserts"
	"{{ shortName }}/pkg/postgres"
	"{{ shortName }}/pkg/postgres/testingdb"
	"testing"
)

func insertMembers(t *testing.T, db *postgres.Database) {
	t.Helper()
	ctx := context.Background()
	t.Cleanup(func() {
		_, err := db.ExecContext(ctx, "DELETE FROM members")
		asserts.NoError(t, err)
	})
	stmt, err := db.PrepareContext(
		ctx,
		"INSERT INTO members (id, email, name, password, is_superuser) VALUES ($1, $2, $3, $4, $5)",
	)
	asserts.NoError(t, err)

	for _, args := range [][]any{
		{"42", "test@example.com", "tester", "password", true},
		{"43", "test2@example.com", "tester2", "password2", false},
	} {
		_, err := stmt.Exec(args...)
		asserts.NoError(t, err)
	}
}

func TestRepository_Create(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertMembers(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name        string
		ctx         context.Context
		data        model.MemberWithPassword
		expectedErr error
	}{
		{
			name: "success",
			ctx:  context.Background(),
			data: model.MemberWithPassword{
				Member: model.Member{ //nolint:exhaustruct
					Email:       "test3@example.com",
					Name:        "tester3",
					IsSuperuser: true,
				},
				HashedPassword: model.MemberHashedPassword("password123"),
			},
			expectedErr: nil,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			data: model.MemberWithPassword{
				Member: model.Member{ //nolint:exhaustruct
					Email:       "test3@example.com",
					Name:        "tester3",
					IsSuperuser: true,
				},
				HashedPassword: model.MemberHashedPassword("password123"),
			},
			expectedErr: errors.New("query: context canceled"),
		},
		{
			name: "email already exists",
			ctx:  context.Background(),
			data: model.MemberWithPassword{
				Member: model.Member{ //nolint:exhaustruct
					Email:       "test@example.com",
					Name:        "tester",
					IsSuperuser: true,
				},
				HashedPassword: model.MemberHashedPassword("password"),
			},
			expectedErr: model.ErrMemberEmailAlreadyExists,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Act
			err := repo.Create(test.ctx, test.data)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestRepository_Update(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	repo := NewRepository(db)

	testCases := []struct {
		name            string
		ctx             context.Context
		data            model.Member
		expectedErr     error
		dataAfterUpdate model.MemberWithPassword
	}{
		{
			name: "success",
			ctx:  context.Background(),
			data: model.Member{
				ID:          model.MemberID("42"),
				Email:       "test3@example.com",
				Name:        "tester3",
				IsSuperuser: false,
			},
			expectedErr: nil,
			dataAfterUpdate: model.MemberWithPassword{
				Member: model.Member{
					ID:          model.MemberID("42"),
					Email:       "test3@example.com",
					Name:        "tester3",
					IsSuperuser: false,
				},
				HashedPassword: model.MemberHashedPassword("password"),
			},
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			data: model.Member{
				ID:          model.MemberID("42"),
				Email:       "test@example.com",
				Name:        "tester3",
				IsSuperuser: false,
			},
			expectedErr: errors.New("query: context canceled"),
			dataAfterUpdate: model.MemberWithPassword{
				Member: model.Member{
					ID:          model.MemberID("42"),
					Email:       "test@example.com",
					Name:        "tester",
					IsSuperuser: true,
				},
				HashedPassword: model.MemberHashedPassword("password"),
			},
		},
		{
			name: "email already exists",
			ctx:  context.Background(),
			data: model.Member{
				ID:          model.MemberID("42"),
				Email:       "test2@example.com",
				Name:        "tester2",
				IsSuperuser: false,
			},
			expectedErr: model.ErrMemberEmailAlreadyExists,
			dataAfterUpdate: model.MemberWithPassword{
				Member: model.Member{
					ID:          model.MemberID("42"),
					Email:       "test@example.com",
					Name:        "tester",
					IsSuperuser: true,
				},
				HashedPassword: model.MemberHashedPassword("password"),
			},
		},
		{
			name: "user not found",
			ctx:  context.Background(),
			data: model.Member{
				ID:          model.MemberID("1"),
				Email:       "test3@example.com",
				Name:        "tester3",
				IsSuperuser: false,
			},
			expectedErr: nil,
			dataAfterUpdate: model.MemberWithPassword{
				Member: model.Member{
					ID:          model.MemberID("42"),
					Email:       "test@example.com",
					Name:        "tester",
					IsSuperuser: true,
				},
				HashedPassword: model.MemberHashedPassword("password"),
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			insertMembers(t, db)

			// Act
			err := repo.Update(test.ctx, test.data)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)

			data, err := repo.GetByEmail(context.Background(), test.dataAfterUpdate.Email)
			asserts.NoError(t, err)
			asserts.Equals(t, test.dataAfterUpdate, *data, "data after update")
		})
	}
}

func TestRepository_GetByEmail(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertMembers(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		email        string
		expectedData *model.MemberWithPassword
		expectedErr  error
	}{
		{
			name:  "success",
			ctx:   context.Background(),
			email: "test@example.com",
			expectedData: &model.MemberWithPassword{
				Member: model.Member{
					ID:          "42",
					Email:       "test@example.com",
					Name:        "tester",
					IsSuperuser: true,
				},
				HashedPassword: model.MemberHashedPassword("password"),
			},
			expectedErr: nil,
		},
		{
			name:         "not found",
			ctx:          context.Background(),
			email:        "notfound@example.com",
			expectedData: nil,
			expectedErr:  model.ErrMemberNotFound,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			email:        "test@example.com",
			expectedData: nil,
			expectedErr:  errors.New("query: context canceled"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Act
			data, err := repo.GetByEmail(test.ctx, test.email)

			// Assert
			asserts.Equals(t, test.expectedData, data)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestRepository_Find(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertMembers(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		offset       int64
		limit        int64
		expectedData []model.Member
		expectedErr  error
	}{
		{
			name:   "success",
			ctx:    context.Background(),
			offset: 0,
			limit:  2,
			expectedData: []model.Member{
				{
					ID:          "42",
					Email:       "test@example.com",
					Name:        "tester",
					IsSuperuser: true,
				},
				{
					ID:          "43",
					Email:       "test2@example.com",
					Name:        "tester2",
					IsSuperuser: false,
				},
			},
			expectedErr: nil,
		},
		{
			name:   "offset = 1",
			ctx:    context.Background(),
			offset: 1,
			limit:  2,
			expectedData: []model.Member{
				{
					ID:          "43",
					Email:       "test2@example.com",
					Name:        "tester2",
					IsSuperuser: false,
				},
			},
			expectedErr: nil,
		},
		{
			name:   "limit = 1",
			ctx:    context.Background(),
			offset: 0,
			limit:  1,
			expectedData: []model.Member{
				{
					ID:          "42",
					Email:       "test@example.com",
					Name:        "tester",
					IsSuperuser: true,
				},
			},
			expectedErr: nil,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			offset:       0,
			limit:        2,
			expectedData: nil,
			expectedErr:  errors.New("query: context canceled"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Act
			data, err := repo.Find(test.ctx, test.offset, test.limit)

			// Assert
			asserts.Equals(t, test.expectedData, data)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestRepository_Total(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertMembers(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		expectedData int
		expectedErr  error
	}{
		{
			name:         "success",
			ctx:          context.Background(),
			expectedData: 2,
			expectedErr:  nil,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			expectedData: 0,
			expectedErr:  errors.New("query: context canceled"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Act
			data, err := repo.Total(test.ctx)

			// Assert
			asserts.Equals(t, test.expectedData, data)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestRepository_UpdateEmail(t *testing.T) {
	db := testingdb.Connect(t)
	repo := NewRepository(db)

	testCases := []struct {
		name          string
		ctx           context.Context
		id            model.MemberID
		email         string
		expectedErr   error
		expectedEmail string
	}{
		{
			name:          "success",
			ctx:           context.Background(),
			id:            "42",
			email:         "changed@example.com",
			expectedEmail: "changed@example.com",
		},
		{
			name:          "not found",
			ctx:           context.Background(),
			id:            "24",
			email:         "changed@example.com",
			expectedEmail: "test@example.com",
		},
		{
			name:          "not unique",
			ctx:           context.Background(),
			id:            "42",
			email:         "test2@example.com",
			expectedErr:   model.ErrMemberEmailAlreadyExists,
			expectedEmail: "test@example.com",
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			id:            "42",
			email:         "changed@example.com",
			expectedErr:   errors.New("query: context canceled"),
			expectedEmail: "test@example.com",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			insertMembers(t, db)

			// Action
			err := repo.UpdateEmail(test.ctx, test.id, test.email)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)

			var email string
			err = db.QueryRowContext(context.Background(), "SELECT email FROM members WHERE id = '42'").Scan(&email)
			asserts.NoError(t, err)
			asserts.Equals(t, test.expectedEmail, email)
		})
	}
}

func TestRepository_UpdateName(t *testing.T) {
	db := testingdb.Connect(t)
	repo := NewRepository(db)

	testCases := []struct {
		caseName     string
		ctx          context.Context
		id           model.MemberID
		name         string
		expectedErr  error
		expectedName string
	}{
		{
			caseName:     "success",
			ctx:          context.Background(),
			id:           "42",
			name:         "changed-tester",
			expectedName: "changed-tester",
		},
		{
			caseName:     "not found",
			ctx:          context.Background(),
			id:           "24",
			name:         "changed-tester",
			expectedName: "tester",
		},
		{
			caseName: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			id:           "42",
			name:         "changed-tester",
			expectedErr:  errors.New("query: context canceled"),
			expectedName: "tester",
		},
	}

	for _, test := range testCases {
		t.Run(test.caseName, func(t *testing.T) {
			// Arrange
			insertMembers(t, db)

			// Action
			err := repo.UpdateName(test.ctx, test.id, test.name)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)

			var name string
			err = db.QueryRowContext(context.Background(), "SELECT name FROM members WHERE id = '42'").Scan(&name)
			asserts.NoError(t, err)
			asserts.Equals(t, test.expectedName, name)
		})
	}
}

func TestRepository_UpdatePassword(t *testing.T) {
	db := testingdb.Connect(t)
	repo := NewRepository(db)

	testCases := []struct {
		caseName         string
		ctx              context.Context
		id               model.MemberID
		password         model.MemberHashedPassword
		expectedErr      error
		expectedPassword string
	}{
		{
			caseName:         "success",
			ctx:              context.Background(),
			id:               "42",
			password:         model.MemberHashedPassword("new-password"),
			expectedPassword: "new-password",
		},
		{
			caseName:         "not found",
			ctx:              context.Background(),
			id:               "24",
			password:         model.MemberHashedPassword("new-password"),
			expectedPassword: "password",
		},
		{
			caseName: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			id:               "42",
			password:         model.MemberHashedPassword("new-password"),
			expectedErr:      errors.New("query: context canceled"),
			expectedPassword: "password",
		},
	}

	for _, test := range testCases {
		t.Run(test.caseName, func(t *testing.T) {
			// Arrange
			insertMembers(t, db)

			// Action
			err := repo.UpdatePassword(test.ctx, test.id, test.password)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)

			var pass string
			err = db.QueryRowContext(context.Background(), "SELECT password FROM members WHERE id = '42'").Scan(&pass)
			asserts.NoError(t, err)
			asserts.Equals(t, test.expectedPassword, pass)
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	db := testingdb.Connect(t)
	repo := NewRepository(db)

	testCases := []struct {
		caseName          string
		ctx               context.Context
		id                model.MemberID
		expectedErr       error
		expectedCountInDB int
	}{
		{
			caseName:          "success",
			ctx:               context.Background(),
			id:                "42",
			expectedCountInDB: 1,
		},
		{
			caseName:          "not found",
			ctx:               context.Background(),
			id:                "24",
			expectedCountInDB: 2,
		},
		{
			caseName: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			id:                "42",
			expectedErr:       errors.New("query: context canceled"),
			expectedCountInDB: 2,
		},
	}

	for _, test := range testCases {
		t.Run(test.caseName, func(t *testing.T) {
			// Arrange
			insertMembers(t, db)

			// Action
			err := repo.Delete(test.ctx, test.id)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			var count int
			err = db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM members").Scan(&count)
			asserts.NoError(t, err)
			asserts.Equals(t, test.expectedCountInDB, count)
		})
	}
}
