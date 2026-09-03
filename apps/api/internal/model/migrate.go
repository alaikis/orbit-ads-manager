package model

import (
	"log"

	"orbit/apps/api/pkg/database"
)

func Migrate() {
	if err := database.RunMigrations(); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("database migrated successfully")
}
