package httpserver

import (
	"github.com/google/uuid"
	"github.com/mozhdekzm/heli-task/internal/domain"
	"time"
)

type CreateTodoRequest struct {
	Description string `json:"description" validate:"required"`
	DueDate     string `json:"dueDate" validate:"required"`
}

type TodoResponse struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
	CreatedAt   time.Time `json:"createdAt"`
}

func GetTodoResponse(todo domain.TodoItem) TodoResponse {
	return TodoResponse{
		ID:          todo.ID,
		Description: todo.Description,
		DueDate:     todo.DueDate,
		CreatedAt:   todo.CreatedAt,
	}
}
