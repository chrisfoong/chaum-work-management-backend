package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// DEDUCTION_TRANSACTION
// Suggested MongoDB collection: "deduction_transactions"
//
// This file is a route-only skeleton (no class modeling yet - see
// README.md). Handlers bind/return plain maps instead of typed structs,
// and don't touch the database yet. Swap in config.Collection("deduction_transactions")
// + real bson.M queries once the data shape for DEDUCTION_TRANSACTION is settled.
// ---------------------------------------------------------------------------

// ListDeductions handles GET /deductions
func ListDeductions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: list deduction documents from the \"deduction_transactions\" collection",
		"data":    []gin.H{},
	})
}

// GetDeduction handles GET /deductions/:id
func GetDeduction(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: fetch one deduction by id from \"deduction_transactions\"",
		"id":      id,
	})
}

// CreateDeduction handles POST /deductions
func CreateDeduction(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "TODO: insert a new deduction into \"deduction_transactions\"",
		"payload": body,
	})
}

// UpdateDeduction handles PUT /deductions/:id
func UpdateDeduction(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: update deduction " + id + " in \"deduction_transactions\"",
		"payload": body,
	})
}

// DeleteDeduction handles DELETE /deductions/:id
func DeleteDeduction(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: delete deduction " + id + " from \"deduction_transactions\"",
	})
}
