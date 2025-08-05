package postgres

import (
	"database/sql"
	"fmt"
	"github.com/mozhdekzm/heli-task/config"
	"log"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dialect    string
	dbConfig   config.Config
	migrations *migrate.FileMigrationSource
}

func NewMigrator(cfg config.Config) Migrator {
	migrations := &migrate.FileMigrationSource{
		Dir: "./migrations",
	}

	return Migrator{
		dbConfig:   cfg,
		dialect:    "postgres",
		migrations: migrations,
	}
}

func (m Migrator) Up() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		m.dbConfig.DBHost,
		m.dbConfig.DBPort,
		m.dbConfig.DBUser,
		m.dbConfig.DBPassword,
		m.dbConfig.DBName,
	)

	db, err := sql.Open(m.dialect, connStr)
	if err != nil {
		log.Fatalf("can't open postgres db: %v", err)
	}
	defer db.Close()

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Up)
	if err != nil {
		log.Fatalf("can't apply migrations: %v", err)
	}
	fmt.Printf("Applied %d migrations!\n", n)
}

func (m Migrator) Down() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		m.dbConfig.DBHost,
		m.dbConfig.DBPort,
		m.dbConfig.DBUser,
		m.dbConfig.DBPassword,
		m.dbConfig.DBName,
	)

	db, err := sql.Open(m.dialect, connStr)
	if err != nil {
		log.Fatalf("can't open postgres db: %v", err)
	}
	defer db.Close()

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Down)
	if err != nil {
		log.Fatalf("can't rollback migrations: %v", err)
	}
	fmt.Printf("Rollbacked %d migrations!\n", n)
}
