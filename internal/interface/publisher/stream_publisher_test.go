package publisher_test

import (
	"context"
	"errors"
	"github.com/mozhdekzm/gqlgql/internal/config"
	redis2 "github.com/mozhdekzm/gqlgql/internal/infrastructure/redis"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mozhdekzm/gqlgql/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestStreamPublisher_Publish_WithGomock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedisClient := &mockRedisClient{}

	cfg := &config.Config{RedisStream: "test-stream"}
	publisher := redis2.NewStreamPublisher(mockRedisClient, cfg)

	todo := domain.TodoItem{
		ID:          1,
		Description: "Test task",
		DueDate:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRedisClient.xAddResult = "1234567890-0"
	mockRedisClient.xAddError = nil

	err := publisher.Publish(context.Background(), todo)
	assert.NoError(t, err)

	mockRedisClient.xAddError = errors.New("redis connection failed")

	err = publisher.Publish(context.Background(), todo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish to redis stream")
}

func TestStreamPublisher_Read_WithGomock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedisClient := &mockRedisClient{}

	cfg := &config.Config{RedisStream: "test-stream"}
	publisher := redis2.NewStreamPublisher(mockRedisClient, cfg)

	mockMessages := []redis.XMessage{
		{
			ID: "1234567890-0",
			Values: map[string]interface{}{
				"id":          "test-uuid-1",
				"description": "Task 1",
				"due_date":    "2024-01-15T10:00:00Z",
			},
		},
		{
			ID: "1234567890-1",
			Values: map[string]interface{}{
				"id":          "test-uuid-2",
				"description": "Task 2",
				"due_date":    "2024-01-16T10:00:00Z",
			},
		},
	}

	mockRedisClient.xReadResult = []redis.XStream{
		{
			Stream:   "test-stream",
			Messages: mockMessages,
		},
	}
	mockRedisClient.xReadError = nil

	messages, err := publisher.Read(context.Background(), "", 0)
	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, "1234567890-0", messages[0].ID)

	mockRedisClient.xReadError = errors.New("redis read failed")
	mockRedisClient.xReadResult = nil

	_, err = publisher.Read(context.Background(), "0-0", 100)
	assert.Error(t, err)
	assert.Equal(t, "redis read failed", err.Error())
}

type mockRedisClient struct {
	xAddResult  string
	xAddError   error
	xReadResult []redis.XStream
	xReadError  error
}

func (m *mockRedisClient) XAdd(ctx context.Context, args *redis.XAddArgs) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx, "xadd", args.Stream)
	if m.xAddError != nil {
		cmd.SetErr(m.xAddError)
	} else {
		cmd.SetVal(m.xAddResult)
	}
	return cmd
}

func (m *mockRedisClient) XRead(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd {
	cmd := redis.NewXStreamSliceCmd(ctx, "xread")
	if m.xReadError != nil {
		cmd.SetErr(m.xReadError)
	} else {
		cmd.SetVal(m.xReadResult)
	}
	return cmd
}
