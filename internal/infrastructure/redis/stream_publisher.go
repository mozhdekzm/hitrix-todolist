package redis

import (
	"context"
	"fmt"
	"github.com/mozhdekzm/gqlgql/internal/config"
	"github.com/mozhdekzm/gqlgql/internal/domain"
	"github.com/mozhdekzm/gqlgql/internal/interface/publisher"
	"github.com/redis/go-redis/v9"
)

// use for tests
type RedisClient interface {
	XAdd(ctx context.Context, args *redis.XAddArgs) *redis.StringCmd
	XRead(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd
}

type streamPublisher struct {
	client RedisClient
	stream string
}

func NewStreamPublisher(client RedisClient, cfg *config.Config) publisher.StreamPublisher {
	return &streamPublisher{
		client: client,
		stream: cfg.RedisStream,
	}
}

func (p *streamPublisher) Publish(ctx context.Context, todo domain.TodoItem) error {
	_, err := p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: p.stream,
		Values: map[string]interface{}{
			"id":          todo.ID,
			"description": todo.Description,
			"due_date":    todo.DueDate.Format("2006-01-02T15:04:05Z07:00"),
		},
	}).Result()
	if err != nil {
		return fmt.Errorf("failed to publish to redis stream: %w", err)
	}
	return nil
}

func (p *streamPublisher) Read(ctx context.Context, startID string, count int64) ([]publisher.StreamMessage, error) {
	if startID == "" {
		startID = "0-0"
	}
	if count <= 0 {
		count = 100
	}

	res, err := p.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{p.stream, startID},
		Count:   count,
		Block:   0,
	}).Result()
	if err != nil {
		return nil, err
	}

	var messages []publisher.StreamMessage
	for _, s := range res {
		for _, m := range s.Messages {
			messages = append(messages, publisher.StreamMessage{
				ID:     m.ID,
				Values: m.Values,
			})
		}
	}
	return messages, nil
}
