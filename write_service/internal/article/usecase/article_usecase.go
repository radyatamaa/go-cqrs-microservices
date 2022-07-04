package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
	"github.com/radyatamaa/go-cqrs-microservices/write_service/internal/domain"
)

type articleUseCase struct {
	zapLogger                  zaplogger.Logger
	contextTimeout             time.Duration
	pgArticleRepository        domain.PgArticleRepository
	messagingArticleRepository domain.MessagingArticleRepository
}

func NewArticleUseCase(timeout time.Duration,
	pgArticleRepository domain.PgArticleRepository,
	messagingArticleRepository domain.MessagingArticleRepository,
	zapLogger zaplogger.Logger) domain.ArticleUseCase {
	return &articleUseCase{
		contextTimeout:             timeout,
		zapLogger:                  zapLogger,
		pgArticleRepository:        pgArticleRepository,
		messagingArticleRepository: messagingArticleRepository,
	}
}

func (a articleUseCase) CreateArticle(c context.Context, command domain.CreateArticleCommand) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	insert, err := a.pgArticleRepository.Store(ctx, command.ToArticle())
	if err != nil {
		a.zapLogger.SetMessageLog(err)
		return err
	}

	createdCommand := domain.CreatedArticleCommand{
		ID:        insert.ID,
		Author:    insert.Author,
		Title:     insert.Title,
		Body:      insert.Body,
		CreatedAt: insert.CreatedAt,
		UpdatedAt: insert.UpdatedAt,
	}

	msg, err := json.Marshal(createdCommand)
	if err != nil {
		a.zapLogger.SetMessageLog(err)
		return err
	}

	err = a.messagingArticleRepository.PushMessageInsertArticle(ctx, msg)
	if err != nil {
		a.zapLogger.SetMessageLog(err)
		return err
	}

	return nil
}
