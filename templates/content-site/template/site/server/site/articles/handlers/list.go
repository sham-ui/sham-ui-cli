package handlers

import (
	"fmt"
	"github.com/go-logr/logr"
	"net/http"
	"site/core/handler"
	"site/proto"
	"strconv"
	"strings"
)

type listHandler struct {
	cmsClient proto.CMSClient
}

type listHandlerData struct {
	tag      string
	category string
	query    string
	offset   int64
	limit    int64
}

func (h *listHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	offset, err := strconv.ParseInt(ctx.Request.URL.Query().Get("offset"), 10, 64)
	if nil != err {
		offset = 0
	}
	limit, err := strconv.ParseInt(ctx.Request.URL.Query().Get("limit"), 10, 64)
	if nil != err {
		limit = 20
	}
	return &listHandlerData{
		offset:   offset,
		limit:    limit,
		tag:      ctx.Request.URL.Query().Get("tag"),
		category: ctx.Request.URL.Query().Get("category"),
		query:    strings.TrimSpace(ctx.Request.URL.Query().Get("q")),
	}, nil
}

func (h *listHandler) Validate(_ *handler.Context, data interface{}) (*handler.Validation, error) {
	validation := handler.NewValidation()
	params := data.(*listHandlerData)
	if params.limit <= 0 {
		validation.AddError("limit must be > 0")
	}
	if params.offset < 0 {
		validation.AddError("offset must be >= 0")
	}
	return validation, nil
}

type listArticlesData struct {
	Title     string   `json:"title"`
	Slug      string   `json:"slug"`
	Category  category `json:"category"`
	Content   string   `json:"content"`
	CreatedAt string   `json:"createdAt"`
}

type category struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (h *listHandler) Process(ctx *handler.Context, data interface{}) (interface{}, error) {
	params := data.(*listHandlerData)
	query := ctx.Request.URL.Query()
	if _, ok := query["q"]; ok {
		return h.processQuery(ctx, params)
	}
	if _, ok := query["category"]; ok {
		return h.processCategory(ctx, params)
	}
	if _, ok := query["tag"]; ok {
		return h.processTag(ctx, params)
	}
	return h.processDefault(ctx, params)
}

func (h *listHandler) processQuery(ctx *handler.Context, params *listHandlerData) (interface{}, error) {
	resp, err := h.cmsClient.ArticleListForQuery(ctx.Request.Context(), &proto.ArticleListForQueryRequest{
		Offset: params.offset,
		Limit:  params.limit,
		Query:  params.query,
	})
	if nil != err {
		return nil, fmt.Errorf("cms: %s", err)
	}
	return map[string]interface{}{
		"meta":     h.buildMeta(params.offset, params.limit, resp.Total),
		"articles": h.buildArticles(resp.Articles),
	}, nil
}

func (h *listHandler) processCategory(ctx *handler.Context, params *listHandlerData) (interface{}, error) {
	resp, err := h.cmsClient.ArticleListForCategory(ctx.Request.Context(), &proto.ArticleListForCategoryRequest{
		Offset:       params.offset,
		Limit:        params.limit,
		CategorySlug: params.category,
	})
	if nil != err {
		return nil, fmt.Errorf("cms: %s", err)
	}
	meta := h.buildMeta(params.offset, params.limit, resp.Total)
	meta["category"] = resp.CategoryName
	return map[string]interface{}{
		"meta":     meta,
		"articles": h.buildArticles(resp.Articles),
	}, nil
}

func (h *listHandler) processTag(ctx *handler.Context, params *listHandlerData) (interface{}, error) {
	resp, err := h.cmsClient.ArticleListForTag(ctx.Request.Context(), &proto.ArticleListForTagRequest{
		Offset:  params.offset,
		Limit:   params.limit,
		TagSlug: params.tag,
	})
	if nil != err {
		return nil, fmt.Errorf("cms: %s", err)
	}
	meta := h.buildMeta(params.offset, params.limit, resp.Total)
	meta["tag"] = resp.TagName
	return map[string]interface{}{
		"meta":     meta,
		"articles": h.buildArticles(resp.Articles),
	}, nil
}

func (h *listHandler) processDefault(ctx *handler.Context, params *listHandlerData) (interface{}, error) {
	resp, err := h.cmsClient.ArticleList(ctx.Request.Context(), &proto.ArticleListRequest{
		Offset: params.offset,
		Limit:  params.limit,
	})
	if nil != err {
		return nil, fmt.Errorf("cms: %s", err)
	}
	return map[string]interface{}{
		"meta":     h.buildMeta(params.offset, params.limit, resp.Total),
		"articles": h.buildArticles(resp.Articles),
	}, nil
}

func (h *listHandler) buildArticles(items []*proto.ArticleListItem) []listArticlesData {
	articles := make([]listArticlesData, len(items))
	for i, item := range items {
		articles[i] = listArticlesData{
			Title: item.Title,
			Slug:  item.Slug,
			Category: category{
				Name: item.Category.Name,
				Slug: item.Category.Slug,
			},
			Content:   item.Content,
			CreatedAt: item.PublishedAt.AsTime().Format("2006-01-02 15:04:05.999999999 -0700 MST"),
		}
	}
	return articles
}

func (h *listHandler) buildMeta(offset, limit, total int64) map[string]interface{} {
	return map[string]interface{}{
		"offset": offset,
		"limit":  limit,
		"total":  total,
	}
}

func NewListHandler(logger logr.Logger, cmsClient proto.CMSClient) http.HandlerFunc {
	return handler.Create(logger, &listHandler{
		cmsClient: cmsClient,
	})
}
