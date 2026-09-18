package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// WORKER (specialization of USER, shares its user_id as PK; (p,e) = not every USER is a WORKER)
// Suggested MongoDB collection: "workers"
//
// This file is a route-only skeleton (no class modeling yet - see
// README.md). Handlers bind/return plain maps instead of typed structs,
// and don't touch the database yet. Swap in config.Collection("workers")
// + real bson.M queries once the data shape for WORKER (specialization of USER, shares its user_id as PK; (p,e) = not every USER is a WORKER) is settled.
// ---------------------------------------------------------------------------

// ListWorkers handles GET /workers
func ListWorkers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: list worker documents from the \"workers\" collection",
		"data":    []gin.H{},
	})
}

// GetWorker handles GET /workers/:id
func GetWorker(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: fetch one worker by id from \"workers\"",
		"id":      id,
	})
}

// CreateWorker handles POST /workers
func CreateWorker(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "TODO: insert a new worker into \"workers\"",
		"payload": body,
	})
}

// UpdateWorker handles PUT /workers/:id
func UpdateWorker(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: update worker " + id + " in \"workers\"",
		"payload": body,
	})
}

// DeleteWorker handles DELETE /workers/:id
func DeleteWorker(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: delete worker " + id + " from \"workers\"",
	})
}

// GetWorkerSchedules handles GET /workers/:id/schedules
// ER relationship: IS_ASSIGNED_TO: one WORKER is assigned to many WORK_SCHEDULE entries
func GetWorkerSchedules(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: IS_ASSIGNED_TO: one WORKER is assigned to many WORK_SCHEDULE entries",
		"id":      id,
		"data":    []gin.H{},
	})
}

// GetWorkerAttendance handles GET /workers/:id/attendance
// ER relationship: LOGS: one WORKER logs many ATTENDANCE records
func GetWorkerAttendance(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: LOGS: one WORKER logs many ATTENDANCE records",
		"id":      id,
		"data":    []gin.H{},
	})
}

// GetWorkerEvidence handles GET /workers/:id/evidence
// ER relationship: CAPTURES: one WORKER submits many WORK_EVIDENCE records
func GetWorkerEvidence(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: CAPTURES: one WORKER submits many WORK_EVIDENCE records",
		"id":      id,
		"data":    []gin.H{},
	})
}

// GetWorkerDeductions handles GET /workers/:id/deductions
// ER relationship: INCURS: one WORKER incurs many DEDUCTION_TRANSACTION records
func GetWorkerDeductions(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: INCURS: one WORKER incurs many DEDUCTION_TRANSACTION records",
		"id":      id,
		"data":    []gin.H{},
	})
}

// GetWorkerLeaveRequests handles GET /workers/:id/leave-requests
// ER relationship: SUBMITS: one WORKER submits many LEAVE_REQUEST records
func GetWorkerLeaveRequests(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: SUBMITS: one WORKER submits many LEAVE_REQUEST records",
		"id":      id,
		"data":    []gin.H{},
	})
}
