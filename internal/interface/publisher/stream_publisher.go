package publisher

import (
	"context"
	"github.com/mozhdekzm/gqlgql/internal/domain"
)

//go:generate mockgen -destination=mock/stream_publisher_mock.go -package=publisher github.com/mozhdekzm/gqlgql/internal/interface/publisher StreamPublisher
type StreamMessage struct {
	ID     string                 `json:"id"`
	Values map[string]interface{} `json:"values"`
}

type StreamPublisher interface {
	Publish(ctx context.Context, todo domain.TodoItem) error
	Read(ctx context.Context, startID string, count int64) ([]StreamMessage, error)
}
