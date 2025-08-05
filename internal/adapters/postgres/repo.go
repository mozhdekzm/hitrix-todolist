package postgres

import (
	"github.com/mozhdekzm/heli-task/internal/domain"
	"gorm.io/gorm"
)

type TodoRepo struct {
	DB *gorm.DB
}

func NewTodoRepo(db *gorm.DB) *TodoRepo {
	return &TodoRepo{DB: db}
}

func (r *TodoRepo) Create(todo *domain.TodoItem) error {
	return r.DB.Create(todo).Error
}

func (r *TodoRepo) GetAll() ([]domain.TodoItem, error) {
	var todos []domain.TodoItem
	result := r.DB.Order("id desc").Find(&todos)
	return todos, result.Error
}
