package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Package-level variable to hold the GORM database instance
var DB *gorm.DB

// Initialize the db connection using the provided DSN(Data Source Name)
func Init(dsn string) (*gorm.DB, error) {
	// Setup GORM Logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Setup GORM Config
	gormConfig := &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Connection Pool settings (optional, GORM provides sensible defaults).
	// sqlDB, err := db.DB() // Get the underlying *sql.DB instance from GORM.
	// if err != nil {
	// return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	// }
	// sqlDB.SetMaxIdleConns(10)           // Maximum number of connections in the idle connection pool.
	// sqlDB.SetMaxOpenConns(100)          // Maximum number of open connections to the database.
	// sqlDB.SetConnMaxLifetime(time.Hour) // Maximum amount of time a connection may be reused.

	return db, nil
}
