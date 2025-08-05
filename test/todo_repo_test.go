package tests

import (
	"fmt"
	postgres2 "github.com/mozhdekzm/heli-task/internal/adapters/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mozhdekzm/heli-task/internal/domain"
)

func TestTodoRepoCreateAndPersistPostgres(t *testing.T) {
	dsn := "host=localhost user=postgres password=todo dbname=todo port=5434 "
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&domain.TodoItem{})
	assert.NoError(t, err)

	repo := postgres2.NewTodoRepo(db)

	tests := []struct {
		name        string
		description string
		dueDate     time.Time
		expectError bool
	}{
		{"Valid TodoItem", "Integration test", time.Now().Add(24 * time.Hour), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo := domain.TodoItem{
				Description: tt.description,
				DueDate:     tt.dueDate,
			}

			err := repo.Create(&todo)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			t.Cleanup(func() {
				db.Exec("DELETE FROM todo_items WHERE id = ?", todo.ID)
			})

			items, err := repo.GetAll()
			assert.NoError(t, err)

			found := false
			for _, item := range items {
				if item.Description == tt.description {
					found = true
					break
				}
			}
			assert.True(t, found, "Created TodoItem was not found in PostgreSQL")
			fmt.Println("Integration test completed for:", tt.description)
		})
	}
}
