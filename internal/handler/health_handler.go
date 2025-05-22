package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Struct to hold dependencies of
// Health check handlers
type HealthHandler struct {
	DB *gorm.DB
}

// Constructor for the the HealthHandler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		DB: db,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	// Create the initial map response
	response := gin.H{
		"api_status": "running",
	}
	// Check DB Stats
	var serverInfoMap map[string]any
	sqlQuery := "SELECT    NOW() as current_db_time,     current_setting('TimeZone') as db_timezone;"
	result := h.DB.Raw(sqlQuery).Scan(&serverInfoMap)

	if result.Error != nil {
		log.Printf("Error: Failed to execute query checking DB Stats. Error: %s", result.Error.Error())
		response["database_status"] = "unhealthy"
		response["database_error"] = result.Error.Error()
		response["api_status"] = "degraded"
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("INFO: Query executed, 0 rows")
		response["database_status"] = "degraded_no_rows"
		response["database_info"] = "Query executed, 0 rows"
		response["api_status"] = "degraded"
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	serverInfoJSON, jsonErr := json.MarshalIndent(serverInfoMap, "", "   ")
	if jsonErr != nil {
		log.Printf("ERROR: Failed to marshal server info map into JSON: %s", jsonErr.Error())
	} else {
		log.Printf("INFO: Database Stats:\n%s", string(serverInfoJSON))
	}
	response["database_status"] = "healthy"
	response["database_info"] = serverInfoMap
	c.JSON(http.StatusOK, response)
}
