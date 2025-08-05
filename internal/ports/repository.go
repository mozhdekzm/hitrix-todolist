package ports

import "github.com/mozhdekzm/heli-task/internal/domain"

type TodoRepository interface {
	Create(todo *domain.TodoItem) error
	GetAll() ([]domain.TodoItem, error)
}
