package domain

import (
	"github.com/google/uuid"
	"time"
)

type TodoItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
	CreatedAt   time.Time `json:"created_at"`
}
