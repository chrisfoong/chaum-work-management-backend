package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// EQUIPMENT_REQUISITION, plus its REQUISITION_ITEM junction sub-resource.
// Suggested MongoDB collections: "equipment_requisitions", "requisition_items"
//
// REQUISITION_ITEM has no PK of its own on the ER diagram - it's the
// associative entity that resolves the EQUIPMENT <-> EQUIPMENT_REQUISITION
// M:M (IS_LISTED_IN + CONTAINS), carrying request_qty/return_qty. It gets
// no top-level resource of its own; it's only reachable nested under a
// requisition, e.g. POST /equipment-requisitions/:id/items.
//
// This file is a route-only skeleton (no class modeling yet - see
// README.md): handlers bind/return plain maps, nothing touches the
// database yet.
// ---------------------------------------------------------------------------

// ListEquipmentRequisitions handles GET /equipment-requisitions
func ListEquipmentRequisitions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: list documents from the \"equipment_requisitions\" collection",
		"data":    []gin.H{},
	})
}

// GetEquipmentRequisition handles GET /equipment-requisitions/:id
func GetEquipmentRequisition(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: fetch one equipment requisition by id",
		"id":      id,
	})
}

// CreateEquipmentRequisition handles POST /equipment-requisitions
func CreateEquipmentRequisition(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "TODO: insert a new equipment requisition",
		"payload": body,
	})
}

// UpdateEquipmentRequisition handles PUT /equipment-requisitions/:id
func UpdateEquipmentRequisition(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: update equipment requisition " + id,
		"payload": body,
	})
}

// DeleteEquipmentRequisition handles DELETE /equipment-requisitions/:id
func DeleteEquipmentRequisition(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "TODO: delete equipment requisition " + id,
	})
}

// -- REQUISITION_ITEM sub-resource (CONTAINS + IS_LISTED_IN) --

// ListRequisitionItems handles GET /equipment-requisitions/:id/items
func ListRequisitionItems(c *gin.Context) {
	requisitionID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message":        "TODO: CONTAINS: list REQUISITION_ITEM lines for this requisition",
		"requisition_id": requisitionID,
		"data":           []gin.H{},
	})
}

// AddRequisitionItem handles POST /equipment-requisitions/:id/items
// Expects an equipment_id + request_qty in the body (and later return_qty
// once the item is returned).
func AddRequisitionItem(c *gin.Context) {
	requisitionID := c.Param("id")
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":        "TODO: add a REQUISITION_ITEM line (equipment_id + request_qty) to this requisition",
		"requisition_id": requisitionID,
		"payload":        body,
	})
}

// RemoveRequisitionItem handles DELETE /equipment-requisitions/:id/items/:itemId
func RemoveRequisitionItem(c *gin.Context) {
	requisitionID := c.Param("id")
	itemID := c.Param("itemId")
	c.JSON(http.StatusOK, gin.H{
		"message":        "TODO: remove REQUISITION_ITEM line from this requisition",
		"requisition_id": requisitionID,
		"item_id":        itemID,
	})
}
