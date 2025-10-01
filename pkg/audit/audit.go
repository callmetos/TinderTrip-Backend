package audit

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLogger handles audit logging
type AuditLogger struct {
	db *gorm.DB
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger() *AuditLogger {
	return &AuditLogger{
		db: database.GetDB(),
	}
}

// LogAction logs an action to audit_logs table
func (a *AuditLogger) LogAction(actorUserID *string, entityTable string, entityID *string, action string, beforeData, afterData interface{}) error {
	var beforeJSON, afterJSON *string

	// Convert before data to JSON
	if beforeData != nil {
		beforeBytes, err := json.Marshal(beforeData)
		if err != nil {
			return fmt.Errorf("failed to marshal before data: %w", err)
		}
		beforeStr := string(beforeBytes)
		beforeJSON = &beforeStr
	}

	// Convert after data to JSON
	if afterData != nil {
		afterBytes, err := json.Marshal(afterData)
		if err != nil {
			return fmt.Errorf("failed to marshal after data: %w", err)
		}
		afterStr := string(afterBytes)
		afterJSON = &afterStr
	}

	// Convert string IDs to UUID pointers
	var actorUUID, entityUUID *uuid.UUID
	if actorUserID != nil {
		if uuid, err := uuid.Parse(*actorUserID); err == nil {
			actorUUID = &uuid
		}
	}
	if entityID != nil {
		if uuid, err := uuid.Parse(*entityID); err == nil {
			entityUUID = &uuid
		}
	}

	auditLog := &models.AuditLog{
		ActorUserID: actorUUID,
		EntityTable: entityTable,
		EntityID:    entityUUID,
		Action:      action,
		BeforeData:  beforeJSON,
		AfterData:   afterJSON,
		CreatedAt:   time.Now(),
	}

	return a.db.Create(auditLog).Error
}

// LogCreate logs a create action
func (a *AuditLogger) LogCreate(actorUserID *string, entityTable string, entityID *string, data interface{}) error {
	return a.LogAction(actorUserID, entityTable, entityID, "CREATE", nil, data)
}

// LogUpdate logs an update action
func (a *AuditLogger) LogUpdate(actorUserID *string, entityTable string, entityID *string, beforeData, afterData interface{}) error {
	return a.LogAction(actorUserID, entityTable, entityID, "UPDATE", beforeData, afterData)
}

// LogDelete logs a delete action
func (a *AuditLogger) LogDelete(actorUserID *string, entityTable string, entityID *string, data interface{}) error {
	return a.LogAction(actorUserID, entityTable, entityID, "DELETE", data, nil)
}

// LogLogin logs a login action
func (a *AuditLogger) LogLogin(actorUserID *string, loginData interface{}) error {
	return a.LogAction(actorUserID, "users", actorUserID, "LOGIN", nil, loginData)
}

// LogLogout logs a logout action
func (a *AuditLogger) LogLogout(actorUserID *string) error {
	return a.LogAction(actorUserID, "users", actorUserID, "LOGOUT", nil, nil)
}

// LogPasswordReset logs a password reset action
func (a *AuditLogger) LogPasswordReset(actorUserID *string) error {
	return a.LogAction(actorUserID, "users", actorUserID, "PASSWORD_RESET", nil, nil)
}

// LogEventJoin logs an event join action
func (a *AuditLogger) LogEventJoin(actorUserID *string, eventID string) error {
	return a.LogAction(actorUserID, "events", &eventID, "JOIN", nil, map[string]string{"event_id": eventID})
}

// LogEventLeave logs an event leave action
func (a *AuditLogger) LogEventLeave(actorUserID *string, eventID string) error {
	return a.LogAction(actorUserID, "events", &eventID, "LEAVE", map[string]string{"event_id": eventID}, nil)
}

// LogEventComplete logs an event completion action
func (a *AuditLogger) LogEventComplete(actorUserID *string, eventID string) error {
	return a.LogAction(actorUserID, "events", &eventID, "COMPLETE", nil, map[string]string{"event_id": eventID})
}

// GetAuditLogs retrieves audit logs with pagination
func (a *AuditLogger) GetAuditLogs(userID *string, entityTable *string, action *string, page, limit int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := a.db.Model(&models.AuditLog{})

	// Apply filters
	if userID != nil {
		if uuid, err := uuid.Parse(*userID); err == nil {
			query = query.Where("actor_user_id = ?", uuid)
		}
	}
	if entityTable != nil {
		query = query.Where("entity_table = ?", *entityTable)
	}
	if action != nil {
		query = query.Where("action = ?", *action)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetEntityAuditHistory gets audit history for a specific entity
func (a *AuditLogger) GetEntityAuditHistory(entityTable, entityID string, page, limit int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	// Parse entity ID to UUID
	entityUUID, err := uuid.Parse(entityID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid entity ID: %w", err)
	}

	query := a.db.Model(&models.AuditLog{}).Where("entity_table = ? AND entity_id = ?", entityTable, entityUUID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// CompareStructs compares two structs and returns the differences
func CompareStructs(before, after interface{}) map[string]interface{} {
	beforeValue := reflect.ValueOf(before)
	afterValue := reflect.ValueOf(after)

	if beforeValue.Kind() != reflect.Ptr || afterValue.Kind() != reflect.Ptr {
		return nil
	}

	beforeElem := beforeValue.Elem()
	afterElem := afterValue.Elem()

	if beforeElem.Type() != afterElem.Type() {
		return nil
	}

	changes := make(map[string]interface{})

	for i := 0; i < beforeElem.NumField(); i++ {
		field := beforeElem.Type().Field(i)
		beforeField := beforeElem.Field(i)
		afterField := afterElem.Field(i)

		// Skip unexported fields
		if !beforeField.CanInterface() || !afterField.CanInterface() {
			continue
		}

		// Compare field values
		if !reflect.DeepEqual(beforeField.Interface(), afterField.Interface()) {
			changes[field.Name] = map[string]interface{}{
				"before": beforeField.Interface(),
				"after":  afterField.Interface(),
			}
		}
	}

	return changes
}
