package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/radyatamaa/go-cqrs-microservices/api_gateway_service/internal/domain"
	kafkaClient "github.com/radyatamaa/go-cqrs-microservices/pkg/kafka"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
	"github.com/segmentio/kafka-go"
)

type commandArticleRepository struct {
	zapLogger       zaplogger.Logger
	producer        kafkaClient.Producer
	confKafkaTopics domain.ConfKafkaTopics
}

func NewCommandArticleRepository(producer kafkaClient.Producer, confKafkaTopics domain.ConfKafkaTopics, zapLogger zaplogger.Logger) domain.CommandArticleRepository {
	return &commandArticleRepository{
		producer:        producer,
		confKafkaTopics: confKafkaTopics,
		zapLogger:       zapLogger,
	}
}

func (m commandArticleRepository) Create(ctx context.Context, command domain.CreateArticleCommand) error {
	msg, err := json.Marshal(command)
	if err != nil {
		return err
	}
	err = m.producer.PublishMessage(ctx, kafka.Message{
		Topic: m.confKafkaTopics.CreateArticle,
		Value: msg,
		Time:  time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	return nil
}
