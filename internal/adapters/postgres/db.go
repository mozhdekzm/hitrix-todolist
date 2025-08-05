package postgres

import (
	"fmt"
	"github.com/mozhdekzm/heli-task/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func NewGormDB(cfg config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	return db
}
