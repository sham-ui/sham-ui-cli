package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"site/core/handler"
	"site/proto"
	"strings"
)

type detailHandler struct {
	cmsClient proto.CMSClient
}

type detailHandlerData struct {
	slug string
}

func (h *detailHandler) ExtractData(ctx *handler.Context) (interface{}, error) {
	slug, ok := mux.Vars(ctx.Request)["slug"]
	if ok {
		slug = strings.TrimSpace(slug)
		slug = strings.ToLower(slug)
	}
	return &detailHandlerData{
		slug: slug,
	}, nil
}

func (h *detailHandler) Validate(_ *handler.Context, data interface{}) (*handler.Validation, error) {
	validation := handler.NewValidation()
	params := data.(*detailHandlerData)
	if "" == params.slug {
		validation.AddError("slug must be not empty")
	}
	return validation, nil
}

type articleData struct {
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	Category     category `json:"category"`
	Tags         []tag    `json:"tags"`
	ShortContent string   `json:"shortContent"`
	Content      string   `json:"content"`
	CreatedAt    string   `json:"createdAt"`
}

type tag struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (h *detailHandler) Process(ctx *handler.Context, data interface{}) (interface{}, error) {
	params := data.(*detailHandlerData)
	resp, err := h.cmsClient.Article(ctx.Request.Context(), &proto.ArticleRequest{
		Slug: params.slug,
	})
	if nil != err {
		return nil, fmt.Errorf("cms: %s", err)
	}
	if _, ok := resp.Response.(*proto.ArticleResponse_NotFound); ok {
		ctx.RespondWithError(http.StatusNotFound)
		return nil, nil
	}
	article := resp.GetArticle()
	tags := make([]tag, len(article.Tags))
	for i, item := range article.Tags {
		tags[i] = tag{
			Name: item.Name,
			Slug: item.Slug,
		}
	}
	return articleData{
		Title: article.Title,
		Slug:  article.Slug,
		Category: category{
			Name: article.Category.Name,
			Slug: article.Category.Slug,
		},
		Tags:         tags,
		ShortContent: article.ShortContent,
		Content:      article.Content,
		CreatedAt:    article.PublishedAt.AsTime().Format("2006-01-02 15:04:05.999999999 -0700 MST"),
	}, nil
}

func NewDetailHandler(cmsClient proto.CMSClient) http.HandlerFunc {
	return handler.Create(&detailHandler{
		cmsClient: cmsClient,
	})
}
