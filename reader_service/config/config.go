package config

import (
	"fmt"
	"os"

	kafkaClient "github.com/radyatamaa/go-cqrs-microservices/pkg/kafka"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/mongodb"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/redis"
	"github.com/spf13/viper"
)

const (
	GrpcPort       = "GRPC_PORT"
	HttpPort       = "HTTP_PORT"
	ConfigPath     = "CONFIG_PATH"
	KafkaBrokers   = "KAFKA_BROKERS"
	JaegerHostPort = "JAEGER_HOST"
	RedisAddr      = "REDIS_ADDR"
	MongoDbURI     = "MONGO_URI"
	PostgresqlHost = "POSTGRES_HOST"
	PostgresqlPort = "POSTGRES_PORT"
)

type Config struct {
	App              AppConfig
	KafkaTopics      KafkaTopics
	Kafka            *kafkaClient.Config
	Mongo            *mongodb.Config
	Redis            *redis.Config
	MongoCollections MongoCollections
	ServiceSettings  ServiceSettings
	GRPC             GRPC
}

type GRPC struct {
	Port        string
	Development bool
}

type AppConfig struct {
	Port                 string
	ServiceName          string
	ExecutionTimeout     int
	CheckIntervalSeconds int
	LogPath              string
	SlackWebHookUrl      string
}
type MongoCollections struct {
	Articles string
}

type KafkaTopics struct {
	ArticleCreate  kafkaClient.TopicConfig
	ArticleCreated kafkaClient.TopicConfig
}

type ServiceSettings struct {
	RedisArticlePrefixKey string
}

func InitConfig() (*Config, error) {

	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath("./reader_service/config")

	// Set config type
	viper.SetConfigType("json")

	// read env
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		return nil, err
	}

	cfg := &Config{
		App: AppConfig{
			Port:                 viper.GetString("app.port"),
			ServiceName:          viper.GetString("app.serviceName"),
			ExecutionTimeout:     viper.GetInt("app.executionTimeout"),
			CheckIntervalSeconds: viper.GetInt("app.checkIntervalSeconds"),
			LogPath:              viper.GetString("app.logPath"),
			SlackWebHookUrl:      viper.GetString("app.slackWebHookUrl"),
		},
		KafkaTopics: KafkaTopics{
			ArticleCreate: kafkaClient.TopicConfig{
				TopicName:         viper.GetString("kafkaTopics.articleCreate.topicName"),
				Partitions:        viper.GetInt("kafkaTopics.articleCreate.partitions"),
				ReplicationFactor: viper.GetInt("kafkaTopics.articleCreate.replicationFactor"),
			},
			ArticleCreated: kafkaClient.TopicConfig{
				TopicName:         viper.GetString("kafkaTopics.articleCreated.topicName"),
				Partitions:        viper.GetInt("kafkaTopics.articleCreated.partitions"),
				ReplicationFactor: viper.GetInt("kafkaTopics.articleCreated.replicationFactor"),
			},
		},
		Kafka: &kafkaClient.Config{
			Brokers:    viper.GetStringSlice("kafka.brokers"),
			GroupID:    viper.GetString("kafka.groupID"),
			InitTopics: viper.GetBool("kafka.initTopics"),
		},
		Mongo: &mongodb.Config{
			URI:      viper.GetString("mongo.uri"),
			User:     viper.GetString("mongo.user"),
			Password: viper.GetString("mongo.password"),
			Db:       viper.GetString("mongo.db"),
		},
		Redis: &redis.Config{
			Addr:     viper.GetString("redis.addr"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.dB"),
			PoolSize: viper.GetInt("redis.poolSize"),
		},
		MongoCollections: MongoCollections{
			Articles: viper.GetString("mongoCollections.articles"),
		},
		ServiceSettings: ServiceSettings{
			RedisArticlePrefixKey: viper.GetString("serviceSettings.redisArticlePrefixKey"),
		},
		GRPC: GRPC{
			Port:        viper.GetString("grpc.port"),
			Development: viper.GetBool("grpc.development"),
		},
	}

	grpcPort := os.Getenv(GrpcPort)
	if grpcPort != "" {
		cfg.GRPC.Port = grpcPort
	}

	mongoURI := os.Getenv(MongoDbURI)
	if mongoURI != "" {
		cfg.Mongo.URI = mongoURI
	}
	redisAddr := os.Getenv(RedisAddr)
	if redisAddr != "" {
		cfg.Redis.Addr = redisAddr
	}

	kafkaBrokers := os.Getenv(KafkaBrokers)
	if kafkaBrokers != "" {
		cfg.Kafka.Brokers = []string{kafkaBrokers}
	}

	return cfg, nil
}
