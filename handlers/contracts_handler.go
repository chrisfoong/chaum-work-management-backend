package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// CONTRACT_TOR
// Suggested MongoDB collection: "contracts"
//
// This file is a route-only skeleton (no class modeling yet - see
// README.md). Handlers bind/return plain maps instead of typed structs,
// and don't touch the database yet. Swap in config.Collection("contracts")
// + real bson.M queries once the data shape for CONTRACT_TOR is settled.
// ---------------------------------------------------------------------------

// ListContracts handles GET /contracts
func ListContracts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: list contract documents from the \"contracts\" collection",
		"data":    []gin.H{},
	})
}

// GetContract handles GET /contracts/:id
func GetContract(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: fetch one contract by id from \"contracts\"",
		"id":      id,
	})
}

// CreateContract handles POST /contracts
func CreateContract(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "TODO: insert a new contract into \"contracts\"",
		"payload": body,
	})
}

// UpdateContract handles PUT /contracts/:id
func UpdateContract(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: update contract " + id + " in \"contracts\"",
		"payload": body,
	})
}

// DeleteContract handles DELETE /contracts/:id
func DeleteContract(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: delete contract " + id + " from \"contracts\"",
	})
}

// GetContractAssignments handles GET /contracts/:id/assignments
// ER relationship: INCLUDES: one CONTRACT_TOR includes many TOR_LOCATION_ASSIGNMENT records
func GetContractAssignments(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: INCLUDES: one CONTRACT_TOR includes many TOR_LOCATION_ASSIGNMENT records",
		"id":      id,
		"data":    []gin.H{},
	})
}

// GetContractInvoices handles GET /contracts/:id/invoices
// ER relationship: BILLS: one CONTRACT_TOR is billed through many COMPANY_INVOICE records
func GetContractInvoices(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: BILLS: one CONTRACT_TOR is billed through many COMPANY_INVOICE records",
		"id":      id,
		"data":    []gin.H{},
	})
}

// GetContractExpenseClaims handles GET /contracts/:id/expense-claims
// ER relationship: HAS_EXPENSE: one CONTRACT_TOR has many EXPENSE_CLAIM records charged to it
func GetContractExpenseClaims(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: HAS_EXPENSE: one CONTRACT_TOR has many EXPENSE_CLAIM records charged to it",
		"id":      id,
		"data":    []gin.H{},
	})
}
