package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// EXPENSE_CLAIM
// Suggested MongoDB collection: "expense_claims"
//
// This file is a route-only skeleton (no class modeling yet - see
// README.md). Handlers bind/return plain maps instead of typed structs,
// and don't touch the database yet. Swap in config.Collection("expense_claims")
// + real bson.M queries once the data shape for EXPENSE_CLAIM is settled.
// ---------------------------------------------------------------------------

// ListExpenseClaims handles GET /expense-claims
func ListExpenseClaims(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: list expense claim documents from the \"expense_claims\" collection",
		"data":    []gin.H{},
	})
}

// GetExpenseClaim handles GET /expense-claims/:id
func GetExpenseClaim(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: fetch one expense claim by id from \"expense_claims\"",
		"id":      id,
	})
}

// CreateExpenseClaim handles POST /expense-claims
func CreateExpenseClaim(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "TODO: insert a new expense claim into \"expense_claims\"",
		"payload": body,
	})
}

// UpdateExpenseClaim handles PUT /expense-claims/:id
func UpdateExpenseClaim(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: update expense claim " + id + " in \"expense_claims\"",
		"payload": body,
	})
}

// DeleteExpenseClaim handles DELETE /expense-claims/:id
func DeleteExpenseClaim(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: delete expense claim " + id + " from \"expense_claims\"",
	})
}
