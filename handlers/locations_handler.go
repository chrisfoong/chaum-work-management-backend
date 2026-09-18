package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// LOCATION
// Suggested MongoDB collection: "locations"
//
// This file is a route-only skeleton (no class modeling yet - see
// README.md). Handlers bind/return plain maps instead of typed structs,
// and don't touch the database yet. Swap in config.Collection("locations")
// + real bson.M queries once the data shape for LOCATION is settled.
// ---------------------------------------------------------------------------

// ListLocations handles GET /locations
func ListLocations(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: list location documents from the \"locations\" collection",
		"data":    []gin.H{},
	})
}

// GetLocation handles GET /locations/:id
func GetLocation(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: fetch one location by id from \"locations\"",
		"id":      id,
	})
}

// CreateLocation handles POST /locations
func CreateLocation(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "TODO: insert a new location into \"locations\"",
		"payload": body,
	})
}

// UpdateLocation handles PUT /locations/:id
func UpdateLocation(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: update location " + id + " in \"locations\"",
		"payload": body,
	})
}

// DeleteLocation handles DELETE /locations/:id
func DeleteLocation(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: delete location " + id + " from \"locations\"",
	})
}

// GetLocationAssignments handles GET /locations/:id/assignments
// ER relationship: HOSTS: one LOCATION hosts many TOR_LOCATION_ASSIGNMENT records
func GetLocationAssignments(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: HOSTS: one LOCATION hosts many TOR_LOCATION_ASSIGNMENT records",
		"id":      id,
		"data":    []gin.H{},
	})
}
