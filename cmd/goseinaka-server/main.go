package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ecabigting/goseinaka/internal/config"
	"github.com/ecabigting/goseinaka/internal/database"
	"github.com/ecabigting/goseinaka/internal/domain"
	"github.com/ecabigting/goseinaka/internal/handler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Instance of the GORM Database connection
var db *gorm.DB

func main() {
	log.Println("Starting API...")
	apiConfig, err := config.Load()
	if err != nil {
		log.Fatalf("Error Loading Config file: %s", err)
	}

	// Print the loaded config if
	// log level is
	// set to 'debug'
	if apiConfig.LogLevel == "debug" {
		apiConfigJSON, jsonErr := json.MarshalIndent(apiConfig, "", " ")
		if jsonErr != nil {
			log.Printf("DEBUG: Loaded config fallback: %+v", apiConfig)
		} else {
			log.Printf("DEBUG: Loaded config (JSON):\n%s\n", string(apiConfigJSON))
		}
	}

	// Init database connection db, err := database.Init(apiConfig.DatabaseURL)
	db, err := database.Init(apiConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	log.Println("Successfully connected to the database..")

	// Graceful Shutdown the DB Connection
	sqlDB, _ := db.DB()
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Fatal("Failed to close db connection with error: ", err)
		} else {
			log.Println("Database connection closed successfully.")
		}
	}()

	// Check DB Stats
	var serverInfoMap map[string]any
	sqlQuery := "SELECT    NOW() as current_db_time,     current_setting('TimeZone') as db_timezone;"
	result := db.Raw(sqlQuery).Scan(&serverInfoMap)

	if result.Error != nil {
		log.Fatal("Error: Failed to execute query checking DB Stats. Error:" + result.Error.Error())
	} else if result.RowsAffected == 0 {
		log.Fatal("INFO: Query executed, 0 rows")
	} else {
		serverInfoJSON, jsonErr := json.MarshalIndent(serverInfoMap, "", "   ")
		if jsonErr != nil {
			log.Fatal("ERROR: Failed to marshal server info map into JSON" + jsonErr.Error())
		} else {
			log.Printf("INFO: Database Stats:\n%s", string(serverInfoJSON))
		}
	}

	// Start DB Migration if needed
	domainsToMigrate := domain.GetRegisteredModels()
	if len(domainsToMigrate) > 0 {
		log.Printf("Found %v number of domains to migrate, starting migration now..", len(domainsToMigrate))
		err := db.AutoMigrate(domainsToMigrate...)
		if err != nil {
			log.Fatalf("Failed to auto-migrate database: %v", err)
		}
		log.Println("Database auto-migrate successfull..")
	} else {
		log.Println("No Models registered for migration.")
	}

	// Setup GIN
	gin.SetMode(apiConfig.GinMode)
	// New Gin Router
	router := gin.New()
	// Health Check Route
	healthH := handler.NewHealthHandler(db)
	router.GET("/health", healthH.Check)

	// Start the HTTP Server
	serverAddr := fmt.Sprintf(":%s", apiConfig.APIPort)
	log.Printf("API Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start Gin server: %s", err)
	}
}
