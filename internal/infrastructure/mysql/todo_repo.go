package mysql

import (
	"context"
	"fmt"
	"github.com/mozhdekzm/gqlgql/internal/interface/repository"
	"log"

	//"fmt"
	//"github.com/google/uuid"
	"github.com/latolukasz/beeorm"
	"github.com/mozhdekzm/gqlgql/internal/domain"
)

type todoRepository struct {
	engine beeorm.Engine
}

func NewTodoRepository(engine beeorm.Engine) repository.TodoRepository {
	return &todoRepository{engine: engine}
}

func (r *todoRepository) Save(ctx context.Context, todo *domain.TodoItem) error {
	log.Printf("ENTITY: %+v", todo)
	fmt.Printf("engine: %+v\n", r.engine)
	fmt.Printf("registered entities: %+v\n", r.engine.GetRegistry())

	r.engine.Flush(todo)
	//flusher := r.engine.NewFlusher()
	//flusher.Track(todo)
	//if err := flusher.FlushWithCheck(); err != nil {
	//	return err
	//}

	return nil
}

func (r *todoRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.TodoItem, error) {
	var todos []domain.TodoItem
	//err := r.db.WithContext(ctx).
	//	Order("created_at desc").
	//	Limit(limit).
	//	Offset(offset).
	//	Find(&todos).Error
	//if err != nil {
	//	return nil, err
	//}
	return todos, nil
}

func (r *todoRepository) FindByID(ctx context.Context, id string) (domain.TodoItem, error) {
	var todo domain.TodoItem
	//uid, err := uuid.Parse(id)
	//if err != nil {
	//	return todo, fmt.Errorf("invalid UUID: %w", err)
	//}
	//if err := r.db.WithContext(ctx).First(&todo, "id = ?", uid).Error; err != nil {
	//	return todo, err
	//}
	return todo, nil
}

func (r *todoRepository) UpdateWithTx(ctx context.Context, todo domain.TodoItem) error {
	//return tx.WithContext(ctx).Save(&todo).Error
	return nil
}

func (r *todoRepository) DeleteWithTx(ctx context.Context, id string) error {
	//uid, err := uuid.Parse(id)
	//if err != nil {
	//	return fmt.Errorf("invalid UUID: %w", err)
	//}
	//return tx.WithContext(ctx).Delete(&domain.TodoItem{}, "id = ?", uid).Error
	return nil
}
