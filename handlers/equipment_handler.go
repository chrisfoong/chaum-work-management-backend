package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// EQUIPMENT
// Suggested MongoDB collection: "equipment"
//
// This file is a route-only skeleton (no class modeling yet - see
// README.md). Handlers bind/return plain maps instead of typed structs,
// and don't touch the database yet. Swap in config.Collection("equipment")
// + real bson.M queries once the data shape for EQUIPMENT is settled.
// ---------------------------------------------------------------------------

// ListEquipmentList handles GET /equipment
func ListEquipmentList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: list equipment item documents from the \"equipment\" collection",
		"data":    []gin.H{},
	})
}

// GetEquipmentItem handles GET /equipment/:id
func GetEquipmentItem(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: fetch one equipment item by id from \"equipment\"",
		"id":      id,
	})
}

// CreateEquipmentItem handles POST /equipment
func CreateEquipmentItem(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "TODO: insert a new equipment item into \"equipment\"",
		"payload": body,
	})
}

// UpdateEquipmentItem handles PUT /equipment/:id
func UpdateEquipmentItem(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: update equipment item " + id + " in \"equipment\"",
		"payload": body,
	})
}

// DeleteEquipmentItem handles DELETE /equipment/:id
func DeleteEquipmentItem(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: delete equipment item " + id + " from \"equipment\"",
	})
}

// GetEquipmentRequisitionHistory handles GET /equipment/:id/requisition-items
// ER relationship: IS_LISTED_IN: one EQUIPMENT item appears on many REQUISITION_ITEM lines
func GetEquipmentRequisitionHistory(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: IS_LISTED_IN: one EQUIPMENT item appears on many REQUISITION_ITEM lines",
		"id":      id,
		"data":    []gin.H{},
	})
}
