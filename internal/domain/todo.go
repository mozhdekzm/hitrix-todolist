package domain

import (
	"time"

	"github.com/google/uuid"
)

type TodoItem struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewTodoItem(description string, dueDate time.Time) TodoItem {
	now := time.Now()
	return TodoItem{
		ID:          uuid.New(),
		Description: description,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (TodoItem) TableName() string {
	return "todos"
}
