package repository

import (
	"context"
	"github.com/mozhdekzm/hitrix-todolist/internal/domain"
)

//go:generate mockgen -destination=mock/todo_repo_mock.go -package=repository github.com/mozhdekzm/gqlgql/internal/interface/repository TodoRepository

type TodoRepository interface {
	Save(ctx context.Context, todo *domain.TodoItem) error
	UpdateWithTx(ctx context.Context, todo domain.TodoItem) error
	DeleteWithTx(ctx context.Context, id uint64) error
	GetAll(ctx context.Context, limit int, offset int) ([]domain.TodoItem, error)
	FindByID(ctx context.Context, id uint64) (domain.TodoItem, error)
}
