package handlers

import (
	"strconv"

	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/utils"
	"TinderTrip-Backend/pkg/audit"

	"github.com/gin-gonic/gin"
)

// AuditHandler handles audit log requests
type AuditHandler struct {
	auditLogger *audit.AuditLogger
}

// NewAuditHandler creates a new audit handler
func NewAuditHandler() *AuditHandler {
	return &AuditHandler{
		auditLogger: audit.NewAuditLogger(),
	}
}

// GetAuditLogs gets audit logs with pagination and filters
// @Summary Get audit logs
// @Description Get audit logs with pagination and optional filters
// @Tags audit
// @Accept json
// @Produce json
// @Param user_id query string false "Filter by user ID"
// @Param entity_table query string false "Filter by entity table"
// @Param action query string false "Filter by action"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dto.AuditLogListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /audit/logs [get]
func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	// Get query parameters
	userID := c.Query("user_id")
	entityTable := c.Query("entity_table")
	action := c.Query("action")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Prepare filter parameters
	var userIDPtr, entityTablePtr, actionPtr *string
	if userID != "" {
		userIDPtr = &userID
	}
	if entityTable != "" {
		entityTablePtr = &entityTable
	}
	if action != "" {
		actionPtr = &action
	}

	// Get audit logs
	logs, total, err := h.auditLogger.GetAuditLogs(userIDPtr, entityTablePtr, actionPtr, page, limit)
	if err != nil {
		utils.SendInternalServerErrorResponse(c, "Failed to get audit logs", err)
		return
	}

	// Convert to response format
	response := dto.AuditLogListResponse{
		Logs:       make([]dto.AuditLogResponse, len(logs)),
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	for i, log := range logs {
		var actorUserID, entityID *string
		if log.ActorUserID != nil {
			actorUserIDStr := log.ActorUserID.String()
			actorUserID = &actorUserIDStr
		}
		if log.EntityID != nil {
			entityIDStr := log.EntityID.String()
			entityID = &entityIDStr
		}

		response.Logs[i] = dto.AuditLogResponse{
			ID:          log.ID.String(),
			ActorUserID: actorUserID,
			EntityTable: log.EntityTable,
			EntityID:    entityID,
			Action:      log.Action,
			BeforeData:  log.BeforeData,
			AfterData:   log.AfterData,
			CreatedAt:   log.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	utils.SendPaginatedResponse(c, response.Logs, total, page, limit)
}

// GetEntityAuditHistory gets audit history for a specific entity
// @Summary Get entity audit history
// @Description Get audit history for a specific entity
// @Tags audit
// @Accept json
// @Produce json
// @Param entity_table path string true "Entity table name"
// @Param entity_id path string true "Entity ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dto.AuditLogListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /audit/entities/{entity_table}/{entity_id} [get]
func (h *AuditHandler) GetEntityAuditHistory(c *gin.Context) {
	entityTable := c.Param("entity_table")
	entityID := c.Param("entity_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get entity audit history
	logs, total, err := h.auditLogger.GetEntityAuditHistory(entityTable, entityID, page, limit)
	if err != nil {
		utils.SendInternalServerErrorResponse(c, "Failed to get entity audit history", err)
		return
	}

	// Convert to response format
	response := dto.AuditLogListResponse{
		Logs:       make([]dto.AuditLogResponse, len(logs)),
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	for i, log := range logs {
		var actorUserID, entityID *string
		if log.ActorUserID != nil {
			actorUserIDStr := log.ActorUserID.String()
			actorUserID = &actorUserIDStr
		}
		if log.EntityID != nil {
			entityIDStr := log.EntityID.String()
			entityID = &entityIDStr
		}

		response.Logs[i] = dto.AuditLogResponse{
			ID:          log.ID.String(),
			ActorUserID: actorUserID,
			EntityTable: log.EntityTable,
			EntityID:    entityID,
			Action:      log.Action,
			BeforeData:  log.BeforeData,
			AfterData:   log.AfterData,
			CreatedAt:   log.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	utils.SendPaginatedResponse(c, response.Logs, total, page, limit)
}
