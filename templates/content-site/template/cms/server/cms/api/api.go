package api

import (
	"cms/config"
	"cms/proto"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"path"
	"strings"
	"time"
)

var _ proto.CMSServer = (*api)(nil)

type api struct {
	proto.UnimplementedCMSServer
	db *sql.DB
}

func (a *api) ArticleList(ctx context.Context, request *proto.ArticleListRequest) (*proto.ArticleListResponse, error) {
	resp := &proto.ArticleListResponse{}
	var err error
	resp.Total, err = a.totalArticles(ctx, "SELECT COUNT(*) FROM article")
	if nil != err {
		return nil, fmt.Errorf("total: %s", err)
	}
	resp.Articles, err = a.articleList(
		ctx,
		`
		SELECT 
			a.title, a.slug, a.short_body, ct.name, ct.slug, a.published_at 
		FROM article a
		JOIN category ct ON ct.id = a.category_id 
		ORDER BY published_at DESC
		LIMIT $1 
		OFFSET $2`,
		request.Limit,
		request.Offset,
	)
	if nil != err {
		return nil, fmt.Errorf("list: %s", err)
	}
	return resp, nil
}

func (a *api) ArticleListForCategory(ctx context.Context, request *proto.ArticleListForCategoryRequest) (*proto.ArticleListForCategoryResponse, error) {
	resp := &proto.ArticleListForCategoryResponse{}
	var err error
	resp.Total, err = a.totalArticles(ctx,
		`
		SELECT COUNT(*) 
		FROM article a 
		JOIN category c ON c.id = a.category_id 
		WHERE c.slug = $1`,
		request.CategorySlug,
	)
	if nil != err {
		return nil, fmt.Errorf("total: %s", err)
	}
	resp.Articles, err = a.articleList(
		ctx,
		`SELECT 
				a.title, a.slug, a.short_body, ct.name, ct.slug, a.published_at 
			FROM article a
			JOIN category ct ON ct.id = a.category_id 
			WHERE ct.slug = $1
			ORDER BY published_at DESC
			LIMIT $2 
			OFFSET $3`,
		request.CategorySlug,
		request.Limit,
		request.Offset,
	)
	if nil != err {
		return nil, fmt.Errorf("list: %s", err)
	}
	resp.CategoryName, err = a.getStringValue(ctx,
		`SELECT name FROM category WHERE slug = $1`,
		request.CategorySlug,
	)
	if nil != err {
		return nil, fmt.Errorf("category: %s", err)
	}
	return resp, nil
}

func (a *api) ArticleListForTag(ctx context.Context, request *proto.ArticleListForTagRequest) (*proto.ArticleListForTagResponse, error) {
	resp := &proto.ArticleListForTagResponse{}
	var err error
	resp.Total, err = a.totalArticles(ctx,
		`
		SELECT COUNT(DISTINCT a.id) 
		FROM article a
		JOIN article_tag a_t ON a_t.article_id = a.id
		JOIN tag t ON t.id = a_t.tag_id 
		WHERE t.slug = $1`,
		request.TagSlug,
	)
	if nil != err {
		return nil, fmt.Errorf("total: %s", err)
	}
	resp.Articles, err = a.articleList(
		ctx,
		`
		SELECT
			a.title, a.slug, a.short_body, ct.name, ct.slug, a.published_at 
		FROM (
			SELECT a.* FROM article a
			JOIN article_tag at ON at.article_id = a.id
			JOIN tag t ON t.id = at.tag_id 
			WHERE t.slug = $1
			GROUP BY a.id
		) a 
		JOIN category ct ON ct.id = a.category_id
		ORDER BY published_at DESC
		LIMIT $2
		OFFSET $3`,
		request.TagSlug,
		request.Limit,
		request.Offset,
	)
	if nil != err {
		return nil, fmt.Errorf("list: %s", err)
	}
	resp.TagName, err = a.getStringValue(ctx,
		`SELECT name FROM tag WHERE slug = $1`,
		request.TagSlug,
	)
	if nil != err {
		return nil, fmt.Errorf("tag: %s", err)
	}
	return resp, nil
}

func (a *api) ArticleListForQuery(ctx context.Context, request *proto.ArticleListForQueryRequest) (*proto.ArticleListResponse, error) {
	resp := &proto.ArticleListResponse{}
	var err error
	resp.Total, err = a.totalArticles(ctx,
		`
		SELECT 
			count(*) 
		FROM article a 
		WHERE 
			lower(a.title) LIKE '%' || $1 || '%' OR
			lower(a.short_body) LIKE '%' || $1 || '%' OR
			lower(a.body) LIKE '%' || $1 || '%';`,
		request.Query,
	)
	if nil != err {
		return nil, fmt.Errorf("total: %s", err)
	}
	resp.Articles, err = a.articleList(
		ctx,
		`
		SELECT 
			a.title, a.slug, a.short_body, ct.name, ct.slug, a.published_at 
		FROM article a
		JOIN category ct ON ct.id = a.category_id
		WHERE 
			lower(a.title) LIKE '%' || $1 || '%' OR
			lower(a.short_body) LIKE '%' || $1 || '%' OR
			lower(a.body) LIKE '%' || $1 || '%'
		ORDER BY published_at DESC
		LIMIT $2 
		OFFSET $3`,
		request.Query,
		request.Limit,
		request.Offset,
	)
	if nil != err {
		return nil, fmt.Errorf("list: %s", err)
	}
	return resp, nil
}

func (a *api) Article(ctx context.Context, request *proto.ArticleRequest) (*proto.ArticleResponse, error) {
	resp := &proto.ArticleResponse{}
	article := &proto.Article{
		Category: &proto.Category{},
	}
	var publishedAt time.Time
	err := a.db.QueryRowContext(
		ctx,
		`
		SELECT a.title, a.slug, a.short_body, a.body, ct.name, ct.slug, a.published_at
		FROM article a
		JOIN category ct ON ct.id = a.category_id
		WHERE a.slug = $1`,
		request.Slug,
	).Scan(&article.Title, &article.Slug, &article.ShortContent, &article.Content, &article.Category.Name, &article.Category.Slug, &publishedAt)
	if nil != err {
		if errors.Is(err, sql.ErrNoRows) {
			resp.Response = &proto.ArticleResponse_NotFound{
				NotFound: &proto.NotFound{},
			}
			return resp, nil
		}
		return nil, fmt.Errorf("select article: %s", err)
	}
	article.PublishedAt = timestamppb.New(publishedAt)
	tagRows, err := a.db.QueryContext(
		ctx,
		`
			SELECT
				t.name, t.slug 
			FROM article a
			JOIN article_tag at ON at.article_id = a.id
			JOIN tag t ON t.id = at.tag_id
			WHERE a.slug = $1`,
		request.Slug,
	)
	if nil != err {
		return nil, fmt.Errorf("select tags: %s", err)
	}
    defer tagRows.Close()
	for tagRows.Next() {
		tag := &proto.Tag{}
		err := tagRows.Scan(&tag.Name, &tag.Slug)
		if nil != err {
			return nil, fmt.Errorf("scan tag row: %s", err)
		}
		article.Tags = append(article.Tags, tag)
	}
	resp.Response = &proto.ArticleResponse_Article{
		Article: article,
	}
	return resp, nil
}

func (a *api) Asset(ctx context.Context, request *proto.AssetRequest) (*proto.AssetResponse, error) {
	resp := &proto.AssetResponse{}
	filePath := strings.ToLower(request.Path)
	filePath = strings.Replace(filePath, "\\", "/", -1)
	if !strings.HasPrefix(filePath, "/") {
		filePath = "/" + filePath
	}
	filePath = path.Clean(filePath)
	filePath = path.Join(config.Upload.Path, filePath)
	_, err := os.Stat(filePath)
	if nil != err {
		if errors.Is(err, os.ErrNotExist) {
			resp.Response = &proto.AssetResponse_NotFound{
				NotFound: &proto.NotFound{},
			}
			return resp, nil
		}
		return resp, fmt.Errorf("stats: %s", err)
	}
	content, err := os.ReadFile(filePath)
	if nil != err {
		return resp, fmt.Errorf("get content: %s", err)
	}
	resp.Response = &proto.AssetResponse_File{
		File: content,
	}
	return resp, nil
}

func (a *api) totalArticles(ctx context.Context, query string, args ...interface{}) (int64, error) {
	var count int64
	err := a.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if nil != err {
		return 0, fmt.Errorf("select count: %s", err)
	}
	return count, nil
}

func (a *api) getStringValue(ctx context.Context, query string, args ...interface{}) (string, error) {
	var result string
	err := a.db.QueryRowContext(ctx, query, args...).Scan(&result)
	if nil != err {
		return "", fmt.Errorf("select string: %s", err)
	}
	return result, nil
}

func (a *api) articleList(ctx context.Context, query string, args ...interface{}) ([]*proto.ArticleListItem, error) {
	rows, err := a.db.QueryContext(ctx, query, args...)
	if nil != err {
		return nil, fmt.Errorf("query: %s", err)
	}
    defer rows.Close()
	var res []*proto.ArticleListItem
	for rows.Next() {
		var publishedAt time.Time
		data := &proto.ArticleListItem{
			Category: &proto.Category{},
		}
		err := rows.Scan(&data.Title, &data.Slug, &data.Content, &data.Category.Name, &data.Category.Slug, &publishedAt)
		if nil != err {
			return nil, fmt.Errorf("scan row: %s", err)
		}
		data.PublishedAt = timestamppb.New(publishedAt)
		res = append(res, data)
	}
	return res, nil
}

func NewAPI(db *sql.DB) proto.CMSServer {
	return &api{
		db: db,
	}
}
