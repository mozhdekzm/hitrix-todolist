package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/latolukasz/beeorm"
	"github.com/mozhdekzm/gqlgql/internal/domain"
	"github.com/mozhdekzm/gqlgql/internal/interface/repository"
)

type outboxRepository struct {
	engine beeorm.Engine
}

func NewOutboxRepository(engine beeorm.Engine) repository.OutboxRepository {
	return &outboxRepository{engine: engine}
}

func (r *outboxRepository) Save(ctx context.Context, event *domain.OutboxEvent) error {
	flusher := r.engine.NewFlusher()
	flusher.Track(event)
	return flusher.FlushWithCheck()
}

func (r *outboxRepository) GetUnpublished(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	var events []*domain.OutboxEvent
	where := beeorm.NewWhere("Published = ?", false)
	pager := beeorm.NewPager(1, limit)
	r.engine.Search(where, pager, &events)

	// Convert pointers to values
	result := make([]domain.OutboxEvent, len(events))
	for i, event := range events {
		result[i] = *event
	}
	return result, nil
}

func (r *outboxRepository) MarkAsPublished(ctx context.Context, eventID uint64) error {
	var event domain.OutboxEvent
	has := r.engine.LoadByID(eventID, &event)
	if !has {
		return fmt.Errorf("outbox event with id %d not found", eventID)
	}

	event.Published = true
	event.UpdatedAt = time.Now()

	flusher := r.engine.NewFlusher()
	flusher.Track(&event)
	return flusher.FlushWithCheck()
}

func (r *outboxRepository) DeletePublished(ctx context.Context, olderThan int64) error {
	var events []*domain.OutboxEvent
	where := beeorm.NewWhere("Published = ? AND CreatedAt < ?", true, time.Unix(olderThan, 0))
	pager := beeorm.NewPager(1, 1000) // Process in batches
	r.engine.Search(where, pager, &events)

	if len(events) == 0 {
		return nil
	}

	flusher := r.engine.NewFlusher()
	for _, event := range events {
		flusher.Delete(event)
	}
	return flusher.FlushWithCheck()
}
