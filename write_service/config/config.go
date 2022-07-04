package config

import (
	"fmt"
	"os"

	"github.com/radyatamaa/go-cqrs-microservices/pkg/database"
	kafkaClient "github.com/radyatamaa/go-cqrs-microservices/pkg/kafka"
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
	App         AppConfig
	Database    database.Config
	KafkaTopics KafkaTopics
	Kafka       *kafkaClient.Config
	GRPC        GRPC
}

type AppConfig struct {
	Port                 string
	ServiceName          string
	ExecutionTimeout     int
	CheckIntervalSeconds int
	LogPath              string
	SlackWebHookUrl      string
}
type KafkaTopics struct {
	ArticleCreate  kafkaClient.TopicConfig
	ArticleCreated kafkaClient.TopicConfig
}

type GRPC struct {
	Port        string
	Development bool
}

func InitConfig() (*Config, error) {

	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath("./write_service/config")

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
		Database: database.Config{
			Driver:                viper.GetString("database.driver"),
			Host:                  viper.GetString("database.host"),
			Port:                  viper.GetString("database.port"),
			Name:                  viper.GetString("database.name"),
			Username:              viper.GetString("database.username"),
			Password:              viper.GetString("database.password"),
			Options:               viper.GetString("database.options"),
			MaxOpenConnection:     viper.GetInt("database.maxOpenConnections"),
			MaxIdleConnection:     viper.GetInt("database.maxIdleConnections"),
			MaxLifeTimeConnection: viper.GetInt("database.maxLifetime"),
			MaxIdleTimeConnection: viper.GetInt("database.maxIdleTime"),
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
		GRPC: GRPC{
			Port:        viper.GetString("grpc.port"),
			Development: viper.GetBool("grpc.development"),
		},
	}

	grpcPort := os.Getenv(GrpcPort)
	if grpcPort != "" {
		cfg.GRPC.Port = grpcPort
	}

	postgresHost := os.Getenv(PostgresqlHost)
	if postgresHost != "" {
		cfg.Database.Host = postgresHost
	}
	postgresPort := os.Getenv(PostgresqlPort)
	if postgresPort != "" {
		cfg.Database.Port = postgresPort
	}

	kafkaBrokers := os.Getenv(KafkaBrokers)
	if kafkaBrokers != "" {
		cfg.Kafka.Brokers = []string{kafkaBrokers}
	}

	return cfg, nil
}
