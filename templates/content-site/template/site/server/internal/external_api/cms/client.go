package cms

import (
	"context"
	"fmt"
	"net"
	"site/pkg/net_addr"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc/credentials/insecure"

	"site/config"
	"site/internal/external_api/cms/proto"
	"site/internal/model"

	"google.golang.org/grpc"
)

type client struct {
	pb cmsClient
}

func (c *client) Asset(ctx context.Context, path string) (*model.Asset, error) {
	resp, err := c.pb.Asset(ctx, &proto.AssetRequest{Path: path})
	if err != nil {
		return nil, NewGRPCError("Asset", err)
	}
	if resp.GetNotFound() != nil {
		return nil, model.NewAssetNotFoundError(path)
	}
	return &model.Asset{
		Path:    path,
		Content: resp.GetFile(),
	}, nil
}

func (c *client) Article(ctx context.Context, slug string) (*model.Article, error) {
	resp, err := c.pb.Article(ctx, &proto.ArticleRequest{Slug: slug})
	if err != nil {
		return nil, NewGRPCError("Article", err)
	}
	if resp.GetNotFound() != nil {
		return nil, model.NewArticleNotFoundError(slug)
	}

	article := resp.GetArticle()

	tags := make([]model.Tag, len(article.GetTags()))
	for i, tag := range article.GetTags() {
		tags[i] = model.Tag{
			Name: tag.GetName(),
			Slug: tag.GetSlug(),
		}
	}

	return &model.Article{
		ShortArticle: model.ShortArticle{
			Title: article.GetTitle(),
			Slug:  article.GetSlug(),
			Category: model.Category{
				Name: article.GetCategory().Name,
				Slug: article.GetCategory().Slug,
			},
			ShortContent: article.GetShortContent(),
			PublishedAt:  article.GetPublishedAt().AsTime(),
		},
		Tags:    tags,
		Content: article.GetContent(),
	}, nil
}

func (c *client) Articles(ctx context.Context, offset, limit int64) (*model.PaginatedArticles, error) {
	resp, err := c.pb.ArticleList(ctx, &proto.ArticleListRequest{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, NewGRPCError("ArticleList", err)
	}
	return &model.PaginatedArticles{
		Articles: articleListItemsToModel(resp.GetArticles()),
		Total:    resp.GetTotal(),
	}, nil
}

func (c *client) ArticleListForQuery(
	ctx context.Context,
	query string,
	offset, limit int64,
) (*model.PaginatedArticles, error) {
	resp, err := c.pb.ArticleListForQuery(ctx, &proto.ArticleListForQueryRequest{
		Offset: offset,
		Limit:  limit,
		Query:  query,
	})
	if err != nil {
		return nil, NewGRPCError("ArticleListForQuery", err)
	}
	return &model.PaginatedArticles{
		Articles: articleListItemsToModel(resp.GetArticles()),
		Total:    resp.GetTotal(),
	}, nil
}

func (c *client) ArticleListForTag(
	ctx context.Context,
	tagSlug string,
	offset, limit int64,
) (*model.PaginatedArticleForTag, error) {
	resp, err := c.pb.ArticleListForTag(ctx, &proto.ArticleListForTagRequest{
		Offset:  offset,
		Limit:   limit,
		TagSlug: tagSlug,
	})
	if err != nil {
		return nil, NewGRPCError("ArticleListForTag", err)
	}
	return &model.PaginatedArticleForTag{
		Articles: articleListItemsToModel(resp.GetArticles()),
		Total:    resp.GetTotal(),
		TagName:  resp.GetTagName(),
	}, nil
}

func (c *client) ArticleListForCategory(
	ctx context.Context,
	categorySlug string,
	offset, limit int64,
) (*model.PaginatedArticleForCategory, error) {
	resp, err := c.pb.ArticleListForCategory(ctx, &proto.ArticleListForCategoryRequest{
		Offset:       offset,
		Limit:        limit,
		CategorySlug: categorySlug,
	})
	if err != nil {
		return nil, NewGRPCError("ArticleListForCategory", err)
	}
	return &model.PaginatedArticleForCategory{
		Articles:     articleListItemsToModel(resp.GetArticles()),
		Total:        resp.GetTotal(),
		CategoryName: resp.GetCategoryName(),
	}, nil
}

func New(
	tracerProvider trace.TracerProvider,
	propagator propagation.TraceContext,
	cfg config.API,
) (*client, error) {
	conn, err := grpc.Dial(
		cfg.Address,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, str string) (net.Conn, error) {
			ctxTimeout, cancel := context.WithTimeout(ctx, cfg.DialTimeout)
			defer cancel()
			var d net.Dialer
			network, addr := net_addr.Resolve(str)
			return d.DialContext(ctxTimeout, network, addr) //nolint:wrapcheck
		}),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(
				otelgrpc.WithTracerProvider(tracerProvider),
				otelgrpc.WithPropagators(propagator),
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("can't connect to %s: %w", cfg.Address, err)
	}
	return &client{
		pb: proto.NewCMSClient(conn),
	}, nil
}
