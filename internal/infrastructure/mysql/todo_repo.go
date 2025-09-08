package mysql

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mozhdekzm/gqlgql/internal/domain"
	"github.com/mozhdekzm/gqlgql/internal/interface/repository"
	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) repository.TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (r *todoRepository) SaveWithTx(ctx context.Context, tx *gorm.DB, todo domain.TodoItem) error {
	return tx.WithContext(ctx).Create(todo).Error
}

func (r *todoRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.TodoItem, error) {
	var todos []domain.TodoItem
	err := r.db.WithContext(ctx).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *todoRepository) FindByID(ctx context.Context, id string) (domain.TodoItem, error) {
	var todo domain.TodoItem
	uid, err := uuid.Parse(id)
	if err != nil {
		return todo, fmt.Errorf("invalid UUID: %w", err)
	}
	if err := r.db.WithContext(ctx).First(&todo, "id = ?", uid).Error; err != nil {
		return todo, err
	}
	return todo, nil
}

func (r *todoRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, todo domain.TodoItem) error {
	return tx.WithContext(ctx).Save(&todo).Error
}

func (r *todoRepository) DeleteWithTx(ctx context.Context, tx *gorm.DB, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}
	return tx.WithContext(ctx).Delete(&domain.TodoItem{}, "id = ?", uid).Error
}
