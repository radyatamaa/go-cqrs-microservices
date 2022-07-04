package grpc

import (
	"context"

	"github.com/radyatamaa/go-cqrs-microservices/pkg/utils"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
	"github.com/radyatamaa/go-cqrs-microservices/reader_service/config"
	"github.com/radyatamaa/go-cqrs-microservices/reader_service/internal/domain"
	readerService "github.com/radyatamaa/go-cqrs-microservices/reader_service/proto/article_reader"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type articleGrpcService struct {
	zapLogger zaplogger.Logger
	useCase   domain.ArticleUseCase
	cfg       *config.Config
}

func NewArticleGrpcService(useCase domain.ArticleUseCase, cfg *config.Config, zapLogger zaplogger.Logger) *articleGrpcService {
	return &articleGrpcService{
		zapLogger: zapLogger,
		useCase:   useCase,
		cfg:       cfg,
	}
}

func (s *articleGrpcService) SearchArticle(ctx context.Context, req *readerService.SearchReq) (*readerService.SearchRes, error) {
	pq := utils.NewPaginationQuery(int(req.GetSize()), int(req.GetPage()))

	query := domain.NewSearchArticleQuery(req.GetSearch(), req.Author, pq)
	articlesList, err := s.useCase.SearchArticle(ctx, query)
	if err != nil {
		s.zapLogger.WarnMsg("ArticleUseCase.SearchArticle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return domain.ArticleListToGrpc(articlesList), nil
}

func (s *articleGrpcService) errResponse(c codes.Code, err error) error {
	return status.Error(c, err.Error())
}
