package usecase

import (
	"context"
	"time"

	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/radyatamaa/go-cqrs-microservices/api_gateway_service/internal/domain"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
)

type articleUseCase struct {
	zapLogger                zaplogger.Logger
	contextTimeout           time.Duration
	articleCommandRepository domain.CommandArticleRepository
	articleQueriesRepository domain.QueriesArticleRepository
}

func NewArticleUseCase(timeout time.Duration,
	zapLogger zaplogger.Logger,
	articleCommandRepository domain.CommandArticleRepository,
	articleQueriesRepository domain.QueriesArticleRepository) domain.ArticleUseCase {
	return &articleUseCase{
		articleCommandRepository: articleCommandRepository,
		articleQueriesRepository: articleQueriesRepository,
		contextTimeout:           timeout,
		zapLogger:                zapLogger,
	}
}

func (a articleUseCase) CreateArticle(beegoCtx *beegoContext.Context, body domain.CreateArticleRequest) error {
	c, cancel := context.WithTimeout(beegoCtx.Request.Context(), a.contextTimeout)
	defer cancel()

	err := a.articleCommandRepository.Create(c, body.ToCreateArticleCommand())
	if err != nil {
		beegoCtx.Input.SetData("stackTrace", a.zapLogger.SetMessageLog(err))
		return err
	}

	return nil
}

func (a articleUseCase) GetArticles(beegoCtx *beegoContext.Context, page int, size int, search string, author string) (*domain.ArticlePaginationResponse, error) {
	c, cancel := context.WithTimeout(beegoCtx.Request.Context(), a.contextTimeout)
	defer cancel()

	result := new(domain.ArticlePaginationResponse)

	list, err := a.articleQueriesRepository.Search(c, page, size, search, author)
	if err != nil {
		beegoCtx.Input.SetData("stackTrace", a.zapLogger.SetMessageLog(err))
		return nil, err
	}

	result = result.ToArticlePaginationResponse(list)
	for i := range list.Articles {
		result.Articles = append(result.Articles, &domain.ArticleResponse{
			ID:        int(list.Articles[i].ID),
			Author:    list.Articles[i].Author,
			Title:     list.Articles[i].Title,
			Body:      list.Articles[i].Body,
			CreatedAt: list.Articles[i].CreatedAt.AsTime(),
			UpdatedAt: list.Articles[i].UpdatedAt.AsTime(),
		})
	}

	return result, nil
}
