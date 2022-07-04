package server

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/interceptors"
	kafkaClient "github.com/radyatamaa/go-cqrs-microservices/pkg/kafka"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/mongodb"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
	"github.com/radyatamaa/go-cqrs-microservices/reader_service/config"
	articleConsumerHandler "github.com/radyatamaa/go-cqrs-microservices/reader_service/internal/article/delivery/kafka"
	"github.com/radyatamaa/go-cqrs-microservices/reader_service/internal/domain"

	"os"
	"os/signal"
	"syscall"
	"time"

	redisClient "github.com/radyatamaa/go-cqrs-microservices/pkg/redis"
	articleRepository "github.com/radyatamaa/go-cqrs-microservices/reader_service/internal/article/repository"
	articlUsecase "github.com/radyatamaa/go-cqrs-microservices/reader_service/internal/article/usecase"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type server struct {
	zapLog         zaplogger.Logger
	cfg            *config.Config
	kafkaConn      *kafka.Conn
	mongoClient    *mongo.Client
	redisClient    redis.UniversalClient
	im             interceptors.InterceptorManager
	articleUsecase domain.ArticleUseCase
}

func NewServer(cfg *config.Config, zapLog zaplogger.Logger) *server {
	return &server{cfg: cfg, zapLog: zapLog}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	s.im = interceptors.NewInterceptorManager(s.zapLog)

	// database initialization
	mongoDBConn, err := mongodb.NewMongoDBConn(ctx, s.cfg.Mongo)
	if err != nil {
		return errors.Wrap(err, "NewMongoDBConn")
	}
	s.mongoClient = mongoDBConn
	defer mongoDBConn.Disconnect(ctx) // nolint: errcheck
	s.zapLog.Infof("Mongo connected: %v", mongoDBConn.NumberSessionsInProgress())

	// cache initialization
	s.redisClient = redisClient.NewUniversalRedisClient(s.cfg.Redis)
	defer s.redisClient.Close() // nolint: errcheck
	s.zapLog.Infof("Redis connected: %+v", s.redisClient.PoolStats())

	kafkaProducer := kafkaClient.NewProducer(s.zapLog, viper.GetStringSlice("kafka.brokers"))
	defer kafkaProducer.Close() // nolint: errcheck

	timeoutContext := time.Duration(s.cfg.App.ExecutionTimeout) * time.Second

	mongoArticleRepo := articleRepository.NewMongoArticleRepository(s.zapLog, s.cfg, s.mongoClient)
	redisArticleRepo := articleRepository.NewRedisRepository(s.zapLog, s.cfg, s.redisClient)

	s.articleUsecase = articlUsecase.NewArticleUseCase(timeoutContext, mongoArticleRepo, redisArticleRepo, s.zapLog)

	kafkaArticleConsumerHandler := articleConsumerHandler.NewArticleConsumer(s.articleUsecase, s.cfg, s.zapLog)

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

	s.runHealthCheck(ctx)

	<-ctx.Done()
	grpcServer.GracefulStop()

	return nil
}
