package article

import (
	"cms/internal/model"
	"cms/pkg/asserts"
	"cms/pkg/postgres"
	"cms/pkg/postgres/testingdb"
	"context"
	"errors"
	"testing"
	"time"
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

func insertArticles(t *testing.T, db *postgres.Database) {
	t.Helper()
	ctx := context.Background()
	t.Cleanup(func() {
		_, err := db.ExecContext(ctx, "DELETE FROM article")
		asserts.NoError(t, err)
	})

	stmt, err := db.PrepareContext(
		ctx,
		"INSERT INTO article (id, title, slug, category_id, short_body, body, published_at) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7)",
	)
	asserts.NoError(t, err)

	for _, args := range [][]any{
		{
			"12",
			"first article title",
			"first-article",
			42,
			"first article short body",
			"first article body",
			time.Date(2022, time.January, 1, 23, 12, 4, 0, time.UTC),
		},
		{

			"13",
			"second article title",
			"second-article",
			43,
			"second article short body",
			"second article body",
			time.Date(2022, time.January, 2, 10, 30, 0, 0, time.UTC),
		},
	} {
		_, err := stmt.Exec(args...)
		asserts.NoError(t, err)
	}
}

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

func insertArticleTags(t *testing.T, db *postgres.Database) {
	t.Helper()
	ctx := context.Background()
	t.Cleanup(func() {
		_, err := db.ExecContext(ctx, "DELETE FROM article_tag")
		asserts.NoError(t, err)
	})
	stmt, err := db.PrepareContext(
		ctx,
		"INSERT INTO article_tag(article_id, tag_id) VALUES ($1, $2)",
	)
	asserts.NoError(t, err)
	for _, args := range [][]any{
		{
			"12",
			"42",
		},
	} {
		_, err := stmt.Exec(args...)
		asserts.NoError(t, err)
	}
}

func TestRepository_Create(t *testing.T) {
	db := testingdb.Connect(t)
	insertCategories(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name        string
		ctx         context.Context
		data        model.Article
		expectedErr error
		idIsEmpty   bool
	}{
		{
			name: "success",
			ctx:  context.Background(),
			data: model.Article{ //nolint:exhaustruct
				Slug:        "new-article",
				Title:       "new article",
				CategoryID:  "42",
				ShortBody:   "new article short body",
				Body:        "new article body",
				PublishedAt: time.Date(2022, time.May, 14, 13, 56, 14, 0, time.UTC),
			},
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			data:        model.Article{}, //nolint:exhaustruct
			expectedErr: errors.New("query: context canceled"),
			idIsEmpty:   true,
		},
		{
			name: "article title already exists",
			ctx:  context.Background(),
			data: model.Article{ //nolint:exhaustruct
				Title:       "first article title",
				Slug:        "first-article1",
				CategoryID:  "42",
				ShortBody:   "new article short body",
				Body:        "new article body",
				PublishedAt: time.Date(2022, time.May, 14, 13, 56, 14, 0, time.UTC),
			},
			expectedErr: model.ErrArticleTitleAlreadyExists,
			idIsEmpty:   true,
		},
		{
			name: "article slug already exists",
			ctx:  context.Background(),
			data: model.Article{ //nolint:exhaustruct
				Title:       "first article1",
				Slug:        "first-article",
				CategoryID:  "42",
				ShortBody:   "new article short body",
				Body:        "new article body",
				PublishedAt: time.Date(2022, time.May, 14, 13, 56, 14, 0, time.UTC),
			},
			expectedErr: model.ErrArticleSlugAlreadyExists,
			idIsEmpty:   true,
		},
		{
			name: "category not found",
			ctx:  context.Background(),
			data: model.Article{ //nolint:exhaustruct
				Title:       "new article",
				Slug:        "new-article",
				CategoryID:  "1",
				ShortBody:   "new article short body",
				Body:        "new article body",
				PublishedAt: time.Date(2022, time.May, 14, 13, 56, 14, 0, time.UTC),
			},
			expectedErr: model.ErrCategoryNotFound,
			idIsEmpty:   true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			insertArticles(t, db)

			// Action
			id, err := repo.Create(test.ctx, test.data)

			// Assert
			asserts.Equals(t, test.idIsEmpty, id == "", "id is empty")
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestRepository_Update(t *testing.T) {
	db := testingdb.Connect(t)
	insertCategories(t, db)
	repo := NewRepository(db)

	oldArticle := model.Article{
		ID:          "12",
		Slug:        "first-article",
		Title:       "first article title",
		CategoryID:  "42",
		ShortBody:   "first article short body",
		Body:        "first article body",
		PublishedAt: time.Date(2022, time.January, 1, 23, 12, 4, 0, time.UTC),
	}

	newArticle := model.Article{
		ID:          "12",
		Slug:        "first-article-changed",
		Title:       "first article changed",
		CategoryID:  "43",
		ShortBody:   "first article short body changed",
		Body:        "first article body changed",
		PublishedAt: time.Date(2023, time.February, 2, 24, 13, 5, 0, time.UTC),
	}

	testCases := []struct {
		name                    string
		ctx                     context.Context
		data                    model.Article
		expectedErr             error
		expectedDataAfterChange model.Article
	}{
		{
			name:                    "success",
			ctx:                     context.Background(),
			data:                    newArticle,
			expectedDataAfterChange: newArticle,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			data:                    newArticle,
			expectedErr:             errors.New("query: context canceled"),
			expectedDataAfterChange: oldArticle,
		},
		{
			name: "article title already exists",
			ctx:  context.Background(),
			data: func() model.Article {
				data := newArticle
				data.Title = "second article title"
				return data
			}(),
			expectedErr:             model.ErrArticleTitleAlreadyExists,
			expectedDataAfterChange: oldArticle,
		},
		{
			name: "article slug already exists",
			ctx:  context.Background(),
			data: func() model.Article {
				data := newArticle
				data.Slug = "second-article"
				return data
			}(),
			expectedErr:             model.ErrArticleSlugAlreadyExists,
			expectedDataAfterChange: oldArticle,
		},
		{
			name: "category not found",
			ctx:  context.Background(),
			data: func() model.Article {
				data := newArticle
				data.CategoryID = "1"
				return data
			}(),
			expectedErr:             model.ErrCategoryNotFound,
			expectedDataAfterChange: oldArticle,
		},
		{
			name: "article not found",
			ctx:  context.Background(),
			data: func() model.Article {
				data := newArticle
				data.ID = "1"
				return data
			}(),
			expectedErr:             nil,
			expectedDataAfterChange: oldArticle,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			insertArticles(t, db)

			// Action
			err := repo.Update(test.ctx, test.data)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			articleAfterUpdate, err := repo.FindByID(context.Background(), oldArticle.ID)
			asserts.NoError(t, err)
			asserts.Equals(t, test.expectedDataAfterChange, articleAfterUpdate)
		})
	}
}

func TestRepository_FindByID(t *testing.T) {
	// Arrange
	firstArticle := model.Article{
		ID:          "12",
		Slug:        "first-article",
		Title:       "first article title",
		CategoryID:  "42",
		ShortBody:   "first article short body",
		Body:        "first article body",
		PublishedAt: time.Date(2022, time.January, 1, 23, 12, 4, 0, time.UTC),
	}
	db := testingdb.Connect(t)
	insertCategories(t, db)
	insertArticles(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		id           model.ArticleID
		expectedData model.Article
		expectedErr  error
	}{
		{
			name:         "success",
			ctx:          context.Background(),
			id:           firstArticle.ID,
			expectedData: firstArticle,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			id:           firstArticle.ID,
			expectedData: model.Article{}, //nolint:exhaustruct
			expectedErr:  context.Canceled,
		},
		{
			name:         "article not found",
			ctx:          context.Background(),
			id:           "1",
			expectedData: model.Article{}, //nolint:exhaustruct
			expectedErr:  model.ErrArticleNotFound,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Action
			data, err := repo.FindByID(test.ctx, test.id)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			asserts.Equals(t, test.expectedData, data)
		})
	}
}

func TestRepository_FindBySlug(t *testing.T) {
	// Arrange
	firstArticle := model.Article{
		ID:          "12",
		Slug:        "first-article",
		Title:       "first article title",
		CategoryID:  "42",
		ShortBody:   "first article short body",
		Body:        "first article body",
		PublishedAt: time.Date(2022, time.January, 1, 23, 12, 4, 0, time.UTC),
	}
	db := testingdb.Connect(t)
	insertCategories(t, db)
	insertArticles(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		slug         model.ArticleSlug
		expectedData model.Article
		expectedErr  error
	}{
		{
			name:         "success",
			ctx:          context.Background(),
			slug:         firstArticle.Slug,
			expectedData: firstArticle,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			slug:         firstArticle.Slug,
			expectedData: model.Article{}, //nolint:exhaustruct
			expectedErr:  context.Canceled,
		},
		{
			name:         "article not found",
			ctx:          context.Background(),
			slug:         model.ArticleSlug("not-found"),
			expectedData: model.Article{}, //nolint:exhaustruct
			expectedErr:  model.ErrArticleNotFound,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Action
			data, err := repo.FindBySlug(test.ctx, test.slug)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			asserts.Equals(t, test.expectedData, data)
		})
	}
}

func TestRepository_FindShortInfo(t *testing.T) {
	// Arrange
	firstArticle := model.ArticleShortInfo{
		ID:          "12",
		Slug:        "first-article",
		Title:       "first article title",
		CategoryID:  "42",
		PublishedAt: time.Date(2022, time.January, 1, 23, 12, 4, 0, time.UTC),
	}

	secondArticle := model.ArticleShortInfo{
		ID:          "13",
		Slug:        "second-article",
		Title:       "second article title",
		CategoryID:  "43",
		PublishedAt: time.Date(2022, time.January, 2, 10, 30, 0, 0, time.UTC),
	}
	db := testingdb.Connect(t)
	insertCategories(t, db)
	insertArticles(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		offset       int64
		limit        int64
		expectedData []model.ArticleShortInfo
		expectedErr  error
	}{
		{
			name:         "success",
			ctx:          context.Background(),
			offset:       0,
			limit:        2,
			expectedData: []model.ArticleShortInfo{secondArticle, firstArticle},
			expectedErr:  nil,
		},
		{
			name:         "offset = 1",
			ctx:          context.Background(),
			offset:       1,
			limit:        2,
			expectedData: []model.ArticleShortInfo{firstArticle},
			expectedErr:  nil,
		},
		{
			name:         "limit = 1",
			ctx:          context.Background(),
			offset:       0,
			limit:        1,
			expectedData: []model.ArticleShortInfo{secondArticle},
			expectedErr:  nil,
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
			// Action
			data, err := repo.FindShortInfo(test.ctx, test.offset, test.limit)

			// Assert
			asserts.Equals(t, test.expectedData, data)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestRepository_FindShortInfoWithCategory(t *testing.T) {
	// Arrange
	firstArticle := model.ArticleShortInfoWithCategory{
		ID:    "12",
		Slug:  "first-article",
		Title: "first article title",
		Category: model.Category{
			ID:   "42",
			Slug: "first-category",
			Name: "first category",
		},
		ShortBody:   "first article short body",
		PublishedAt: time.Date(2022, time.January, 1, 23, 12, 4, 0, time.UTC),
	}

	secondArticle := model.ArticleShortInfoWithCategory{
		ID:    "13",
		Slug:  "second-article",
		Title: "second article title",
		Category: model.Category{
			ID:   "43",
			Slug: "second-category",
			Name: "second category",
		},
		ShortBody:   "second article short body",
		PublishedAt: time.Date(2022, time.January, 2, 10, 30, 0, 0, time.UTC),
	}
	db := testingdb.Connect(t)
	insertCategories(t, db)
	insertArticles(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		offset       int64
		limit        int64
		expectedData []model.ArticleShortInfoWithCategory
		expectedErr  error
	}{
		{
			name:         "success",
			ctx:          context.Background(),
			offset:       0,
			limit:        2,
			expectedData: []model.ArticleShortInfoWithCategory{secondArticle, firstArticle},
			expectedErr:  nil,
		},
		{
			name:         "offset = 1",
			ctx:          context.Background(),
			offset:       1,
			limit:        2,
			expectedData: []model.ArticleShortInfoWithCategory{firstArticle},
			expectedErr:  nil,
		},
		{
			name:         "limit = 1",
			ctx:          context.Background(),
			offset:       0,
			limit:        1,
			expectedData: []model.ArticleShortInfoWithCategory{secondArticle},
			expectedErr:  nil,
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
			// Action
			data, err := repo.FindShortInfoWithCategory(test.ctx, test.offset, test.limit)

			// Assert
			asserts.Equals(t, test.expectedData, data)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestRepository_FindShortInfoWithCategoryForCategory(t *testing.T) {
	// Arrange
	firstArticle := model.ArticleShortInfoWithCategory{
		ID:    "12",
		Slug:  "first-article",
		Title: "first article title",
		Category: model.Category{
			ID:   "42",
			Slug: "first-category",
			Name: "first category",
		},
		ShortBody:   "first article short body",
		PublishedAt: time.Date(2022, time.January, 1, 23, 12, 4, 0, time.UTC),
	}
	db := testingdb.Connect(t)
	insertCategories(t, db)
	insertArticles(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		categorySlug model.CategorySlug
		offset       int64
		limit        int64
		expectedData []model.ArticleShortInfoWithCategory
		expectedErr  error
	}{
		{
			name:         "success",
			ctx:          context.Background(),
			categorySlug: firstArticle.Category.Slug,
			offset:       0,
			limit:        2,
			expectedData: []model.ArticleShortInfoWithCategory{firstArticle},
			expectedErr:  nil,
		},
		{
			name:         "offset = 1",
			ctx:          context.Background(),
			categorySlug: firstArticle.Category.Slug,
			offset:       1,
			limit:        2,
			expectedData: nil,
			expectedErr:  nil,
		},
		{
			name:         "limit = 1",
			ctx:          context.Background(),
			categorySlug: firstArticle.Category.Slug,
			offset:       0,
			limit:        1,
			expectedData: []model.ArticleShortInfoWithCategory{firstArticle},
			expectedErr:  nil,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			categorySlug: firstArticle.Category.Slug,
			offset:       0,
			limit:        2,
			expectedData: nil,
			expectedErr:  errors.New("query: context canceled"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Action
			data, err := repo.FindShortInfoWithCategoryForCategory(
				test.ctx,
				test.categorySlug,
				test.offset,
				test.limit,
			)

			// Assert
			asserts.Equals(t, test.expectedData, data)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestRepository_FindShortInfoWithCategoryForTag(t *testing.T) {
	// Arrange
	firstArticle := model.ArticleShortInfoWithCategory{
		ID:    "12",
		Slug:  "first-article",
		Title: "first article title",
		Category: model.Category{
			ID:   "42",
			Slug: "first-category",
			Name: "first category",
		},
		ShortBody:   "first article short body",
		PublishedAt: time.Date(2022, time.January, 1, 23, 12, 4, 0, time.UTC),
	}
	tagSlug := model.TagSlug("first-tag")
	db := testingdb.Connect(t)
	insertCategories(t, db)
	insertArticles(t, db)
	insertTags(t, db)
	insertArticleTags(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		tagSlug      model.TagSlug
		offset       int64
		limit        int64
		expectedData []model.ArticleShortInfoWithCategory
		expectedErr  error
	}{
		{
			name:         "success",
			ctx:          context.Background(),
			tagSlug:      tagSlug,
			offset:       0,
			limit:        2,
			expectedData: []model.ArticleShortInfoWithCategory{firstArticle},
			expectedErr:  nil,
		},
		{
			name:         "offset = 1",
			ctx:          context.Background(),
			tagSlug:      tagSlug,
			offset:       1,
			limit:        2,
			expectedData: nil,
			expectedErr:  nil,
		},
		{
			name:         "limit = 1",
			ctx:          context.Background(),
			tagSlug:      tagSlug,
			offset:       0,
			limit:        1,
			expectedData: []model.ArticleShortInfoWithCategory{firstArticle},
			expectedErr:  nil,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			tagSlug:      tagSlug,
			offset:       0,
			limit:        2,
			expectedData: nil,
			expectedErr:  errors.New("query: context canceled"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Action
			data, err := repo.FindShortInfoWithCategoryForTag(
				test.ctx,
				test.tagSlug,
				test.offset,
				test.limit,
			)

			// Assert
			asserts.Equals(t, test.expectedData, data)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestRepository_FindShortInfoWithCategoryForQuery(t *testing.T) {
	// Arrange
	firstArticle := model.ArticleShortInfoWithCategory{
		ID:    "12",
		Slug:  "first-article",
		Title: "first article title",
		Category: model.Category{
			ID:   "42",
			Slug: "first-category",
			Name: "first category",
		},
		ShortBody:   "first article short body",
		PublishedAt: time.Date(2022, time.January, 1, 23, 12, 4, 0, time.UTC),
	}
	db := testingdb.Connect(t)
	insertCategories(t, db)
	insertArticles(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		query        string
		offset       int64
		limit        int64
		expectedData []model.ArticleShortInfoWithCategory
		expectedErr  error
	}{
		{
			name:         "success",
			ctx:          context.Background(),
			query:        "first",
			offset:       0,
			limit:        2,
			expectedData: []model.ArticleShortInfoWithCategory{firstArticle},
			expectedErr:  nil,
		},
		{
			name:         "ignore case",
			ctx:          context.Background(),
			query:        "FIRST",
			offset:       0,
			limit:        2,
			expectedData: []model.ArticleShortInfoWithCategory{firstArticle},
			expectedErr:  nil,
		},
		{
			name:         "title",
			ctx:          context.Background(),
			query:        "FIRST ARTICLE TITLE",
			offset:       0,
			limit:        2,
			expectedData: []model.ArticleShortInfoWithCategory{firstArticle},
			expectedErr:  nil,
		},
		{
			name:         "short body",
			ctx:          context.Background(),
			query:        "FIRST ARTICLE SHORT BODY",
			offset:       0,
			limit:        2,
			expectedData: []model.ArticleShortInfoWithCategory{firstArticle},
			expectedErr:  nil,
		},
		{
			name:         "body",
			ctx:          context.Background(),
			query:        "FIRST ARTICLE BODY",
			offset:       0,
			limit:        2,
			expectedData: []model.ArticleShortInfoWithCategory{firstArticle},
			expectedErr:  nil,
		},
		{
			name:         "offset = 1",
			ctx:          context.Background(),
			query:        "first",
			offset:       1,
			limit:        2,
			expectedData: nil,
			expectedErr:  nil,
		},
		{
			name:         "limit = 1",
			ctx:          context.Background(),
			query:        "first",
			offset:       0,
			limit:        1,
			expectedData: []model.ArticleShortInfoWithCategory{firstArticle},
			expectedErr:  nil,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			query:        "first",
			offset:       0,
			limit:        2,
			expectedData: nil,
			expectedErr:  errors.New("query: context canceled"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Action
			data, err := repo.FindShortInfoWithCategoryForQuery(
				test.ctx,
				test.query,
				test.offset,
				test.limit,
			)

			// Assert
			asserts.Equals(t, test.expectedData, data)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestRepository_Total(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertCategories(t, db)
	insertArticles(t, db)
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

func TestRepository_TotalForCategory(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertCategories(t, db)
	insertArticles(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		slug         model.CategorySlug
		expectedData int
		expectedErr  error
	}{
		{
			name:         "success",
			ctx:          context.Background(),
			slug:         model.CategorySlug("first-category"),
			expectedData: 1,
			expectedErr:  nil,
		},
		{
			name:         "not found",
			ctx:          context.Background(),
			slug:         model.CategorySlug("not-found"),
			expectedData: 0,
			expectedErr:  nil,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			slug:         model.CategorySlug("first-category"),
			expectedData: 0,
			expectedErr:  errors.New("query: context canceled"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Act
			data, err := repo.TotalForCategory(test.ctx, test.slug)

			// Assert
			asserts.Equals(t, test.expectedData, data)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestRepository_TotalForTag(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertCategories(t, db)
	insertArticles(t, db)
	insertTags(t, db)
	insertArticleTags(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		slug         model.TagSlug
		expectedData int
		expectedErr  error
	}{
		{
			name:         "success",
			ctx:          context.Background(),
			slug:         model.TagSlug("first-tag"),
			expectedData: 1,
			expectedErr:  nil,
		},
		{
			name:         "not found",
			ctx:          context.Background(),
			slug:         model.TagSlug("not-found"),
			expectedData: 0,
			expectedErr:  nil,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			slug:         model.TagSlug("first-tag"),
			expectedData: 0,
			expectedErr:  errors.New("query: context canceled"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Act
			data, err := repo.TotalForTag(test.ctx, test.slug)

			// Assert
			asserts.Equals(t, test.expectedData, data)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestRepository_TotalForQuery(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertCategories(t, db)
	insertArticles(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name         string
		ctx          context.Context
		query        string
		expectedData int
		expectedErr  error
	}{
		{
			name:         "success",
			ctx:          context.Background(),
			query:        "first",
			expectedData: 1,
			expectedErr:  nil,
		},
		{
			name:         "ignore case",
			ctx:          context.Background(),
			query:        "FIRST",
			expectedData: 1,
			expectedErr:  nil,
		},
		{
			name:         "not found",
			ctx:          context.Background(),
			query:        "not-found",
			expectedData: 0,
			expectedErr:  nil,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			query:        "first",
			expectedData: 0,
			expectedErr:  errors.New("query: context canceled"),
		},
		{
			name:         "title",
			ctx:          context.Background(),
			query:        "FIRST ARTICLE TITLE",
			expectedData: 1,
			expectedErr:  nil,
		},
		{
			name:         "short body",
			ctx:          context.Background(),
			query:        "FIRST ARTICLE SHORT BODY",
			expectedData: 1,
			expectedErr:  nil,
		},
		{
			name:         "body",
			ctx:          context.Background(),
			query:        "FIRST ARTICLE BODY",
			expectedData: 1,
			expectedErr:  nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Act
			data, err := repo.TotalForQuery(test.ctx, test.query)

			// Assert
			asserts.Equals(t, test.expectedData, data)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertCategories(t, db)
	insertTags(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		name             string
		ctx              context.Context
		id               model.ArticleID
		expectedErr      error
		countAfterDelete int
	}{
		{
			name:             "success",
			ctx:              context.Background(),
			id:               "12",
			expectedErr:      nil,
			countAfterDelete: 1,
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			id:               "12",
			expectedErr:      errors.New("query: context canceled"),
			countAfterDelete: 2,
		},
		{
			name:             "article not found",
			ctx:              context.Background(),
			id:               "11",
			expectedErr:      nil,
			countAfterDelete: 2,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			insertArticles(t, db)
			insertArticleTags(t, db)

			// Action
			err := repo.Delete(test.ctx, test.id)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			var count int
			err = db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM article").Scan(&count)
			asserts.NoError(t, err)
			asserts.Equals(t, test.countAfterDelete, count)
		})
	}
}
