package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// COMPANY_INVOICE
// Suggested MongoDB collection: "company_invoices"
//
// This file is a route-only skeleton (no class modeling yet - see
// README.md). Handlers bind/return plain maps instead of typed structs,
// and don't touch the database yet. Swap in config.Collection("company_invoices")
// + real bson.M queries once the data shape for COMPANY_INVOICE is settled.
// ---------------------------------------------------------------------------

// ListInvoices handles GET /invoices
func ListInvoices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: list invoice documents from the \"company_invoices\" collection",
		"data":    []gin.H{},
	})
}

// GetInvoice handles GET /invoices/:id
func GetInvoice(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: fetch one invoice by id from \"company_invoices\"",
		"id":      id,
	})
}

// CreateInvoice handles POST /invoices
func CreateInvoice(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "TODO: insert a new invoice into \"company_invoices\"",
		"payload": body,
	})
}

// UpdateInvoice handles PUT /invoices/:id
func UpdateInvoice(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: update invoice " + id + " in \"company_invoices\"",
		"payload": body,
	})
}

// DeleteInvoice handles DELETE /invoices/:id
func DeleteInvoice(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: delete invoice " + id + " from \"company_invoices\"",
	})
}
