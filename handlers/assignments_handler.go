package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// TOR_LOCATION_ASSIGNMENT
// Suggested MongoDB collection: "tor_location_assignments"
//
// This file is a route-only skeleton (no class modeling yet - see
// README.md). Handlers bind/return plain maps instead of typed structs,
// and don't touch the database yet. Swap in config.Collection("tor_location_assignments")
// + real bson.M queries once the data shape for TOR_LOCATION_ASSIGNMENT is settled.
// ---------------------------------------------------------------------------

// ListAssignments handles GET /assignments
func ListAssignments(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: list assignment documents from the \"tor_location_assignments\" collection",
		"data":    []gin.H{},
	})
}

// GetAssignment handles GET /assignments/:id
func GetAssignment(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: fetch one assignment by id from \"tor_location_assignments\"",
		"id":      id,
	})
}

// CreateAssignment handles POST /assignments
func CreateAssignment(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "TODO: insert a new assignment into \"tor_location_assignments\"",
		"payload": body,
	})
}

// UpdateAssignment handles PUT /assignments/:id
func UpdateAssignment(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: update assignment " + id + " in \"tor_location_assignments\"",
		"payload": body,
	})
}

// DeleteAssignment handles DELETE /assignments/:id
func DeleteAssignment(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: delete assignment " + id + " from \"tor_location_assignments\"",
	})
}

// GetAssignmentSchedules handles GET /assignments/:id/schedules
// ER relationship: GENERATES: one TOR_LOCATION_ASSIGNMENT generates many WORK_SCHEDULE entries
func GetAssignmentSchedules(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: GENERATES: one TOR_LOCATION_ASSIGNMENT generates many WORK_SCHEDULE entries",
		"id":      id,
		"data":    []gin.H{},
	})
}

// GetAssignmentEvidence handles GET /assignments/:id/evidence
// ER relationship: COLLECTS: one TOR_LOCATION_ASSIGNMENT collects many WORK_EVIDENCE records
func GetAssignmentEvidence(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: COLLECTS: one TOR_LOCATION_ASSIGNMENT collects many WORK_EVIDENCE records",
		"id":      id,
		"data":    []gin.H{},
	})
}

// GetAssignmentEquipmentRequisitions handles GET /assignments/:id/equipment-requisitions
// ER relationship: REQUIRES: one TOR_LOCATION_ASSIGNMENT requires many EQUIPMENT_REQUISITION records
func GetAssignmentEquipmentRequisitions(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: REQUIRES: one TOR_LOCATION_ASSIGNMENT requires many EQUIPMENT_REQUISITION records",
		"id":      id,
		"data":    []gin.H{},
	})
}
