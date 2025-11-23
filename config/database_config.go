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
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
        DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        TranslateError: true,
    })

    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    db = db.Session(&gorm.Session{
        PrepareStmt: false,
    })

    return &databaseConfig{db: db}
}

func (cfg *databaseConfig) AutoMigrateAll(entities ...interface{}) error {
    if err := cfg.db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
        log.Printf("Warning: Could not create uuid-ossp extension: %v", err)
    }

    // Run migrations
    err := cfg.db.AutoMigrate(entities...)
    if err != nil {
        return fmt.Errorf("migration failed: %w", err)
    }

    return nil
}

func (cfg *databaseConfig) GetInstance() *gorm.DB {
    return cfg.db
}