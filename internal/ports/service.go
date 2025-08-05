package ports

import (
	"github.com/mozhdekzm/heli-task/internal/domain"
	"time"
)

type TodoService interface {
	Create(description string, dueDate time.Time) (domain.TodoItem, error)
	List() ([]domain.TodoItem, error)
}
