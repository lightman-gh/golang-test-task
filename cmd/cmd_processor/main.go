package main

import (
	"context"
	"gamelight/internal/eventbus"
	"gamelight/internal/storage"
	"gamelight/internal/types"

	"github.com/sethvargo/go-envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AMQPUrl       string `env:"amqp_url, default=amqp://user:password@localhost:7001/"`
	AMQPQueue     string `env:"amqp_queue, default=messages"`
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

	consumer, err := eventbus.NewAMQPConsumer(cfg.AMQPUrl, cfg.AMQPQueue)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	store := storage.NewRedisStorage(cfg.RedisHost, cfg.RedisUsername, cfg.RedisPassword)

	err = consumer.Consume(ctx, func(ctx context.Context, message *types.Message) error {
		return store.Save(ctx, message)
	})

	logrus.Exit(0)
}
