package application

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mozhdekzm/heli-task/internal/domain"
	"github.com/mozhdekzm/heli-task/internal/ports"
	"time"
)

type TodoService struct {
	Repo  ports.TodoRepository
	Queue ports.Queue
}

func NewTodoService(repo ports.TodoRepository, queue ports.Queue) *TodoService {
	return &TodoService{
		Repo:  repo,
		Queue: queue,
	}
}

func (s *TodoService) Create(description string, dueDate time.Time) (domain.TodoItem, error) {
	if description == "" {
		return domain.TodoItem{}, errors.New("description is required")
	}
	if dueDate.Before(time.Now()) {
		return domain.TodoItem{}, errors.New("due date cannot be in the past")
	}

	todo := domain.TodoItem{
		ID:          uuid.New(),
		Description: description,
		DueDate:     dueDate,
	}
	if err := s.Repo.Create(&todo); err != nil {
		return domain.TodoItem{}, err
	}
	_ = s.Queue.Publish(todo)
	return todo, nil
}

func (s *TodoService) List() ([]domain.TodoItem, error) {
	return s.Repo.GetAll()
}
