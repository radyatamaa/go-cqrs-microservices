package kafka

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
	"github.com/radyatamaa/go-cqrs-microservices/write_service/config"
	"github.com/radyatamaa/go-cqrs-microservices/write_service/internal/domain"
	"github.com/segmentio/kafka-go"
)

const (
	retryAttempts = 3
	retryDelay    = 300 * time.Millisecond
	PoolSize      = 30
)

var (
	retryOptions = []retry.Option{retry.Attempts(retryAttempts), retry.Delay(retryDelay), retry.DelayType(retry.BackOffDelay)}
)

type articleConsumer struct {
	zapLogger zaplogger.Logger
	useCase   domain.ArticleUseCase
	cfg       *config.Config
}

func NewArticleConsumer(useCase domain.ArticleUseCase, cfg *config.Config, zapLogger zaplogger.Logger) *articleConsumer {
	return &articleConsumer{
		zapLogger: zapLogger,
		useCase:   useCase,
		cfg:       cfg,
	}
}

func (s *articleConsumer) ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m, err := r.FetchMessage(ctx)
		if err != nil {
			s.zapLogger.Warnf("workerID: %v, err: %v", workerID, err)
			continue
		}

		s.logProcessMessage(m, workerID)

		switch m.Topic {
		case s.cfg.KafkaTopics.ArticleCreate.TopicName:
			s.processCreateArticle(ctx, r, m)
		}
	}
}

func (s *articleConsumer) processCreateArticle(ctx context.Context, r *kafka.Reader, m kafka.Message) {

	var command domain.CreateArticleCommand
	if err := json.Unmarshal(m.Value, &command); err != nil {
		s.zapLogger.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	if err := retry.Do(func() error {
		return s.useCase.CreateArticle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.zapLogger.WarnMsg("ArticleUseCase.CreateArticle", err)
		return
	}

	s.commitMessage(ctx, r, m)
}

func (s *articleConsumer) logProcessMessage(m kafka.Message, workerID int) {
	s.zapLogger.KafkaProcessMessage(m.Topic, m.Partition, string(m.Value), workerID, m.Offset, m.Time)
}

func (s *articleConsumer) commitMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.zapLogger.KafkaLogCommittedMessage(m.Topic, m.Partition, m.Offset)
	if err := r.CommitMessages(ctx, m); err != nil {
		s.zapLogger.WarnMsg("commitMessage", err)
	}
}

func (s *articleConsumer) commitErrMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.zapLogger.KafkaLogCommittedMessage(m.Topic, m.Partition, m.Offset)
	if err := r.CommitMessages(ctx, m); err != nil {
		s.zapLogger.WarnMsg("commitMessage", err)
	}
}
