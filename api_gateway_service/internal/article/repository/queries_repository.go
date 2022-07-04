package repository

import (
	"context"

	"github.com/radyatamaa/go-cqrs-microservices/api_gateway_service/internal/domain"
	readerService "github.com/radyatamaa/go-cqrs-microservices/api_gateway_service/proto/article_reader"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
)

type queriesArticleRepository struct {
	zapLogger zaplogger.Logger
	rsClient  readerService.ReaderServiceClient
}

func NewQueriesArticleRepository(rsClient readerService.ReaderServiceClient, zapLogger zaplogger.Logger) domain.QueriesArticleRepository {
	return &queriesArticleRepository{
		rsClient:  rsClient,
		zapLogger: zapLogger,
	}
}

func (q queriesArticleRepository) Search(ctx context.Context, page int, size int, search string, author string) (*readerService.SearchRes, error) {
	res, err := q.rsClient.SearchArticle(ctx, &readerService.SearchReq{
		Author: author,
		Search: search,
		Page:   int64(page),
		Size:   int64(size),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
