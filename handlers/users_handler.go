package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// USER
// Suggested MongoDB collection: "users"
//
// This file is a route-only skeleton (no class modeling yet - see
// README.md). Handlers bind/return plain maps instead of typed structs,
// and don't touch the database yet. Swap in config.Collection("users")
// + real bson.M queries once the data shape for USER is settled.
// ---------------------------------------------------------------------------

// ListUsers handles GET /users
func ListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: list user documents from the \"users\" collection",
		"data":    []gin.H{},
	})
}

// GetUser handles GET /users/:id
func GetUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: fetch one user by id from \"users\"",
		"id":      id,
	})
}

// CreateUser handles POST /users
func CreateUser(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "TODO: insert a new user into \"users\"",
		"payload": body,
	})
}

// UpdateUser handles PUT /users/:id
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: update user " + id + " in \"users\"",
		"payload": body,
	})
}

// DeleteUser handles DELETE /users/:id
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: delete user " + id + " from \"users\"",
	})
}

// GetUserPayroll handles GET /users/:id/payroll
// ER relationship: RECEIVES: one USER receives many PAYROLL records
func GetUserPayroll(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: RECEIVES: one USER receives many PAYROLL records",
		"id":      id,
		"data":    []gin.H{},
	})
}

// GetUserExpenseClaims handles GET /users/:id/expense-claims
// ER relationship: CLAIMS: one USER files many EXPENSE_CLAIM records
func GetUserExpenseClaims(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: CLAIMS: one USER files many EXPENSE_CLAIM records",
		"id":      id,
		"data":    []gin.H{},
	})
}

// GetUserEquipmentRequisitions handles GET /users/:id/equipment-requisitions
// ER relationship: REQUEST: one USER submits many EQUIPMENT_REQUISITION records
func GetUserEquipmentRequisitions(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: REQUEST: one USER submits many EQUIPMENT_REQUISITION records",
		"id":      id,
		"data":    []gin.H{},
	})
}
