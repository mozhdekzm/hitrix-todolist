package mysql

import (
	"fmt"
	"log"
	"time"

	"github.com/mozhdekzm/gqlgql/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGormClient(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var db *gorm.DB
	var err error

	maxAttempts := 12
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			return db
		}
		log.Printf("[mysql] attempt %d/%d failed: %v", attempt, maxAttempts, err)
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	log.Fatalf("[mysql] failed to connect after %d attempts: %v", maxAttempts, err)
	return nil
}
