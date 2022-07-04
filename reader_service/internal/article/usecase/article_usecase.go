package usecase

import (
	"context"
	"time"

	"github.com/radyatamaa/go-cqrs-microservices/pkg/helper"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
	"github.com/radyatamaa/go-cqrs-microservices/reader_service/internal/domain"
)

type articleUseCase struct {
	zapLogger              zaplogger.Logger
	contextTimeout         time.Duration
	mongoArticleRepository domain.MongoArticleRepository
	redisArticleRepository domain.RedisArticleRepository
}

func NewArticleUseCase(timeout time.Duration,
	mongoArticleRepository domain.MongoArticleRepository,
	redisArticleRepository domain.RedisArticleRepository,
	zapLogger zaplogger.Logger) domain.ArticleUseCase {
	return &articleUseCase{
		contextTimeout:         timeout,
		zapLogger:              zapLogger,
		mongoArticleRepository: mongoArticleRepository,
		redisArticleRepository: redisArticleRepository,
	}
}

func (a articleUseCase) CreateArticle(c context.Context, command domain.CreatedArticleCommand) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	insert, err := a.mongoArticleRepository.Create(ctx, command.ToArticle())
	if err != nil {
		a.zapLogger.SetMessageLog(err)
		return err
	}

	a.redisArticleRepository.Put(ctx, helper.IntToString(insert.ID), insert)

	return nil
}

func (a articleUseCase) SearchArticle(c context.Context, query domain.SearchArticleQuery) (*domain.ArticlesList, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	return a.mongoArticleRepository.Search(ctx, query.Text, query.Author, query.Pagination)
}
