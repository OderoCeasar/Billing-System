package db

import (
	"fmt"
	"log"

	"github.com/OderoCeasar/system/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)


var DB *gorm.DB

func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)


	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})


	if err != nil {
		return fmt.Errorf("Failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}


func GetDB() *gorm.DB {
	return DB
}

