package article

import (
	"cms/internal/model"
	"cms/pkg/asserts"
	"cms/pkg/transactor"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestService_FindShortInfo(t *testing.T) {
	article := model.ArticleShortInfo{
		ID:          "1",
		Slug:        "article-slug",
		Title:       "article title",
		CategoryID:  "42",
		PublishedAt: time.Date(2022, time.January, 1, 23, 12, 4, 0, time.UTC),
	}

	testCases := []struct {
		name         string
		ctx          context.Context
		offset       int64
		limit        int64
		articleRepo  func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository
		expectedData []model.ArticleShortInfo
		expectedErr  error
	}{
		{
			name:   "success",
			ctx:    context.Background(),
			offset: 10,
			limit:  20,
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("FindShortInfo", mock.Anything, int64(10), int64(20)).
					Return([]model.ArticleShortInfo{article}, nil).
					Once()
				return m
			},
			expectedData: []model.ArticleShortInfo{article},
		},
		{
			name:   "fail",
			ctx:    context.Background(),
			offset: 10,
			limit:  20,
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("FindShortInfo", mock.Anything, int64(10), int64(20)).
					Return(nil, errors.New("test")).
					Once()
				return m
			},
			expectedErr: errors.New("find short info: test"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			repo := New(
				transactor.NewDummyTransactor(),
				test.articleRepo(t),
				NewMockTagRepository(t),
				NewMockArticleTagRepository(t),
			)

			// Act
			data, err := repo.FindShortInfo(test.ctx, test.offset, test.limit)

			// Assert
			asserts.Equals(t, test.expectedData, data)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestService_FindByID(t *testing.T) {
	article := model.Article{ //nolint:exhaustruct
		ID:          "1",
		Slug:        "article-slug",
		Title:       "article title",
		CategoryID:  "42",
		PublishedAt: time.Date(2022, time.January, 1, 23, 12, 4, 0, time.UTC),
	}

	testCases := []struct {
		name         string
		ctx          context.Context
		id           model.ArticleID
		articleRepo  func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository
		expectedData model.Article
		expectedErr  error
	}{
		{
			name: "success",
			ctx:  context.Background(),
			id:   "1",
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("FindByID", mock.Anything, model.ArticleID("1")).
					Return(article, nil).
					Once()
				return m
			},
			expectedData: article,
		},
		{
			name: "fail",
			ctx:  context.Background(),
			id:   "1",
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("FindByID", mock.Anything, model.ArticleID("1")).
					Return(model.Article{}, errors.New("test")). //nolint:exhaustruct
					Once()
				return m
			},
			expectedErr: errors.New("find by id: test"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			repo := New(
				transactor.NewDummyTransactor(),
				test.articleRepo(t),
				NewMockTagRepository(t),
				NewMockArticleTagRepository(t),
			)

			// Act
			data, err := repo.FindByID(test.ctx, test.id)

			// Assert
			asserts.Equals(t, test.expectedData, data)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestService_GetTags(t *testing.T) {
	firstTag := model.Tag{
		ID:   "12",
		Slug: "first-tag",
		Name: "First tag",
	}
	secondTag := model.Tag{
		ID:   "34",
		Slug: "second-tag",
		Name: "Second tag",
	}

	testCases := []struct {
		name           string
		ctx            context.Context
		id             model.ArticleID
		articleTagRepo func(t mockConstructorTestingTNewMockArticleTagRepository) *MockArticleTagRepository
		expectedData   []model.Tag
		expectedErr    error
	}{
		{
			name: "success",
			ctx:  context.Background(),
			id:   "1",
			articleTagRepo: func(t mockConstructorTestingTNewMockArticleTagRepository) *MockArticleTagRepository {
				m := NewMockArticleTagRepository(t)
				m.
					On("GetTagForArticle", mock.Anything, model.ArticleID("1")).
					Return([]model.Tag{firstTag, secondTag}, nil).
					Once()
				return m
			},
			expectedData: []model.Tag{firstTag, secondTag},
		},
		{
			name: "fail",
			ctx:  context.Background(),
			id:   "1",
			articleTagRepo: func(t mockConstructorTestingTNewMockArticleTagRepository) *MockArticleTagRepository {
				m := NewMockArticleTagRepository(t)
				m.
					On("GetTagForArticle", mock.Anything, model.ArticleID("1")).
					Return(nil, errors.New("test")).
					Once()
				return m
			},
			expectedErr: errors.New("get tags for article: test"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			repo := New(
				transactor.NewDummyTransactor(),
				NewMockArticleRepository(t),
				NewMockTagRepository(t),
				test.articleTagRepo(t),
			)

			// Act
			data, err := repo.GetTags(test.ctx, test.id)

			// Assert
			asserts.Equals(t, test.expectedData, data)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestService_Total(t *testing.T) {
	testCases := []struct {
		name         string
		ctx          context.Context
		articleRepo  func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository
		expectedData int
		expectedErr  error
	}{
		{
			name: "success",
			ctx:  context.Background(),
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Total", mock.Anything).
					Return(10, nil).
					Once()
				return m
			},
			expectedData: 10,
		},
		{
			name: "fail",
			ctx:  context.Background(),
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Total", mock.Anything).
					Return(0, errors.New("test")).
					Once()
				return m
			},
			expectedErr: errors.New("total: test"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			repo := New(
				transactor.NewDummyTransactor(),
				test.articleRepo(t),
				NewMockTagRepository(t),
				NewMockArticleTagRepository(t),
			)

			// Act
			data, err := repo.Total(test.ctx)

			// Assert
			asserts.Equals(t, test.expectedData, data)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestService_Create(t *testing.T) {
	newTag := model.Tag{Name: "new tag", Slug: "new-tag"} //nolint:exhaustruct
	existedTag := model.Tag{ID: "1", Name: "existed tag", Slug: "existed-tag"}
	tags := []model.Tag{newTag, existedTag}
	article := model.Article{ //nolint:exhaustruct
		Slug:        "new-article",
		Title:       "new article",
		CategoryID:  "42",
		ShortBody:   "new article short body",
		Body:        "full content of new article",
		PublishedAt: time.Date(2024, time.November, 10, 23, 0, 0, 0, time.UTC),
	}

	testCases := []struct {
		name           string
		ctx            context.Context
		article        model.Article
		tags           []model.Tag
		articleRepo    func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository
		tagRepo        func(t mockConstructorTestingTNewMockTagRepository) *MockTagRepository
		articleTagRepo func(t mockConstructorTestingTNewMockArticleTagRepository) *MockArticleTagRepository
		expectedID     model.ArticleID
		expectedErr    error
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			article: article,
			tags:    tags,
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Create", mock.Anything, article).
					Return(model.ArticleID("4"), nil).
					Once()
				return m
			},
			tagRepo: func(t mockConstructorTestingTNewMockTagRepository) *MockTagRepository {
				m := NewMockTagRepository(t)
				m.
					On("GetBySlug", mock.Anything, newTag.Slug).
					Return(model.Tag{}, model.ErrTagNotFound). //nolint:exhaustruct
					Once()
				m.
					On("Create", mock.Anything, newTag).
					Return(model.TagID("2"), nil).
					Once()
				m.
					On("GetBySlug", mock.Anything, existedTag.Slug).
					Return(existedTag, nil).
					Once()
				return m
			},
			articleTagRepo: func(t mockConstructorTestingTNewMockArticleTagRepository) *MockArticleTagRepository {
				m := NewMockArticleTagRepository(t)
				m.
					On("Create", mock.Anything, model.ArticleID("4"), model.TagID("2")).
					Return(nil).
					Once()
				m.
					On("Create", mock.Anything, model.ArticleID("4"), model.TagID("1")).
					Return(nil).
					Once()
				return m
			},
			expectedID:  "4",
			expectedErr: nil,
		},
		{
			name:    "create article fail",
			ctx:     context.Background(),
			article: article,
			tags:    tags,
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Create", mock.Anything, article).
					Return(model.ArticleID(""), errors.New("test")).
					Once()
				return m
			},
			tagRepo:        NewMockTagRepository,
			articleTagRepo: NewMockArticleTagRepository,
			expectedErr:    errors.New("within transaction: create article: test"),
		},
		{
			name:    "get tag by slug fail",
			ctx:     context.Background(),
			article: article,
			tags:    tags,
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Create", mock.Anything, article).
					Return(model.ArticleID("4"), nil).
					Once()
				return m
			},
			tagRepo: func(t mockConstructorTestingTNewMockTagRepository) *MockTagRepository {
				m := NewMockTagRepository(t)
				m.
					On("GetBySlug", mock.Anything, newTag.Slug).
					Return(model.Tag{}, errors.New("test")). //nolint:exhaustruct
					Once()
				return m
			},
			articleTagRepo: NewMockArticleTagRepository,
			expectedErr: errors.New(
				"within transaction: get or create tags: get tag by slug(slug=new-tag): test",
			),
		},
		{
			name:    "create tag fail",
			ctx:     context.Background(),
			article: article,
			tags:    tags,
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Create", mock.Anything, article).
					Return(model.ArticleID("4"), nil).
					Once()
				return m
			},
			tagRepo: func(t mockConstructorTestingTNewMockTagRepository) *MockTagRepository {
				m := NewMockTagRepository(t)
				m.
					On("GetBySlug", mock.Anything, newTag.Slug).
					Return(model.Tag{}, model.ErrTagNotFound). //nolint:exhaustruct
					Once()
				m.
					On("Create", mock.Anything, newTag).
					Return(model.TagID(""), errors.New("test")).
					Once()
				return m
			},
			articleTagRepo: NewMockArticleTagRepository,
			expectedErr:    errors.New("within transaction: get or create tags: create tag(slug=new-tag): test"),
		},
		{
			name:    "create article tag fail",
			ctx:     context.Background(),
			article: article,
			tags:    tags,
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Create", mock.Anything, article).
					Return(model.ArticleID("4"), nil).
					Once()
				return m
			},
			tagRepo: func(t mockConstructorTestingTNewMockTagRepository) *MockTagRepository {
				m := NewMockTagRepository(t)
				m.
					On("GetBySlug", mock.Anything, newTag.Slug).
					Return(model.Tag{}, model.ErrTagNotFound). //nolint:exhaustruct
					Once()
				m.
					On("Create", mock.Anything, newTag).
					Return(model.TagID("2"), nil).
					Once()
				m.
					On("GetBySlug", mock.Anything, existedTag.Slug).
					Return(existedTag, nil).
					Once()
				return m
			},
			articleTagRepo: func(t mockConstructorTestingTNewMockArticleTagRepository) *MockArticleTagRepository {
				m := NewMockArticleTagRepository(t)
				m.
					On("Create", mock.Anything, model.ArticleID("4"), model.TagID("2")).
					Return(errors.New("test")).
					Once()
				return m
			},
			expectedErr: errors.New("within transaction: create article tag(tag_id=2, article_id=4): test"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			srv := New(
				transactor.NewDummyTransactor(),
				test.articleRepo(t),
				test.tagRepo(t),
				test.articleTagRepo(t),
			)

			// Action
			id, err := srv.Create(test.ctx, test.article, test.tags)

			// Assert
			asserts.Equals(t, test.expectedID, id)
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestService_Update(t *testing.T) { //nolint:maintidx
	newTag := model.Tag{Name: "new tag", Slug: "new-tag"} //nolint:exhaustruct
	firstExistedTag := model.Tag{ID: "1", Name: "existed tag1", Slug: "existed-tag-1"}
	secondExistedTag := model.Tag{ID: "2", Name: "existed tag2", Slug: "existed-tag-2"}
	tags := []model.Tag{newTag, firstExistedTag}
	article := model.Article{
		ID:          "12",
		Slug:        "new-article",
		Title:       "new article",
		CategoryID:  "42",
		ShortBody:   "new article short body",
		Body:        "full content of new article",
		PublishedAt: time.Date(2024, time.November, 10, 23, 0, 0, 0, time.UTC),
	}
	existedArticleTagIDs := []model.TagID{firstExistedTag.ID, secondExistedTag.ID}

	testCases := []struct {
		name           string
		ctx            context.Context
		article        model.Article
		tags           []model.Tag
		articleRepo    func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository
		tagRepo        func(t mockConstructorTestingTNewMockTagRepository) *MockTagRepository
		articleTagRepo func(t mockConstructorTestingTNewMockArticleTagRepository) *MockArticleTagRepository
		expectedErr    error
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			article: article,
			tags:    tags,
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Update", mock.Anything, article).
					Return(nil).
					Once()
				return m
			},
			tagRepo: func(t mockConstructorTestingTNewMockTagRepository) *MockTagRepository {
				m := NewMockTagRepository(t)
				m.
					On("GetBySlug", mock.Anything, newTag.Slug).
					Return(model.Tag{}, model.ErrTagNotFound). //nolint:exhaustruct
					Once()
				m.
					On("Create", mock.Anything, newTag).
					Return(model.TagID("3"), nil).
					Once()
				m.
					On("GetBySlug", mock.Anything, firstExistedTag.Slug).
					Return(firstExistedTag, nil).
					Once()
				return m
			},
			articleTagRepo: func(t mockConstructorTestingTNewMockArticleTagRepository) *MockArticleTagRepository {
				m := NewMockArticleTagRepository(t)
				m.
					On("GetTagIDs", mock.Anything, article.ID).
					Return(existedArticleTagIDs, nil).
					Once()
				m.
					On("Create", mock.Anything, article.ID, model.TagID("3")).
					Return(nil).
					Once()
				m.
					On("Delete", mock.Anything, article.ID, secondExistedTag.ID).
					Return(nil).
					Once()
				return m
			},
			expectedErr: nil,
		},
		{
			name:    "update article fail",
			ctx:     context.Background(),
			article: article,
			tags:    tags,
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Update", mock.Anything, article).
					Return(errors.New("test")).
					Once()
				return m
			},
			tagRepo:        NewMockTagRepository,
			articleTagRepo: NewMockArticleTagRepository,
			expectedErr:    errors.New("within transaction: update article: test"),
		},
		{
			name:    "get existed tag ids fail",
			ctx:     context.Background(),
			article: article,
			tags:    tags,
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Update", mock.Anything, article).
					Return(nil).
					Once()
				return m
			},
			tagRepo: NewMockTagRepository,
			articleTagRepo: func(t mockConstructorTestingTNewMockArticleTagRepository) *MockArticleTagRepository {
				m := NewMockArticleTagRepository(t)
				m.
					On("GetTagIDs", mock.Anything, article.ID).
					Return(nil, errors.New("test")).
					Once()
				return m
			},
			expectedErr: errors.New("within transaction: get article tag ids(article_id=12): test"),
		},
		{
			name:    "get tag by slug fail",
			ctx:     context.Background(),
			article: article,
			tags:    tags,
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Update", mock.Anything, article).
					Return(nil).
					Once()
				return m
			},
			tagRepo: func(t mockConstructorTestingTNewMockTagRepository) *MockTagRepository {
				m := NewMockTagRepository(t)
				m.
					On("GetBySlug", mock.Anything, newTag.Slug).
					Return(model.Tag{}, errors.New("test")). //nolint:exhaustruct
					Once()
				return m
			},
			articleTagRepo: func(t mockConstructorTestingTNewMockArticleTagRepository) *MockArticleTagRepository {
				m := NewMockArticleTagRepository(t)
				m.
					On("GetTagIDs", mock.Anything, article.ID).
					Return(existedArticleTagIDs, nil).
					Once()
				return m
			},
			expectedErr: errors.New("within transaction: get or create tags: get tag by slug(slug=new-tag): test"),
		},
		{
			name:    "create tag fail",
			ctx:     context.Background(),
			article: article,
			tags:    tags,
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Update", mock.Anything, article).
					Return(nil).
					Once()
				return m
			},
			tagRepo: func(t mockConstructorTestingTNewMockTagRepository) *MockTagRepository {
				m := NewMockTagRepository(t)
				m.
					On("GetBySlug", mock.Anything, newTag.Slug).
					Return(model.Tag{}, model.ErrTagNotFound). //nolint:exhaustruct
					Once()
				m.
					On("Create", mock.Anything, newTag).
					Return(model.TagID(""), errors.New("test")).
					Once()
				return m
			},
			articleTagRepo: func(t mockConstructorTestingTNewMockArticleTagRepository) *MockArticleTagRepository {
				m := NewMockArticleTagRepository(t)
				m.
					On("GetTagIDs", mock.Anything, article.ID).
					Return(existedArticleTagIDs, nil).
					Once()
				return m
			},
			expectedErr: errors.New("within transaction: get or create tags: create tag(slug=new-tag): test"),
		},
		{
			name:    "create article tag fail",
			ctx:     context.Background(),
			article: article,
			tags:    tags,
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Update", mock.Anything, article).
					Return(nil).
					Once()
				return m
			},
			tagRepo: func(t mockConstructorTestingTNewMockTagRepository) *MockTagRepository {
				m := NewMockTagRepository(t)
				m.
					On("GetBySlug", mock.Anything, newTag.Slug).
					Return(model.Tag{}, model.ErrTagNotFound). //nolint:exhaustruct
					Once()
				m.
					On("Create", mock.Anything, newTag).
					Return(model.TagID("3"), nil).
					Once()
				m.
					On("GetBySlug", mock.Anything, firstExistedTag.Slug).
					Return(firstExistedTag, nil).
					Once()
				return m
			},
			articleTagRepo: func(t mockConstructorTestingTNewMockArticleTagRepository) *MockArticleTagRepository {
				m := NewMockArticleTagRepository(t)
				m.
					On("GetTagIDs", mock.Anything, article.ID).
					Return(existedArticleTagIDs, nil).
					Once()
				m.
					On("Create", mock.Anything, article.ID, model.TagID("3")).
					Return(errors.New("test")).
					Once()
				return m
			},
			expectedErr: errors.New("within transaction: create article tag(tag_id=3, article_id=12): test"),
		},
		{
			name:    "delete article tag fail",
			ctx:     context.Background(),
			article: article,
			tags:    tags,
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Update", mock.Anything, article).
					Return(nil).
					Once()
				return m
			},
			tagRepo: func(t mockConstructorTestingTNewMockTagRepository) *MockTagRepository {
				m := NewMockTagRepository(t)
				m.
					On("GetBySlug", mock.Anything, newTag.Slug).
					Return(model.Tag{}, model.ErrTagNotFound). //nolint:exhaustruct
					Once()
				m.
					On("Create", mock.Anything, newTag).
					Return(model.TagID("3"), nil).
					Once()
				m.
					On("GetBySlug", mock.Anything, firstExistedTag.Slug).
					Return(firstExistedTag, nil).
					Once()
				return m
			},
			articleTagRepo: func(t mockConstructorTestingTNewMockArticleTagRepository) *MockArticleTagRepository {
				m := NewMockArticleTagRepository(t)
				m.
					On("GetTagIDs", mock.Anything, article.ID).
					Return(existedArticleTagIDs, nil).
					Once()
				m.
					On("Create", mock.Anything, article.ID, model.TagID("3")).
					Return(nil).
					Once()
				m.
					On("Delete", mock.Anything, article.ID, secondExistedTag.ID).
					Return(errors.New("test")).
					Once()
				return m
			},
			expectedErr: errors.New("within transaction: delete article tag(tag_id=2, article_id=12): test"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			srv := New(
				transactor.NewDummyTransactor(),
				test.articleRepo(t),
				test.tagRepo(t),
				test.articleTagRepo(t),
			)

			// Action
			err := srv.Update(test.ctx, test.article, test.tags)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}

func TestService_Delete(t *testing.T) {
	testCases := []struct {
		name        string
		ctx         context.Context
		id          model.ArticleID
		articleRepo func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository
		expectedErr error
	}{
		{
			name: "success",
			ctx:  context.Background(),
			id:   model.ArticleID("12"),
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Delete", mock.Anything, model.ArticleID("12")).
					Return(nil).
					Once()
				return m
			},
			expectedErr: nil,
		},
		{
			name: "fail",
			ctx:  context.Background(),
			id:   model.ArticleID("12"),
			articleRepo: func(t mockConstructorTestingTNewMockArticleRepository) *MockArticleRepository {
				m := NewMockArticleRepository(t)
				m.
					On("Delete", mock.Anything, model.ArticleID("12")).
					Return(errors.New("test")).
					Once()
				return m
			},
			expectedErr: errors.New("delete article: test"),
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			srv := New(
				transactor.NewDummyTransactor(),
				test.articleRepo(t),
				NewMockTagRepository(t),
				NewMockArticleTagRepository(t),
			)

			// Action
			err := srv.Delete(test.ctx, test.id)

			// Assert
			asserts.ErrorsEqual(t, test.expectedErr, err)
		})
	}
}
