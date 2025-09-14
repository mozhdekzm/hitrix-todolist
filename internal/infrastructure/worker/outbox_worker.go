package worker

import (
	"context"
	"log"
	"time"

	"github.com/mozhdekzm/hitrix-todolist/internal/interface/publisher"
	"github.com/mozhdekzm/hitrix-todolist/internal/interface/repository"
)

type OutboxWorker struct {
	outboxRepo      repository.OutboxRepository
	streamPublisher publisher.StreamPublisher
	interval        time.Duration
	batchSize       int
}

func NewOutboxWorker(outboxRepo repository.OutboxRepository, streamPublisher publisher.StreamPublisher) *OutboxWorker {
	return &OutboxWorker{
		outboxRepo:      outboxRepo,
		streamPublisher: streamPublisher,
		interval:        5 * time.Second, // Process every 5 seconds
		batchSize:       100,             // Process 100 events at a time
	}
}

func (w *OutboxWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Println("Outbox worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Outbox worker stopped")
			return
		case <-ticker.C:
			w.processOutboxEvents(ctx)
		}
	}
}

func (w *OutboxWorker) processOutboxEvents(ctx context.Context) {
	events, err := w.outboxRepo.GetUnpublished(ctx, w.batchSize)
	if err != nil {
		log.Printf("Failed to get unpublished events: %v", err)
		return
	}

	if len(events) == 0 {
		return
	}

	log.Printf("Processing %d outbox events", len(events))

	for _, event := range events {
		if err := w.streamPublisher.PublishOutboxEvent(ctx, event); err != nil {
			log.Printf("Failed to publish event %d: %v", event.ID, err)
			continue
		}

		if err := w.outboxRepo.MarkAsPublished(ctx, event.ID); err != nil {
			log.Printf("Failed to mark event %d as published: %v", event.ID, err)
		}
	}
}

func (w *OutboxWorker) Cleanup(ctx context.Context) {
	// Clean up events older than 24 hours
	olderThan := time.Now().Add(-24 * time.Hour).Unix()
	if err := w.outboxRepo.DeletePublished(ctx, olderThan); err != nil {
		log.Printf("Failed to cleanup old events: %v", err)
	} else {
		log.Println("Cleaned up old published events")
	}
}
