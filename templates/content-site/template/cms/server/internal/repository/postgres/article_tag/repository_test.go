package article_tag

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
		{"24", "first category", "first-category"},
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
			"first article",
			"first-article",
			"24",
			"first article short body",
			"first article body",
			time.Date(2022, time.January, 1, 23, 12, 4, 0, time.UTC),
		},
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
			"43",
		},
	} {
		_, err := stmt.Exec(args...)
		asserts.NoError(t, err)
	}
}

func TestRepository_Create(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertTags(t, db)
	insertCategories(t, db)
	insertArticles(t, db)
	repo := NewRepository(db)

	type articleTag struct {
		ArticleID string
		TagID     string
	}

	testCases := []struct {
		name               string
		ctx                context.Context
		articleID          model.ArticleID
		tagID              model.TagID
		expectedErr        error
		recordsAfterCreate []articleTag
	}{
		{
			name:      "create article tag",
			ctx:       context.Background(),
			articleID: "12",
			tagID:     "42",
			recordsAfterCreate: []articleTag{
				{ArticleID: "12", TagID: "42"},
				{ArticleID: "12", TagID: "43"},
			},
		},
		{
			name:        "create duplicate article tag",
			ctx:         context.Background(),
			articleID:   "12",
			tagID:       "43",
			expectedErr: model.ErrArticleTagAlreadyExists,
			recordsAfterCreate: []articleTag{
				{ArticleID: "12", TagID: "43"},
			},
		},
		{
			name:        "create article tag with invalid article id",
			ctx:         context.Background(),
			articleID:   "13",
			tagID:       "42",
			expectedErr: model.ErrArticleNotFound,
			recordsAfterCreate: []articleTag{
				{ArticleID: "12", TagID: "43"},
			},
		},
		{
			name:        "create article tag with invalid tag id",
			ctx:         context.Background(),
			articleID:   "12",
			tagID:       "44",
			expectedErr: model.ErrTagNotFound,
			recordsAfterCreate: []articleTag{
				{ArticleID: "12", TagID: "43"},
			},
		},
		{
			name: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			articleID:   "12",
			tagID:       "42",
			expectedErr: errors.New("query: context canceled"),
			recordsAfterCreate: []articleTag{
				{ArticleID: "12", TagID: "43"},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			insertArticleTags(t, db)

			// Action
			err := repo.Create(test.ctx, test.articleID, test.tagID)

			// Arrange
			asserts.ErrorsEqual(t, test.expectedErr, err)
			rows, err := db.QueryContext(
				context.Background(),
				"SELECT article_id, tag_id FROM article_tag ORDER BY article_id, tag_id",
			)
			asserts.NoError(t, err)
			var items []articleTag
			for rows.Next() {
				var item articleTag
				err := rows.Scan(&item.ArticleID, &item.TagID)
				asserts.NoError(t, err)
				items = append(items, item)
			}
			asserts.Equals(t, test.recordsAfterCreate, items)
		})
	}
}

func TestRepository_GetTagIDs(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertTags(t, db)
	insertCategories(t, db)
	insertArticles(t, db)
	insertArticleTags(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		caseName    string
		ctx         context.Context
		articleID   model.ArticleID
		expected    []model.TagID
		expectedErr error
	}{
		{
			caseName:  "success",
			ctx:       context.Background(),
			articleID: "12",
			expected:  []model.TagID{"43"},
		},
		{
			caseName:  "article not found",
			ctx:       context.Background(),
			articleID: "13",
		},
		{
			caseName: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			articleID:   "12",
			expectedErr: errors.New("query: context canceled"),
		},
	}

	for _, test := range testCases {
		t.Run(test.caseName, func(t *testing.T) {
			// Action
			ids, err := repo.GetTagIDs(test.ctx, test.articleID)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			asserts.Equals(t, test.expected, ids)
		})
	}
}

func TestRepository_GetTagForArticle(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	insertTags(t, db)
	insertCategories(t, db)
	insertArticles(t, db)
	insertArticleTags(t, db)
	repo := NewRepository(db)

	testCases := []struct {
		caseName     string
		ctx          context.Context
		articleID    model.ArticleID
		expectedTags []model.Tag
		expectedErr  error
	}{
		{
			caseName:  "success",
			ctx:       context.Background(),
			articleID: "12",
			expectedTags: []model.Tag{
				{ID: "43", Name: "second tag", Slug: "second-tag"},
			},
		},
		{
			caseName:  "article not found",
			ctx:       context.Background(),
			articleID: "13",
		},
		{
			caseName: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			articleID:   "12",
			expectedErr: errors.New("query: context canceled"),
		},
	}

	for _, test := range testCases {
		t.Run(test.caseName, func(t *testing.T) {
			// Action
			tags, err := repo.GetTagForArticle(test.ctx, test.articleID)

			// Assert
			asserts.Equals(t, test.expectedTags, tags)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	db := testingdb.Connect(t)
	repo := NewRepository(db)

	testCases := []struct {
		caseName          string
		ctx               context.Context
		articleID         model.ArticleID
		tagID             model.TagID
		expectedErr       error
		expectedCountInDB int
	}{
		{
			caseName:          "success",
			ctx:               context.Background(),
			articleID:         "12",
			tagID:             "43",
			expectedCountInDB: 0,
		},
		{
			caseName:          "article not found",
			ctx:               context.Background(),
			articleID:         "13",
			tagID:             "43",
			expectedCountInDB: 1,
		},
		{
			caseName:          "tag not found",
			ctx:               context.Background(),
			articleID:         "12",
			tagID:             "41",
			expectedCountInDB: 1,
		},
		{
			caseName: "context cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			articleID:         "12",
			tagID:             "43",
			expectedErr:       errors.New("query: context canceled"),
			expectedCountInDB: 1,
		},
	}

	for _, test := range testCases {
		t.Run(test.caseName, func(t *testing.T) {
			// Arrange
			insertTags(t, db)
			insertCategories(t, db)
			insertArticles(t, db)
			insertArticleTags(t, db)

			// Action
			err := repo.Delete(test.ctx, test.articleID, test.tagID)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
			var count int
			err = db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM article_tag").Scan(&count)
			asserts.NoError(t, err)
			asserts.Equals(t, test.expectedCountInDB, count)
		})
	}
}
