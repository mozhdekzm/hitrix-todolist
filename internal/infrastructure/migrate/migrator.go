package migrate

import (
	"database/sql"
	"fmt"
	"github.com/mozhdekzm/hitrix-todolist/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dialect    string
	dbConfig   *config.Config
	migrations *migrate.FileMigrationSource
}

func NewMigrator(cfg *config.Config) Migrator {
	migrations := &migrate.FileMigrationSource{
		Dir: "./migrations",
	}

	return Migrator{
		dbConfig:   cfg,
		dialect:    "mysql",
		migrations: migrations,
	}
}

func (m Migrator) Up() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		m.dbConfig.DBUser,
		m.dbConfig.DBPassword,
		m.dbConfig.DBHost,
		m.dbConfig.DBPort,
		m.dbConfig.DBName,
	)

	db, err := sql.Open(m.dialect, dsn)
	if err != nil {
		log.Fatalf("can't open mysql db: %v", err)
	}
	defer db.Close()

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Up)
	if err != nil {
		log.Fatalf("can't apply migrations: %v", err)
	}
	fmt.Printf("✅ Applied %d migrations\n", n)
}

func (m Migrator) Down() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		m.dbConfig.DBUser,
		m.dbConfig.DBPassword,
		m.dbConfig.DBHost,
		m.dbConfig.DBPort,
		m.dbConfig.DBName,
	)

	db, err := sql.Open(m.dialect, dsn)
	if err != nil {
		log.Fatalf("can't open mysql db: %v", err)
	}
	defer db.Close()

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Down)
	if err != nil {
		log.Fatalf("can't rollback migrations: %v", err)
	}
	fmt.Printf("⬅️  Rolled back %d migrations\n", n)
}
