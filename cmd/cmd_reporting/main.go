package main

import (
	"context"
	"errors"
	"gamelight/internal/api"
	"gamelight/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	ApiHost       string `env:"api_host, default=0.0.0.0:8082"`
	RedisHost     string `env:"redis_host, default=localhost:6379"`
	RedisUsername string `env:"redis_username"`
	RedisPassword string `env:"redis_password"`
}

func main() {
	logrus.SetLevel(logrus.TraceLevel)

	ctx := context.Background()

	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		logrus.Fatalf("can not read config, err: %s", err.Error())
	}

	engine := gin.Default()
	store := storage.NewRedisStorage(cfg.RedisHost, cfg.RedisUsername, cfg.RedisPassword)

	handler := api.NewReportingHandler(engine, store)
	handler.RegisterRoutes()

	if err := http.ListenAndServe(cfg.ApiHost, engine); errors.Is(err, http.ErrServerClosed) {
		logrus.Info("http server listen error")
	} else if err != nil {
		logrus.Fatalf("http server err: %s", err.Error())
	}

	logrus.Exit(0)
}
