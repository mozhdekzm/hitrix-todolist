package domain

import (
	"github.com/latolukasz/beeorm"
	"time"

	"github.com/google/uuid"
)

type TodoItem struct {
	beeorm.ORM  `orm:"table=todos;primary_key=id"`
	ID          string    `json:"id"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewTodoItem(description string, dueDate time.Time) TodoItem {
	now := time.Now()
	return TodoItem{
		ID:          uuid.New().String(),
		Description: description,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (TodoItem) TableName() string {
	return "todos"
}
