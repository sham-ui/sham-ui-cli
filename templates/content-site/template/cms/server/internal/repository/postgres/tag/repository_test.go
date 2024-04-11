package tag

import (
	"cms/internal/model"
	"cms/pkg/asserts"
	"cms/pkg/postgres"
	"cms/pkg/postgres/testingdb"
	"context"
	"errors"
	"testing"
)

func insertTags(t *testing.T, db *postgres.Database) {
	t.Helper()
	ctx := context.Background()
	t.Cleanup(func() {
		_, err := db.ExecContext(ctx, "DELETE FROM tag")
		asserts.NoError(t, err)
	})

	stmt, err := db.PrepareContext(ctx, "INSERT INTO tag (id, name, slug) VALUES ($1, $2, $3)")
	asserts.NoError(t, err)

	for _, args := range [][]any{
		{"42", "first tag", "first-tag"},
		{"43", "second tag", "second-tag"},
	} {
		_, err := stmt.Exec(args...)
		asserts.NoError(t, err)
	}
}

func TestRepository_Create(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertTags(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name        string
		ctx         context.Context
		data        model.Tag
		expectedErr error
		idIsEmpty   bool
	}{
		{
			name: "success",
			ctx:  context.Background(),
			data: model.Tag{ //nolint:exhaustruct
				Name: "new tag",
				Slug: "new-tag",
			},
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			data:        model.Tag{}, //nolint:exhaustruct
			expectedErr: errors.New("query: context canceled"),
			idIsEmpty:   true,
		},
		{
			name: "tag name already exists",
			ctx:  context.Background(),
			data: model.Tag{ //nolint:exhaustruct
				Name: "first tag",
				Slug: "first-tag1",
			},
			expectedErr: model.ErrTagNameAlreadyExists,
			idIsEmpty:   true,
		},
		{
			name: "tag slug already exists",
			ctx:  context.Background(),
			data: model.Tag{ //nolint:exhaustruct
				Name: "first tag1",
				Slug: "first-tag",
			},
			expectedErr: model.ErrTagSlugAlreadyExists,
			idIsEmpty:   true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Action
			id, err := repo.Create(test.ctx, test.data)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			asserts.Equals(t, test.idIsEmpty, id == "", "id is empty")
		})
	}
}

func TestRepository_GetBySlug(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertTags(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		slug         model.TagSlug
		expectedData model.Tag
		expectedErr  error
	}{
		{
			name:         "success",
			ctx:          context.Background(),
			slug:         "first-tag",
			expectedData: model.Tag{ID: "42", Name: "first tag", Slug: "first-tag"},
			expectedErr:  nil,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			expectedData: model.Tag{}, //nolint:exhaustruct
			expectedErr:  errors.New("query: context canceled"),
		},
		{
			name:         "not found",
			ctx:          context.Background(),
			slug:         "not-found",
			expectedData: model.Tag{}, //nolint:exhaustruct
			expectedErr:  model.ErrTagNotFound,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Action
			data, err := repo.GetBySlug(test.ctx, test.slug)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			asserts.Equals(t, test.expectedData, data, "data")
		})
	}
}

func TestRepository_All(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertTags(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		expectedData []model.Tag
		expectedErr  error
	}{
		{
			name: "success",
			ctx:  context.Background(),
			expectedData: []model.Tag{
				{ID: "42", Name: "first tag", Slug: "first-tag"},
				{ID: "43", Name: "second tag", Slug: "second-tag"},
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
			expectedData: []model.Tag(nil),
			expectedErr:  errors.New("query: context canceled"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Action
			data, err := repo.All(test.ctx)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			asserts.Equals(t, test.expectedData, data, "data")
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
		data            model.Tag
		expectedErr     error
		dataAfterUpdate tag
	}{
		{
			name: "success",
			ctx:  context.Background(),
			data: model.Tag{
				ID:   "42",
				Name: "new tag",
				Slug: "new-tag",
			},
			dataAfterUpdate: tag{
				ID:   "42",
				Name: "new tag",
				Slug: "new-tag",
			},
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			data: model.Tag{
				ID:   "42",
				Name: "first tag",
				Slug: "first-tag",
			},
			expectedErr: errors.New("query: context canceled"),
			dataAfterUpdate: tag{
				ID:   "42",
				Name: "first tag",
				Slug: "first-tag",
			},
		},
		{
			name: "tag name already exists",
			ctx:  context.Background(),
			data: model.Tag{
				ID:   "42",
				Name: "second tag",
				Slug: "second-tag1",
			},
			expectedErr: model.ErrTagNameAlreadyExists,
			dataAfterUpdate: tag{
				ID:   "42",
				Name: "first tag",
				Slug: "first-tag",
			},
		},
		{
			name: "tag slug already exists",
			ctx:  context.Background(),
			data: model.Tag{
				ID:   "42",
				Name: "second tag1",
				Slug: "second-tag",
			},
			expectedErr: model.ErrTagSlugAlreadyExists,
			dataAfterUpdate: tag{
				ID:   "42",
				Name: "first tag",
				Slug: "first-tag",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			insertTags(t, db)

			// Action
			err := repo.Update(test.ctx, test.data)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			var after tag
			err = db.QueryRowContext(
				context.Background(),
				"SELECT id, name, slug FROM tag WHERE id = 42",
			).Scan(
				&after.ID,
				&after.Name,
				&after.Slug,
			)
			asserts.NoError(t, err)
			asserts.Equals(t, test.dataAfterUpdate, after, "after update")
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	db := testingdb.Connect(t)
	repo := NewRepository(db)

	testCases := []struct {
		caseName          string
		ctx               context.Context
		id                model.TagID
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
			insertTags(t, db)

			// Action
			err := repo.Delete(test.ctx, test.id)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			var count int
			err = db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM tag").Scan(&count)
			asserts.NoError(t, err)
			asserts.Equals(t, test.expectedCountInDB, count)
		})
	}
}
