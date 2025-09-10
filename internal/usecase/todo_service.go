package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/latolukasz/beeorm"
	"github.com/mozhdekzm/gqlgql/internal/domain"
	"github.com/mozhdekzm/gqlgql/internal/interface/publisher"
	"github.com/mozhdekzm/gqlgql/internal/interface/repository"
)

type TodoService struct {
	todoRepo       repository.TodoRepository
	outboxRepo     repository.OutboxRepository
	redisPublisher publisher.StreamPublisher
	engine         beeorm.Engine
}

func NewTodoService(todoRepo repository.TodoRepository, outboxRepo repository.OutboxRepository, redisPublisher publisher.StreamPublisher, engine beeorm.Engine) *TodoService {
	return &TodoService{
		todoRepo:       todoRepo,
		outboxRepo:     outboxRepo,
		redisPublisher: redisPublisher,
		engine:         engine,
	}
}

func (s *TodoService) Create(ctx context.Context, todo domain.TodoItem) (domain.TodoItem, error) {
	// Create outbox event
	outboxEvent, err := domain.NewOutboxEvent("CREATE", "TodoItem", todo.ID, todo)
	if err != nil {
		return domain.TodoItem{}, fmt.Errorf("failed to create outbox event: %w", err)
	}

	// Use beeorm transaction to ensure consistency
	flusher := s.engine.NewFlusher()

	// Track todo and outbox event together
	flusher.Track(&todo)
	flusher.Track(outboxEvent)

	// Commit both in single transaction
	if err := flusher.FlushWithCheck(); err != nil {
		return domain.TodoItem{}, fmt.Errorf("failed to save todo and outbox event: %w", err)
	}

	// Publish event asynchronously (best effort)
	go func() {
		if err := s.redisPublisher.PublishOutboxEvent(context.Background(), *outboxEvent); err != nil {
			fmt.Printf("Failed to publish outbox event: %v\n", err)
		} else {
			// Mark as published if successful
			s.outboxRepo.MarkAsPublished(context.Background(), outboxEvent.ID)
		}
	}()

	return todo, nil
}

func (s *TodoService) GetAll(ctx context.Context, limit int, offset int) ([]domain.TodoItem, error) {
	todos, err := s.todoRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (s *TodoService) FindByID(ctx context.Context, id uint64) (domain.TodoItem, error) {
	todo, err := s.todoRepo.FindByID(ctx, id)
	if err != nil {
		return domain.TodoItem{}, fmt.Errorf("failed to find todo by id: %w", err)
	}
	return todo, nil
}

func (s *TodoService) Update(ctx context.Context, todo domain.TodoItem) (domain.TodoItem, error) {
	todo.UpdatedAt = time.Now()

	// Create outbox event
	outboxEvent, err := domain.NewOutboxEvent("UPDATE", "TodoItem", todo.ID, todo)
	if err != nil {
		return domain.TodoItem{}, fmt.Errorf("failed to create outbox event: %w", err)
	}

	// Load existing todo first for beeorm tracking
	existingTodo, err := s.todoRepo.FindByID(ctx, todo.ID)
	if err != nil {
		return domain.TodoItem{}, fmt.Errorf("failed to find existing todo: %w", err)
	}

	// Update fields
	existingTodo.Description = todo.Description
	existingTodo.DueDate = todo.DueDate
	existingTodo.UpdatedAt = todo.UpdatedAt

	// Use beeorm transaction
	flusher := s.engine.NewFlusher()
	flusher.Track(&existingTodo)
	flusher.Track(outboxEvent)

	if err := flusher.FlushWithCheck(); err != nil {
		return domain.TodoItem{}, fmt.Errorf("failed to update todo and save outbox event: %w", err)
	}

	// Publish event asynchronously
	go func() {
		if err := s.redisPublisher.PublishOutboxEvent(context.Background(), *outboxEvent); err != nil {
			fmt.Printf("Failed to publish update outbox event: %v\n", err)
		} else {
			s.outboxRepo.MarkAsPublished(context.Background(), outboxEvent.ID)
		}
	}()

	return existingTodo, nil
}

func (s *TodoService) Delete(ctx context.Context, id uint64) error {
	// Create outbox event
	outboxEvent, err := domain.NewOutboxEvent("DELETE", "TodoItem", id, map[string]interface{}{"id": id})
	if err != nil {
		return fmt.Errorf("failed to create outbox event: %w", err)
	}

	// Load existing todo for beeorm tracking
	var existingTodo domain.TodoItem
	has := s.engine.LoadByID(id, &existingTodo)
	if !has {
		return fmt.Errorf("todo with id %d not found", id)
	}

	// Use beeorm transaction
	flusher := s.engine.NewFlusher()
	flusher.Delete(&existingTodo)
	flusher.Track(outboxEvent)

	if err := flusher.FlushWithCheck(); err != nil {
		return fmt.Errorf("failed to delete todo and save outbox event: %w", err)
	}

	// Publish event asynchronously
	go func() {
		if err := s.redisPublisher.PublishOutboxEvent(context.Background(), *outboxEvent); err != nil {
			fmt.Printf("Failed to publish delete outbox event: %v\n", err)
		} else {
			s.outboxRepo.MarkAsPublished(context.Background(), outboxEvent.ID)
		}
	}()

	return nil
}
