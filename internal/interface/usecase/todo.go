package usecase

import (
	"context"

	"github.com/mozhdekzm/gqlgql/internal/domain"
)

type TodoService interface {
	Create(ctx context.Context, todo domain.TodoItem) (domain.TodoItem, error)
	GetAll(ctx context.Context, limit int, offset int) ([]domain.TodoItem, error)
	FindByID(ctx context.Context, id string) (domain.TodoItem, error)
	Update(ctx context.Context, todo domain.TodoItem) (domain.TodoItem, error)
	Delete(ctx context.Context, id string) error
}
