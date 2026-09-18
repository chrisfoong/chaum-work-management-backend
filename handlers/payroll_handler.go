package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// PAYROLL
// Suggested MongoDB collection: "payroll"
//
// This file is a route-only skeleton (no class modeling yet - see
// README.md). Handlers bind/return plain maps instead of typed structs,
// and don't touch the database yet. Swap in config.Collection("payroll")
// + real bson.M queries once the data shape for PAYROLL is settled.
// ---------------------------------------------------------------------------

// ListPayrollRecords handles GET /payroll
func ListPayrollRecords(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: list payroll documents from the \"payroll\" collection",
		"data":    []gin.H{},
	})
}

// GetPayroll handles GET /payroll/:id
func GetPayroll(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: fetch one payroll by id from \"payroll\"",
		"id":      id,
	})
}

// CreatePayroll handles POST /payroll
func CreatePayroll(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "TODO: insert a new payroll into \"payroll\"",
		"payload": body,
	})
}

// UpdatePayroll handles PUT /payroll/:id
func UpdatePayroll(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: update payroll " + id + " in \"payroll\"",
		"payload": body,
	})
}

// DeletePayroll handles DELETE /payroll/:id
func DeletePayroll(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: delete payroll " + id + " from \"payroll\"",
	})
}

// GetPayrollDeductions handles GET /payroll/:id/deductions
// ER relationship: CLEARS: one PAYROLL run clears many DEDUCTION_TRANSACTION records
func GetPayrollDeductions(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: CLEARS: one PAYROLL run clears many DEDUCTION_TRANSACTION records",
		"id":      id,
		"data":    []gin.H{},
	})
}
