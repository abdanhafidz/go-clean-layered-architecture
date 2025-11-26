package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig interface {
	AutoMigrateAll(entities ...interface{}) error
	GetInstance() *gorm.DB
}
type databaseConfig struct {
	db *gorm.DB
}

func NewDatabaseConfig(DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT string) DatabaseConfig {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta "+
			"prefer_simple_protocol=true",
		DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})

	db = db.Session(&gorm.Session{
		PrepareStmt: false,
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return &databaseConfig{db: db}
}

func (cfg *databaseConfig) AutoMigrateAll(entities ...interface{}) error {

	err := cfg.db.AutoMigrate(
		entities...,
	)

	return err

}

func (cfg *databaseConfig) GetInstance() *gorm.DB {
	return cfg.db
}
