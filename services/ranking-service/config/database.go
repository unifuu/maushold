package config

import (
	"fmt"
	"log"
	"os"

	"maushold/ranking-service/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	err = db.AutoMigrate(&model.PlayerRanking{}, &model.LeaderboardEntry{}, &model.Player{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Execute SQL migrations (Materialized View, etc.)
	migrationPath := "migrations/001_add_combat_power_and_optimize.sql"
	content, err := os.ReadFile(migrationPath)
	if err == nil {
		db.Exec(string(content))
		log.Println("Applied SQL migrations")
	} else {
		log.Printf("Warning: Could not read migration file %s: %v", migrationPath, err)
	}

	log.Println("Database connected and migrated")
	return db
}
