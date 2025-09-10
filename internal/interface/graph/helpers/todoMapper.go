package helpers

import (
	"github.com/mozhdekzm/gqlgql/internal/domain"
	"github.com/mozhdekzm/gqlgql/internal/interface/graph/model"
)

func MapDomainTodoToModel(t domain.TodoItem) *model.Todo {
	return &model.Todo{
		ID:          int64(t.ID),
		Description: t.Description,
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func MapDomainTodosToModels(todos []domain.TodoItem) []*model.Todo {
	result := make([]*model.Todo, len(todos))
	for i, t := range todos {
		result[i] = MapDomainTodoToModel(t)
	}
	return result
}
