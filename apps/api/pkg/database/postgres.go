package database

import (
	"fmt"
	"log"
	"orbit/apps/api/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	cfg := config.AppCfg
	dsn := cfg.Database.URL
	if dsn == "" {
		dsn = "postgres://orbit:orbit@localhost:5432/orbit?sslmode=disable"
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	fmt.Println("database connected")
}
