package repository

import (
	"context"
	"time"

	kafkaClient "github.com/radyatamaa/go-cqrs-microservices/pkg/kafka"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
	"github.com/radyatamaa/go-cqrs-microservices/write_service/config"
	"github.com/radyatamaa/go-cqrs-microservices/write_service/internal/domain"
	"github.com/segmentio/kafka-go"
)

type messagingArticleRepository struct {
	zapLogger zaplogger.Logger
	producer  kafkaClient.Producer
	cfg       *config.Config
}

func NewMessagingArticleRepository(producer kafkaClient.Producer, cfg *config.Config, zapLogger zaplogger.Logger) domain.MessagingArticleRepository {
	return &messagingArticleRepository{
		producer:  producer,
		zapLogger: zapLogger,
		cfg:       cfg,
	}
}

func (m messagingArticleRepository) PushMessageInsertArticle(ctx context.Context, msg []byte) error {
	return m.producer.PublishMessage(ctx, kafka.Message{
		Topic: m.cfg.KafkaTopics.ArticleCreated.TopicName,
		Value: msg,
		Time:  time.Now().UTC(),
	})
}
