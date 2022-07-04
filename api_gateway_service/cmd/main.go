package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/filter/cors"
	"github.com/beego/i18n"
	"github.com/radyatamaa/go-cqrs-microservices/api_gateway_service/internal/client"
	"github.com/radyatamaa/go-cqrs-microservices/api_gateway_service/internal/domain"
	"github.com/radyatamaa/go-cqrs-microservices/api_gateway_service/internal/middlewares"
	readerService "github.com/radyatamaa/go-cqrs-microservices/api_gateway_service/proto/article_reader"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/interceptors"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/kafka"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/response"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"

	articleHandler "github.com/radyatamaa/go-cqrs-microservices/api_gateway_service/internal/article/delivery/http/v1"
	articleRepository "github.com/radyatamaa/go-cqrs-microservices/api_gateway_service/internal/article/repository"
	articleUsecase "github.com/radyatamaa/go-cqrs-microservices/api_gateway_service/internal/article/usecase"
)

// @title Api Gateway V1
// @version v1
// @contact.name radyatama
// @contact.email mohradyatama24@gmail.com
// @description api "API Gateway v1"
// @BasePath /api
// @query.collection.format multi

func main() {
	err := beego.LoadAppConfig("ini", "api_gateway_service/conf/app.ini")
	if err != nil {
		panic(err)
	}
	// global execution timeout
	serverTimeout := beego.AppConfig.DefaultInt64("serverTimeout", 60)
	// global execution timeout
	requestTimeout := beego.AppConfig.DefaultInt("executionTimeout", 5)
	// web hook to slack error log
	slackWebHookUrl := beego.AppConfig.DefaultString("slackWebhookUrlLog", "")
	// app version
	appVersion := beego.AppConfig.DefaultString("version", "1")
	// log path
	logPath := beego.AppConfig.DefaultString("logPath", "./logs/api_gateway_service.log")
	// grpc Reader Service Port
	grpcReaderServiceHost := beego.AppConfig.DefaultString("grpcReaderServiceHost", "localhost:5003")
	// brokers
	brokers := beego.AppConfig.DefaultStrings("brokers", []string{"localhost:9092"})
	// article create topic
	createArticleTopic := beego.AppConfig.DefaultString("createArticleTopic", "article_create")

	grpcReaderService := os.Getenv("READER_SERVICE")
	if grpcReaderService != "" {
		grpcReaderServiceHost = grpcReaderService
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers != "" {
		brokers = []string{kafkaBrokers}
	}

	// language
	lang := beego.AppConfig.DefaultString("lang", "en|id")
	languages := strings.Split(lang, "|")
	for _, value := range languages {
		if err := i18n.SetMessage(value, "./api_gateway_service/conf/"+value+".ini"); err != nil {
			panic("Failed to set message file for l10n")
		}
	}

	// global execution timeout to second
	timeoutContext := time.Duration(requestTimeout) * time.Second

	// beego config
	beego.BConfig.Log.AccessLogs = false
	beego.BConfig.Log.EnableStaticLogs = false
	beego.BConfig.Listen.ServerTimeOut = serverTimeout

	// zap logger
	zapLog := zaplogger.NewZapLogger(logPath, slackWebHookUrl)

	im := interceptors.NewInterceptorManager(zapLog)

	// init grpc
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	readerServiceConn, err := client.NewReaderServiceConn(ctx, grpcReaderServiceHost, im)
	if err != nil {
		panic(err)
	}
	defer readerServiceConn.Close() // nolint: errcheck
	rsClient := readerService.NewReaderServiceClient(readerServiceConn)

	// init kafka
	kafkaProducer := kafka.NewProducer(zapLog, brokers)
	defer kafkaProducer.Close() // nolint: errcheck
	confKafka := domain.ConfKafkaTopics{
		CreateArticle: createArticleTopic,
	}

	if beego.BConfig.RunMode != "prod" {
		// static files swagger
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	// middleware init
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowMethods:    []string{http.MethodGet, http.MethodPost},
		AllowAllOrigins: true,
	}))

	beego.InsertFilterChain("*", middlewares.RequestID())
	beego.InsertFilterChain("/api/*", middlewares.BodyDumpWithConfig(middlewares.NewAccessLogMiddleware(zapLog, appVersion).Logger()))

	// health check
	beego.Get("/health", func(ctx *beegoContext.Context) {
		ctx.Output.SetStatus(http.StatusOK)
		ctx.Output.JSON(beego.M{"status": "alive"}, beego.BConfig.RunMode != "prod", false)
	})

	// default error handler
	beego.ErrorController(&response.ErrorController{})

	// init repository
	articleQueriesRepository := articleRepository.NewQueriesArticleRepository(rsClient, zapLog)
	articleCommandRepository := articleRepository.NewCommandArticleRepository(kafkaProducer, confKafka, zapLog)

	// init usecase
	articleUcase := articleUsecase.NewArticleUseCase(timeoutContext, zapLog, articleCommandRepository, articleQueriesRepository)

	// init handler
	articleHandler.NewArticleHandler(articleUcase, zapLog)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		beego.Run()
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit

	pid := syscall.Getpid()

	switch sig {
	case syscall.SIGINT:
		log.Println(pid, "Received SIGINT.")
		log.Println(pid, "Waiting for connections to finish...")
		if err := beego.BeeApp.Server.Shutdown(ctx); err != nil {
			log.Fatal("failed shutdown server:", err)
		}
	case syscall.SIGTERM:
		log.Println(pid, "Received SIGTERM.")
		log.Println(pid, "Waiting for connections to finish...")
		if err := beego.BeeApp.Server.Shutdown(ctx); err != nil {
			log.Fatal("failed shutdown server:", err)
		}
	default:
		log.Printf("Received %v: nothing i care about...\n", sig)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("shutdown server success.")
	}
	log.Println("server exiting")
}
