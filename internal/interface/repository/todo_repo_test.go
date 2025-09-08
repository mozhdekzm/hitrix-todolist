package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/mozhdekzm/gqlgql/internal/domain"
	repository "github.com/mozhdekzm/gqlgql/internal/interface/repository/mock"
	"github.com/stretchr/testify/assert"
)

func TestMySQLTodoRepository_Save_ValidationAndPersistence(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockTodoRepository(ctrl)

	tests := []struct {
		name        string
		todo        domain.TodoItem
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid todo with all fields",
			todo: domain.TodoItem{
				ID:          uuid.New(),
				Description: "Complete task with all required fields",
				DueDate:     time.Date(2024, 8, 15, 10, 0, 0, 0, time.UTC),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			expectError: false,
		},
		{
			name: "todo with empty description",
			todo: domain.TodoItem{
				ID:          uuid.New(),
				Description: "",
				DueDate:     time.Date(2024, 8, 15, 10, 0, 0, 0, time.UTC),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			expectError: false,
		},
		{
			name: "todo with very long description",
			todo: domain.TodoItem{
				ID: uuid.New(),
				Description: "This is a very long description that tests the database field limits. " +
					"It contains many characters to ensure the repository can handle long text properly. " +
					"This should be saved successfully in the database without any truncation issues.",
				DueDate:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: false,
		},
		{
			name: "todo with nil UUID",
			todo: domain.TodoItem{
				ID:          uuid.Nil,
				Description: "Task with nil UUID",
				DueDate:     time.Date(2024, 8, 15, 10, 0, 0, 0, time.UTC),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			expectError: true,
			errorMsg:    "invalid UUID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectError {
				mockRepo.EXPECT().
					Save(gomock.Any(), gomock.Eq(tt.todo)).
					Return(errors.New(tt.errorMsg)).
					Times(1)
			} else {
				mockRepo.EXPECT().
					Save(gomock.Any(), gomock.Eq(tt.todo)).
					Return(nil).
					Times(1)
			}

			err := mockRepo.Save(context.Background(), tt.todo)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			if !tt.expectError {
				assert.NotEqual(t, uuid.Nil, tt.todo.ID)
			}
			assert.NotNil(t, tt.todo.CreatedAt)
			assert.NotNil(t, tt.todo.UpdatedAt)
		})
	}
}

func TestMySQLTodoRepository_GetAll_WithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockTodoRepository(ctrl)

	todos := []domain.TodoItem{
		{
			ID:          uuid.New(),
			Description: "Task 1",
			DueDate:     time.Date(2024, 9, 15, 10, 0, 0, 0, time.UTC),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Description: "Task 2",
			DueDate:     time.Date(2024, 9, 12, 10, 0, 0, 0, time.UTC),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	t.Run("success", func(t *testing.T) {
		limit := 10
		offset := 0
		mockRepo.EXPECT().
			GetAll(gomock.Any(), limit, offset).
			Return(todos, nil).
			Times(1)

		result, err := mockRepo.GetAll(context.Background(), limit, offset)
		assert.NoError(t, err)
		assert.Equal(t, todos, result)
	})

	t.Run("error from repo", func(t *testing.T) {
		limit := 5
		offset := 1
		expectedErr := errors.New("db error")
		mockRepo.EXPECT().
			GetAll(gomock.Any(), limit, offset).
			Return(nil, expectedErr).
			Times(1)

		result, err := mockRepo.GetAll(context.Background(), limit, offset)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}
