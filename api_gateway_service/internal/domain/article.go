package domain

import (
	"context"

	beegoContext "github.com/beego/beego/v2/server/web/context"
	readerService "github.com/radyatamaa/go-cqrs-microservices/api_gateway_service/proto/article_reader"
)

type CreateArticleCommand struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type ConfKafkaTopics struct {
	CreateArticle string
}

// ArticleUseCase UseCase Interface
type ArticleUseCase interface {
	CreateArticle(beegoCtx *beegoContext.Context, body CreateArticleRequest) error
	GetArticles(beegoCtx *beegoContext.Context, page int, size int, search string, author string) (*ArticlePaginationResponse, error)
}

// CommandArticleRepository Repository Interface
type CommandArticleRepository interface {
	Create(ctx context.Context, command CreateArticleCommand) error
}

// QueriesArticleRepository Repository Interface
type QueriesArticleRepository interface {
	Search(ctx context.Context, page int, size int, search string, author string) (*readerService.SearchRes, error)
}

// Mapper
func (a ArticlePaginationResponse) ToArticlePaginationResponse(r *readerService.SearchRes) *ArticlePaginationResponse {
	result := &ArticlePaginationResponse{
		TotalCount: r.TotalCount,
		TotalPages: r.TotalPages,
		Page:       r.Page,
		Size:       r.Size,
		HasMore:    r.HasMore,
	}
	return result
}
