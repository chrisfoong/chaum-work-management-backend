package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// ATTENDANCE
// Suggested MongoDB collection: "attendance"
//
// This file is a route-only skeleton (no class modeling yet - see
// README.md). Handlers bind/return plain maps instead of typed structs,
// and don't touch the database yet. Swap in config.Collection("attendance")
// + real bson.M queries once the data shape for ATTENDANCE is settled.
// ---------------------------------------------------------------------------

// ListAttendanceRecords handles GET /attendance
func ListAttendanceRecords(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: list attendance documents from the \"attendance\" collection",
		"data":    []gin.H{},
	})
}

// GetAttendance handles GET /attendance/:id
func GetAttendance(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: fetch one attendance by id from \"attendance\"",
		"id":      id,
	})
}

// CreateAttendance handles POST /attendance
func CreateAttendance(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "TODO: insert a new attendance into \"attendance\"",
		"payload": body,
	})
}

// UpdateAttendance handles PUT /attendance/:id
func UpdateAttendance(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: update attendance " + id + " in \"attendance\"",
		"payload": body,
	})
}

// DeleteAttendance handles DELETE /attendance/:id
func DeleteAttendance(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: delete attendance " + id + " from \"attendance\"",
	})
}

// GetAttendanceDeductions handles GET /attendance/:id/deductions
// ER relationship: TRIGGERS: one ATTENDANCE record can trigger many DEDUCTION_TRANSACTION records
func GetAttendanceDeductions(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: TRIGGERS: one ATTENDANCE record can trigger many DEDUCTION_TRANSACTION records",
		"id":      id,
		"data":    []gin.H{},
	})
}
