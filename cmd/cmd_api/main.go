package main

import (
	"context"
	"errors"
	"gamelight/internal/api"
	"gamelight/internal/eventbus"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AMQPUrl   string `env:"amqp_url, default=amqp://user:password@localhost:7001/"`
	AMQPQueue string `env:"amqp_queue, default=messages"`
	ApiHost   string `env:"api_host, default=0.0.0.0:8083"`
}

func main() {
	ctx := context.Background()

	logrus.SetLevel(logrus.TraceLevel)

	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		logrus.Fatalf("can not read config, err: %s", err.Error())
	}

	engine := gin.Default()

	producer, err := eventbus.NewAMQPProducer(cfg.AMQPUrl, cfg.AMQPQueue)
	if err != nil {
		panic(err)
	}

	handler := api.NewMessageHandler(engine, producer)
	handler.RegisterRoutes()

	if err := http.ListenAndServe(cfg.ApiHost, engine); errors.Is(err, http.ErrServerClosed) {
		logrus.Errorf("http server listen error")
	} else if err != nil {
		logrus.Fatal(err.Error())
	}

	logrus.Exit(0)
}
