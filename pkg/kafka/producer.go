package kafka

import (
	"context"

	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishMessage(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type producer struct {
	brokers []string
	w       *kafka.Writer
	log     zaplogger.Logger
}

// NewProducer create new kafka producer
func NewProducer(log zaplogger.Logger, brokers []string) *producer {
	return &producer{log: log, brokers: brokers, w: NewWriter(brokers, kafka.LoggerFunc(log.Errorf))}
}

func (p *producer) PublishMessage(ctx context.Context, msgs ...kafka.Message) error {
	return p.w.WriteMessages(ctx, msgs...)
}

func (p *producer) Close() error {
	return p.w.Close()
}
