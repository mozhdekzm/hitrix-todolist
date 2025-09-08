package redis

import (
	"context"
	"github.com/mozhdekzm/gqlgql/internal/config"
	"github.com/redis/go-redis/v9"
	"log"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("[redis] failed to connect: %v", err)
	}

	return client
}

func ReadStream(ctx context.Context, client *redis.Client, stream string, startID string, count int64) ([]redis.XMessage, error) {
	if startID == "" {
		startID = "0-0"
	}
	if count <= 0 {
		count = 100
	}

	res, err := client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{stream, startID},
		Count:   count,
		Block:   0,
	}).Result()
	if err != nil {
		return nil, err
	}

	var messages []redis.XMessage
	for _, s := range res {
		messages = append(messages, s.Messages...)
	}
	return messages, nil
}
