package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// WORK_EVIDENCE
// Suggested MongoDB collection: "work_evidence"
//
// This file is a route-only skeleton (no class modeling yet - see
// README.md). Handlers bind/return plain maps instead of typed structs,
// and don't touch the database yet. Swap in config.Collection("work_evidence")
// + real bson.M queries once the data shape for WORK_EVIDENCE is settled.
// ---------------------------------------------------------------------------

// ListEvidenceRecords handles GET /evidence
func ListEvidenceRecords(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: list evidence documents from the \"work_evidence\" collection",
		"data":    []gin.H{},
	})
}

// GetEvidence handles GET /evidence/:id
func GetEvidence(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: fetch one evidence by id from \"work_evidence\"",
		"id":      id,
	})
}

// CreateEvidence handles POST /evidence
func CreateEvidence(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "TODO: insert a new evidence into \"work_evidence\"",
		"payload": body,
	})
}

// UpdateEvidence handles PUT /evidence/:id
func UpdateEvidence(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: update evidence " + id + " in \"work_evidence\"",
		"payload": body,
	})
}

// DeleteEvidence handles DELETE /evidence/:id
func DeleteEvidence(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: delete evidence " + id + " from \"work_evidence\"",
	})
}
