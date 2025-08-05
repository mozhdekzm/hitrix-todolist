package tests

import (
	"github.com/mozhdekzm/heli-task/internal/application"
	"testing"
	"time"

	"github.com/mozhdekzm/heli-task/internal/domain"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct{ Items []domain.TodoItem }

func (m *mockRepo) Create(item *domain.TodoItem) error {
	m.Items = append(m.Items, *item)
	return nil
}
func (m *mockRepo) GetAll() ([]domain.TodoItem, error) { return m.Items, nil }

type mockQueue struct{ Published []domain.TodoItem }

func (q *mockQueue) Publish(item domain.TodoItem) error {
	q.Published = append(q.Published, item)
	return nil
}

func TestTodoService_Create_TableDriven(t *testing.T) {
	repo := &mockRepo{}
	queue := &mockQueue{}
	service := application.NewTodoService(repo, queue)

	tests := []struct {
		name        string
		description string
		dueDate     time.Time
		expectError bool
	}{
		{"valid todo", "Task 1", time.Now().Add(24 * time.Hour), false},
		{"empty description", "", time.Now().Add(24 * time.Hour), true},
		{"past due date", "Old Task", time.Now().Add(-24 * time.Hour), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo := domain.TodoItem{
				Description: tt.description,
				DueDate:     tt.dueDate,
			}

			todo, err := service.Create(todo.Description, todo.DueDate)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	assert.Equal(t, 1, len(repo.Items))
	assert.Equal(t, 1, len(queue.Published))
}
