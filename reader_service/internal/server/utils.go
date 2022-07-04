package server

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/heptiolabs/healthcheck"
	"github.com/pkg/errors"
	kafkaClient "github.com/radyatamaa/go-cqrs-microservices/pkg/kafka"
	"github.com/segmentio/kafka-go"
)

const (
	stackSize = 1 << 10 // 1 KB
)

func (s *server) connectKafkaBrokers(ctx context.Context) error {
	kafkaConn, err := kafkaClient.NewKafkaConn(ctx, s.cfg.Kafka)
	if err != nil {
		return errors.Wrap(err, "kafka.NewKafkaCon")
	}

	s.kafkaConn = kafkaConn

	brokers, err := kafkaConn.Brokers()
	if err != nil {
		return errors.Wrap(err, "kafkaConn.Brokers")
	}

	s.zapLog.Infof("kafka connected to brokers: %+v", brokers)

	return nil
}

func (s *server) initKafkaTopics(ctx context.Context) {
	controller, err := s.kafkaConn.Controller()
	if err != nil {
		s.zapLog.WarnMsg("kafkaConn.Controller", err)
		return
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	s.zapLog.Infof("kafka controller uri: %s", controllerURI)

	conn, err := kafka.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		s.zapLog.WarnMsg("initKafkaTopics.DialContext", err)
		return
	}
	defer conn.Close() // nolint: errcheck

	s.zapLog.Infof("established new kafka controller connection: %s", controllerURI)

	articleCreateTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.ArticleCreate.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.ArticleCreate.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.ArticleCreate.ReplicationFactor,
	}

	articleCreatedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.ArticleCreated.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.ArticleCreated.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.ArticleCreated.ReplicationFactor,
	}

	if err := conn.CreateTopics(
		articleCreateTopic,
		articleCreatedTopic,
	); err != nil {
		s.zapLog.WarnMsg("kafkaConn.CreateTopics", err)
		return
	}

	s.zapLog.Infof("kafka topics created or already exists: %+v", []kafka.TopicConfig{articleCreateTopic, articleCreatedTopic})
}

func (s *server) runHealthCheck(ctx context.Context) {
	health := healthcheck.NewHandler()

	health.AddLivenessCheck(s.cfg.App.ServiceName, healthcheck.AsyncWithContext(ctx, func() error {
		return nil
	}, time.Duration(s.cfg.App.CheckIntervalSeconds)*time.Second))

	health.AddReadinessCheck("redis", healthcheck.AsyncWithContext(ctx, func() error {
		return s.redisClient.Ping(ctx).Err()
	}, time.Duration(s.cfg.App.CheckIntervalSeconds)*time.Second))

	health.AddReadinessCheck("mongo", healthcheck.AsyncWithContext(ctx, func() error {
		return s.mongoClient.Ping(ctx, nil)
	}, time.Duration(s.cfg.App.CheckIntervalSeconds)*time.Second))

	health.AddReadinessCheck("kafka", healthcheck.AsyncWithContext(ctx, func() error {
		_, err := s.kafkaConn.Brokers()
		if err != nil {
			return err
		}
		return nil
	}, time.Duration(s.cfg.App.CheckIntervalSeconds)*time.Second))

	go func() {
		s.zapLog.Infof("Reader microservice Kubernetes probes listening on port: %s", s.cfg.App.Port)
		if err := http.ListenAndServe(s.cfg.App.Port, health); err != nil {
			s.zapLog.WarnMsg("ListenAndServe", err)
		}
	}()
}

func (s *server) getConsumerGroupTopics() []string {
	return []string{
		s.cfg.KafkaTopics.ArticleCreated.TopicName,
	}
}
