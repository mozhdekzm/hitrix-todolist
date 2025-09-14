package publisher

import (
	"context"
	"github.com/mozhdekzm/hitrix-todolist/internal/domain"
)

//go:generate mockgen -destination=mock/stream_publisher_mock.go -package=publisher github.com/mozhdekzm/gqlgql/internal/interface/publisher StreamPublisher
type StreamMessage struct {
	ID     string
	Values map[string]interface{}
}

type StreamPublisher interface {
	PublishCreate(ctx context.Context, todo domain.TodoItem) error
	PublishUpdate(ctx context.Context, todo domain.TodoItem) error
	PublishDelete(ctx context.Context, id uint64) error
	PublishOutboxEvent(ctx context.Context, event domain.OutboxEvent) error
	Read(ctx context.Context, startID string, count int64) ([]StreamMessage, error)
}
