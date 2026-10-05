package contract

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"chrisfoong/chaum-work-management-backend/internal/apperr"
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"chrisfoong/chaum-work-management-backend/internal/httpx"
)

// Handler exposes the 1S operations over HTTP.
type Handler struct {
	svc *Service
}

// NewHandler returns a Handler for svc.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes adds the 1S routes to the authenticated web group (supervisor only).
func RegisterRoutes(web *gin.RouterGroup, h *Handler) {
	sup := web.Group("", auth.RequireRole(auth.RoleSupervisor))
	sup.POST("/contracts/info", h.SubmitContractInfo)
	sup.POST("/contracts/scope", h.SubmitScopeData)
	sup.POST("/contracts/confirm", h.ConfirmContract)
	sup.GET("/locations", h.SearchLocations)
}

// SubmitContractInfo handles POST /contracts/info (uc 1S step 4). Stores nothing.
func (h *Handler) SubmitContractInfo(c *gin.Context) {
	var in ContractInfo
	if err := httpx.BindJSON(c, &in); err != nil {
		httpx.WriteError(c, err)
		return
	}
	if err := h.svc.SubmitContractInfo(c.Request.Context(), in); err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"valid": true})
}

// SubmitScopeData handles POST /contracts/scope (uc 1S step 7). Stores nothing.
func (h *Handler) SubmitScopeData(c *gin.Context) {
	var in Scope
	if err := httpx.BindJSON(c, &in); err != nil {
		httpx.WriteError(c, err)
		return
	}
	if err := h.svc.SubmitScopeData(c.Request.Context(), in); err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"valid": true})
}

// ConfirmContract handles POST /contracts/confirm (uc 1S step 10). The 201 body
// carries the data for the contract detail page.
func (h *Handler) ConfirmContract(c *gin.Context) {
	p, ok := auth.PrincipalFrom(c)
	if !ok {
		httpx.WriteError(c, apperr.Internal(errors.New("confirm contract without principal")))
		return
	}
	var in ConfirmRequest
	if err := httpx.BindJSON(c, &in); err != nil {
		httpx.WriteError(c, err)
		return
	}
	out, err := h.svc.ConfirmContract(c.Request.Context(), p.UserID, in)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// SearchLocations handles GET /locations?q= (user addition for the 1S picker).
func (h *Handler) SearchLocations(c *gin.Context) {
	locations, err := h.svc.SearchLocations(c.Request.Context(), c.Query("q"))
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": locations})
}
