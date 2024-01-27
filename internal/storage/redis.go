package storage

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gamelight/internal/types"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const KeyDialogMessages = "messages"

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr, username, password string) *RedisStorage {
	opts := &redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
	}

	client := redis.NewClient(opts)

	return &RedisStorage{client: client}
}

// We need to build a key to get whole dialog between two parties

func (r *RedisStorage) buildDialogKey(sender, receiver string) string {
	var key string

	if receiver > sender {
		key = fmt.Sprintf("%s-%s", receiver, sender)
	} else {
		key = fmt.Sprintf("%s-%s", sender, receiver)
	}

	b64 := base64.StdEncoding.EncodeToString([]byte(key))

	return fmt.Sprintf("%s.%s", KeyDialogMessages, b64)
}

func (r *RedisStorage) Save(ctx context.Context, message *types.Message) error {
	dialogKey := r.buildDialogKey(message.Sender, message.Receiver)

	b, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("can not marshal message, err: %v", err)
	}

	if err := r.client.LPush(ctx, dialogKey, b).Err(); err != nil {
		return fmt.Errorf("can not save message, err: %v", err)
	}

	return nil
}

func (r *RedisStorage) List(ctx context.Context, sender, receiver string) ([]*types.Message, error) {
	dialogKey := r.buildDialogKey(sender, receiver)

	messages, err := r.client.LRange(ctx, dialogKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("can not retrive messages, err: %v", err)
	}

	rsp := make([]*types.Message, 0, len(messages))
	for _, msg := range messages {
		var message types.Message

		if err = json.Unmarshal([]byte(msg), &message); err != nil {
			logrus.Errorf("can not parse message: %s, err: %s. Skip", msg, err.Error())

			continue
		}

		rsp = append(rsp, &message)
	}

	return rsp, nil
}
