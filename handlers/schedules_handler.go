package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// WORK_SCHEDULE
// Suggested MongoDB collection: "work_schedules"
//
// This file is a route-only skeleton (no class modeling yet - see
// README.md). Handlers bind/return plain maps instead of typed structs,
// and don't touch the database yet. Swap in config.Collection("work_schedules")
// + real bson.M queries once the data shape for WORK_SCHEDULE is settled.
// ---------------------------------------------------------------------------

// ListSchedules handles GET /schedules
func ListSchedules(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: list schedule documents from the \"work_schedules\" collection",
		"data":    []gin.H{},
	})
}

// GetSchedule handles GET /schedules/:id
func GetSchedule(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: fetch one schedule by id from \"work_schedules\"",
		"id":      id,
	})
}

// CreateSchedule handles POST /schedules
func CreateSchedule(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "TODO: insert a new schedule into \"work_schedules\"",
		"payload": body,
	})
}

// UpdateSchedule handles PUT /schedules/:id
func UpdateSchedule(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: update schedule " + id + " in \"work_schedules\"",
		"payload": body,
	})
}

// DeleteSchedule handles DELETE /schedules/:id
func DeleteSchedule(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: delete schedule " + id + " from \"work_schedules\"",
	})
}

// GetScheduleAttendance handles GET /schedules/:id/attendance
// ER relationship: RECORDS: one WORK_SCHEDULE records one ATTENDANCE entry (1:1)
func GetScheduleAttendance(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: RECORDS: one WORK_SCHEDULE records one ATTENDANCE entry (1:1)",
		"id":      id,
		"data":    []gin.H{},
	})
}
