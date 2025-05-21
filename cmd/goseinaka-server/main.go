package main

import (
	"encoding/json"
	"log"

	"github.com/ecabigting/goseinaka/internal/config"
	"github.com/ecabigting/goseinaka/internal/database"
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
			log.Printf("DEBUG: Loaded config fallback: %s", apiConfig)
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
	var serverInfoMap map[string]interface{}
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

	log.Println("API at port", apiConfig.APIPort)
}
