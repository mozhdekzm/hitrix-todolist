package repository

import (
	"context"
	"gorm.io/gorm"

	"github.com/mozhdekzm/gqlgql/internal/domain"
)

//go:generate mockgen -destination=mock/todo_repo_mock.go -package=repository github.com/mozhdekzm/gqlgql/internal/interface/repository TodoRepository

type TodoRepository interface {
	BeginTx(ctx context.Context) (*gorm.DB, error)
	SaveWithTx(ctx context.Context, tx *gorm.DB, todo domain.TodoItem) error
	GetAll(ctx context.Context, limit int, offset int) ([]domain.TodoItem, error)
	FindByID(ctx context.Context, id string) (domain.TodoItem, error)
	UpdateWithTx(ctx context.Context, tx *gorm.DB, todo domain.TodoItem) error
	DeleteWithTx(ctx context.Context, tx *gorm.DB, id string) error
}
