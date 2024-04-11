package category

import (
	"cms/internal/model"
	"cms/pkg/asserts"
	"cms/pkg/postgres"
	"cms/pkg/postgres/testingdb"
	"context"
	"errors"
	"testing"
)

func insertCategories(t *testing.T, db *postgres.Database) {
	t.Helper()
	ctx := context.Background()
	t.Cleanup(func() {
		_, err := db.ExecContext(ctx, "DELETE FROM category")
		asserts.NoError(t, err)
	})

	stmt, err := db.PrepareContext(ctx, "INSERT INTO category (id, name, slug) VALUES ($1, $2, $3)")
	asserts.NoError(t, err)

	for _, args := range [][]any{
		{"42", "first category", "first-category"},
		{"43", "second category", "second-category"},
	} {
		_, err := stmt.Exec(args...)
		asserts.NoError(t, err)
	}
}

func TestRepository_Create(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertCategories(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name        string
		ctx         context.Context
		data        model.Category
		expectedErr error
	}{
		{
			name: "success",
			ctx:  context.Background(),
			data: model.Category{ //nolint:exhaustruct
				Name: "new category",
				Slug: "new-category",
			},
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			data:        model.Category{}, //nolint:exhaustruct
			expectedErr: errors.New("query: context canceled"),
		},
		{
			name: "category name already exists",
			ctx:  context.Background(),
			data: model.Category{ //nolint:exhaustruct
				Name: "first category",
				Slug: "first-category1",
			},
			expectedErr: model.ErrCategoryNameAlreadyExists,
		},
		{
			name: "category slug already exists",
			ctx:  context.Background(),
			data: model.Category{ //nolint:exhaustruct
				Name: "first category1",
				Slug: "first-category",
			},
			expectedErr: model.ErrCategorySlugAlreadyExists,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Action
			err := repo.Create(test.ctx, test.data)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestRepository_GetByID(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertCategories(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		categoryID   model.CategoryID
		expectedData *model.Category
		expectedErr  error
	}{
		{
			name:       "success",
			ctx:        context.Background(),
			categoryID: model.CategoryID("42"),
			expectedData: &model.Category{
				ID:   "42",
				Name: "first category",
				Slug: "first-category",
			},
			expectedErr: nil,
		},
		{
			name:        "not found",
			ctx:         context.Background(),
			categoryID:  model.CategoryID("1"),
			expectedErr: model.ErrCategoryNotFound,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			categoryID:  model.CategoryID("42"),
			expectedErr: errors.New("query: context canceled"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Action
			data, err := repo.GetByID(test.ctx, test.categoryID)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			asserts.Equals(t, test.expectedData, data)
		})
	}
}

func TestRepository_GetBySlug(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertCategories(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		categorySlug model.CategorySlug
		expectedData *model.Category
		expectedErr  error
	}{
		{
			name:         "success",
			ctx:          context.Background(),
			categorySlug: "first-category",
			expectedData: &model.Category{
				ID:   "42",
				Name: "first category",
				Slug: "first-category",
			},
			expectedErr: nil,
		},
		{
			name:         "not found",
			ctx:          context.Background(),
			categorySlug: "not-found",
			expectedErr:  model.ErrCategoryNotFound,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			categorySlug: "first-category",
			expectedErr:  errors.New("query: context canceled"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Action
			data, err := repo.GetBySlug(test.ctx, test.categorySlug)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			asserts.Equals(t, test.expectedData, data)
		})
	}
}

func TestRepository_All(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertCategories(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		expectedData []model.Category
		expectedErr  error
	}{
		{
			name: "success",
			ctx:  context.Background(),
			expectedData: []model.Category{
				{ID: "42", Name: "first category", Slug: "first-category"},
				{ID: "43", Name: "second category", Slug: "second-category"},
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
			expectedData: []model.Category(nil),
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
		data            model.Category
		expectedErr     error
		dataAfterUpdate category
	}{
		{
			name: "success",
			ctx:  context.Background(),
			data: model.Category{
				ID:   "42",
				Name: "new category",
				Slug: "new-category",
			},
			dataAfterUpdate: category{
				ID:   "42",
				Name: "new category",
				Slug: "new-category",
			},
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			data: model.Category{
				ID:   "42",
				Name: "first category",
				Slug: "first-category",
			},
			expectedErr: errors.New("query: context canceled"),
			dataAfterUpdate: category{
				ID:   "42",
				Name: "first category",
				Slug: "first-category",
			},
		},
		{
			name: "category name already exists",
			ctx:  context.Background(),
			data: model.Category{
				ID:   "42",
				Name: "second category",
				Slug: "second-category1",
			},
			expectedErr: model.ErrCategoryNameAlreadyExists,
			dataAfterUpdate: category{
				ID:   "42",
				Name: "first category",
				Slug: "first-category",
			},
		},
		{
			name: "category slug already exists",
			ctx:  context.Background(),
			data: model.Category{
				ID:   "42",
				Name: "second category1",
				Slug: "second-category",
			},
			expectedErr: model.ErrCategorySlugAlreadyExists,
			dataAfterUpdate: category{
				ID:   "42",
				Name: "first category",
				Slug: "first-category",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			insertCategories(t, db)

			// Action
			err := repo.Update(test.ctx, test.data)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			var after category
			err = db.QueryRowContext(
				context.Background(),
				"SELECT id, name, slug FROM category WHERE id = 42",
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
		id                model.CategoryID
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
			insertCategories(t, db)

			// Action
			err := repo.Delete(test.ctx, test.id)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			var count int
			err = db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM category").Scan(&count)
			asserts.NoError(t, err)
			asserts.Equals(t, test.expectedCountInDB, count)
		})
	}
}
