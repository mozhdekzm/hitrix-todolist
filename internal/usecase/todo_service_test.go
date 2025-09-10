package usecase_test

import (
	"context"
	"errors"
	publisher "github.com/mozhdekzm/gqlgql/internal/interface/publisher/mock"
	repository "github.com/mozhdekzm/gqlgql/internal/interface/repository/mock"
	"github.com/mozhdekzm/gqlgql/internal/usecase"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/latolukasz/beeorm"
	"github.com/mozhdekzm/gqlgql/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestTodoService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := repository.NewMockTodoRepository(ctrl)
	outboxRepoMock := repository.NewMockOutboxRepository(ctrl)
	pubMock := publisher.NewMockStreamPublisher(ctrl)

	// Create a mock BeeORM engine - for testing we'll use nil
	var mockEngine beeorm.Engine

	service := usecase.NewTodoService(repoMock, outboxRepoMock, pubMock, mockEngine)

	todo := domain.TodoItem{
		ID:          1,
		Description: "test task",
		DueDate:     time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		repoMock.EXPECT().Save(gomock.Any(), todo).Return(nil)
		pubMock.EXPECT().Publish(gomock.Any(), todo).Return(nil)

		result, err := service.Create(context.Background(), todo)
		assert.NoError(t, err)
		assert.Equal(t, todo, result)
	})

	t.Run("repo error", func(t *testing.T) {
		repoMock.EXPECT().Save(gomock.Any(), todo).Return(errors.New("db error"))

		_, err := service.Create(context.Background(), todo)
		assert.Error(t, err)
	})

	t.Run("publish error", func(t *testing.T) {
		repoMock.EXPECT().Save(gomock.Any(), todo).Return(nil)
		pubMock.EXPECT().Publish(gomock.Any(), todo).Return(errors.New("redis error"))

		_, err := service.Create(context.Background(), todo)
		assert.Error(t, err)
	})
}
