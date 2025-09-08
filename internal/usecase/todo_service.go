package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/mozhdekzm/gqlgql/internal/domain"
	"github.com/mozhdekzm/gqlgql/internal/interface/publisher"
	"github.com/mozhdekzm/gqlgql/internal/interface/repository"
)

type TodoService struct {
	todoRepo       repository.TodoRepository
	redisPublisher publisher.StreamPublisher
}

func NewTodoService(todoRepo repository.TodoRepository, redisPublisher publisher.StreamPublisher) *TodoService {
	return &TodoService{
		todoRepo:       todoRepo,
		redisPublisher: redisPublisher,
	}
}

func (s *TodoService) Create(ctx context.Context, todo domain.TodoItem) (domain.TodoItem, error) {
	tx, err := s.todoRepo.BeginTx(ctx)
	if err != nil {
		return domain.TodoItem{}, fmt.Errorf("failed to begin tx: %w", err)
	}

	if err := s.todoRepo.SaveWithTx(ctx, tx, todo); err != nil {
		tx.Rollback()
		return domain.TodoItem{}, fmt.Errorf("failed to save todo: %w", err)
	}

	if err := s.redisPublisher.Publish(ctx, todo); err != nil {
		tx.Rollback()
		return domain.TodoItem{}, fmt.Errorf("failed to publish to redis: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return domain.TodoItem{}, fmt.Errorf("failed to commit tx: %w", err)
	}

	return todo, nil
}

func (s *TodoService) GetAll(ctx context.Context, limit int, offset int) ([]domain.TodoItem, error) {
	todos, err := s.todoRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (s *TodoService) FindByID(ctx context.Context, id string) (domain.TodoItem, error) {
	todo, err := s.todoRepo.FindByID(ctx, id)
	if err != nil {
		return domain.TodoItem{}, fmt.Errorf("failed to find todo by id: %w", err)
	}
	return todo, nil
}

func (s *TodoService) Update(ctx context.Context, todo domain.TodoItem) (domain.TodoItem, error) {
	tx, err := s.todoRepo.BeginTx(ctx)
	if err != nil {
		return domain.TodoItem{}, fmt.Errorf("failed to begin tx: %w", err)
	}

	todo.UpdatedAt = time.Now()

	if err := s.todoRepo.UpdateWithTx(ctx, tx, todo); err != nil {
		tx.Rollback()
		return domain.TodoItem{}, fmt.Errorf("failed to update todo: %w", err)
	}

	if err := s.redisPublisher.Publish(ctx, todo); err != nil {
		tx.Rollback()
		return domain.TodoItem{}, fmt.Errorf("failed to publish update to redis: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return domain.TodoItem{}, fmt.Errorf("failed to commit tx: %w", err)
	}

	return todo, nil
}

func (s *TodoService) Delete(ctx context.Context, id string) error {
	tx, err := s.todoRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	if err := s.todoRepo.DeleteWithTx(ctx, tx, id); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	if err := s.redisPublisher.Publish(ctx, domain.TodoItem{ID: uuid.MustParse(id)}); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to publish delete to redis: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}
