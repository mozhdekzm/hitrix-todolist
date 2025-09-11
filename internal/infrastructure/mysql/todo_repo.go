package mysql

import (
	"context"
	"fmt"
	"github.com/mozhdekzm/gqlgql/internal/interface/repository"

	"git.ice.global/packages/beeorm/v4"
	"github.com/mozhdekzm/gqlgql/internal/domain"
)

type todoRepository struct {
	engine beeorm.Engine
}

func NewTodoRepository(engine beeorm.Engine) repository.TodoRepository {
	return &todoRepository{engine: engine}
}

func (r *todoRepository) Save(ctx context.Context, todo *domain.TodoItem) error {
	flusher := r.engine.NewFlusher()
	flusher.Track(todo)
	return flusher.FlushWithCheck()
}

func (r *todoRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.TodoItem, error) {
	var todos []*domain.TodoItem
	where := beeorm.NewWhere("1 = 1")
	pager := beeorm.NewPager(offset+1, limit)
	r.engine.Search(where, pager, &todos)
	// Convert pointers to values
	result := make([]domain.TodoItem, len(todos))
	for i, todo := range todos {
		result[i] = *todo
	}
	return result, nil
}

func (r *todoRepository) FindByID(ctx context.Context, id uint64) (domain.TodoItem, error) {
	var todo domain.TodoItem
	has := r.engine.LoadByID(id, &todo)
	if !has {
		return todo, fmt.Errorf("todo with id %d not found", id)
	}
	return todo, nil
}

func (r *todoRepository) UpdateWithTx(ctx context.Context, todo domain.TodoItem) error {
	// First load the existing entity to ensure it's tracked by beeorm
	var existingTodo domain.TodoItem
	has := r.engine.LoadByID(todo.ID, &existingTodo)
	if !has {
		return fmt.Errorf("todo with id %d not found", todo.ID)
	}

	// Update the fields
	existingTodo.Description = todo.Description
	existingTodo.DueDate = todo.DueDate
	existingTodo.UpdatedAt = todo.UpdatedAt

	// Track and flush the changes
	flusher := r.engine.NewFlusher()
	flusher.Track(&existingTodo)
	return flusher.FlushWithCheck()
}

func (r *todoRepository) DeleteWithTx(ctx context.Context, id uint64) error {
	var todo domain.TodoItem
	has := r.engine.LoadByID(id, &todo)
	if !has {
		return fmt.Errorf("todo with id %d not found", id)
	}
	flusher := r.engine.NewFlusher()
	flusher.Delete(&todo)
	return flusher.FlushWithCheck()
}
