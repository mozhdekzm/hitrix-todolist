package mysql

import (
	"fmt"

	"github.com/latolukasz/beeorm"
	"github.com/mozhdekzm/gqlgql/internal/config"
	"github.com/mozhdekzm/gqlgql/internal/domain"
)

func NewBeeORMEngine(cfg *config.Config) (beeorm.Engine, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	reg := beeorm.NewRegistry()

	reg.RegisterMySQLPool(dsn, "default")

	reg.RegisterEntity(&domain.TodoItem{})
	reg.RegisterEntity(&domain.OutboxEvent{})

	validated, err := reg.Validate()
	if err != nil {
		return nil, fmt.Errorf("beeorm registry validate: %w", err)
	}

	engine := validated.CreateEngine()
	return engine, nil
}
