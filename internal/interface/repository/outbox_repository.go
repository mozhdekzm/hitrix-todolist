package repository

import (
	"context"
	"github.com/mozhdekzm/hitrix-todolist/internal/domain"
)

type OutboxRepository interface {
	Save(ctx context.Context, event *domain.OutboxEvent) error
	GetUnpublished(ctx context.Context, limit int) ([]domain.OutboxEvent, error)
	MarkAsPublished(ctx context.Context, eventID uint64) error
	DeletePublished(ctx context.Context, olderThan int64) error
}
