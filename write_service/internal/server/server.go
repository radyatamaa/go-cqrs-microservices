package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/database"
	kafkaClient "github.com/radyatamaa/go-cqrs-microservices/pkg/kafka"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
	"github.com/radyatamaa/go-cqrs-microservices/write_service/config"
	articleConsumerHandler "github.com/radyatamaa/go-cqrs-microservices/write_service/internal/article/delivery/kafka"
	articleRepository "github.com/radyatamaa/go-cqrs-microservices/write_service/internal/article/repository"
	articlUsecase "github.com/radyatamaa/go-cqrs-microservices/write_service/internal/article/usecase"
	"github.com/radyatamaa/go-cqrs-microservices/write_service/internal/domain"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type server struct {
	zapLog    zaplogger.Logger
	cfg       *config.Config
	db        *gorm.DB
	kafkaConn *kafka.Conn
}

func NewServer(cfg *config.Config, zapLog zaplogger.Logger) *server {
	return &server{cfg: cfg, zapLog: zapLog}
}

func (s *server) Run() error {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// database initialization
	conn, err := database.New(
		func(config *database.Config) {
			config.Driver = s.cfg.Database.Driver
			config.Host = s.cfg.Database.Host
			config.Port = s.cfg.Database.Port
			config.Name = s.cfg.Database.Name
			config.Username = s.cfg.Database.Username
			config.Password = s.cfg.Database.Password
			config.Options = s.cfg.Database.Options
			config.MaxOpenConnection = s.cfg.Database.MaxOpenConnection
			config.MaxIdleConnection = s.cfg.Database.MaxIdleConnection
			config.MaxLifeTimeConnection = s.cfg.Database.MaxLifeTimeConnection
			config.MaxIdleTimeConnection = s.cfg.Database.MaxIdleTimeConnection
		},
	)
	if err != nil {
		panic(err)
	}

	s.db = conn[s.cfg.Database.Name]

	// db auto migrate dev environment
	if err := s.db.AutoMigrate(
		&domain.Article{}); err != nil {
		panic(err)
	}

	kafkaProducer := kafkaClient.NewProducer(s.zapLog, s.cfg.Kafka.Brokers)
	defer kafkaProducer.Close() // nolint: errcheck

	timeoutContext := time.Duration(s.cfg.App.ExecutionTimeout) * time.Second

	messagingArticleRepo := articleRepository.NewMessagingArticleRepository(kafkaProducer, s.cfg, s.zapLog)
	pgArticleRepo := articleRepository.NewPgArticleRepository(s.db, s.zapLog)

	articleUcase := articlUsecase.NewArticleUseCase(timeoutContext, pgArticleRepo, messagingArticleRepo, s.zapLog)

	kafkaArticleConsumerHandler := articleConsumerHandler.NewArticleConsumer(articleUcase, s.cfg, s.zapLog)

	s.zapLog.Infof("Starting Reader Kafka consumers")
	consumerGroup := kafkaClient.NewConsumerGroup(s.cfg.Kafka.Brokers, s.cfg.Kafka.GroupID, s.zapLog)
	go consumerGroup.ConsumeTopic(ctx, s.getConsumerGroupTopics(), articleConsumerHandler.PoolSize, kafkaArticleConsumerHandler.ProcessMessages)

	if err := s.connectKafkaBrokers(ctx); err != nil {
		return errors.Wrap(err, "s.connectKafkaBrokers")
	}
	defer s.kafkaConn.Close() // nolint: errcheck

	closeGrpcServer, grpcServer, err := s.newReaderGrpcServer()
	if err != nil {
		return errors.Wrap(err, "NewScmGrpcServer")
	}
	defer closeGrpcServer() // nolint: errcheck

	if s.cfg.Kafka.InitTopics {
		s.initKafkaTopics(ctx)
	}

	s.runHealthCheck(ctx)

	<-ctx.Done()
	grpcServer.GracefulStop()

	return nil
}
