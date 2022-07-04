package main

import (
	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
	"github.com/radyatamaa/go-cqrs-microservices/write_service/config"
	"github.com/radyatamaa/go-cqrs-microservices/write_service/internal/server"
)

func main() {

	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	// zap logger
	zaplog := zaplogger.NewZapLogger(cfg.App.LogPath, cfg.App.SlackWebHookUrl)
	zaplog.WithName("WriterService")

	s := server.NewServer(cfg, zaplog)
	zaplog.Fatalf("running server : %s", s.Run())

}
