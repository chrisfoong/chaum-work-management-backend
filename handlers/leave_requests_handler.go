package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// LEAVE_REQUEST
// Suggested MongoDB collection: "leave_requests"
//
// This file is a route-only skeleton (no class modeling yet - see
// README.md). Handlers bind/return plain maps instead of typed structs,
// and don't touch the database yet. Swap in config.Collection("leave_requests")
// + real bson.M queries once the data shape for LEAVE_REQUEST is settled.
// ---------------------------------------------------------------------------

// ListLeaveRequests handles GET /leave-requests
func ListLeaveRequests(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: list leave request documents from the \"leave_requests\" collection",
		"data":    []gin.H{},
	})
}

// GetLeaveRequest handles GET /leave-requests/:id
func GetLeaveRequest(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: fetch one leave request by id from \"leave_requests\"",
		"id":      id,
	})
}

// CreateLeaveRequest handles POST /leave-requests
func CreateLeaveRequest(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "TODO: insert a new leave request into \"leave_requests\"",
		"payload": body,
	})
}

// UpdateLeaveRequest handles PUT /leave-requests/:id
func UpdateLeaveRequest(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: update leave request " + id + " in \"leave_requests\"",
		"payload": body,
	})
}

// DeleteLeaveRequest handles DELETE /leave-requests/:id
func DeleteLeaveRequest(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: delete leave request " + id + " from \"leave_requests\"",
	})
}
