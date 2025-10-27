package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Error codes for machine-readable errors
const (
	ErrCodeValidation              = "VALIDATION_ERROR"
	ErrCodeUnauthorized            = "UNAUTHORIZED"
	ErrCodeForbidden               = "FORBIDDEN"
	ErrCodeNotFound                = "NOT_FOUND"
	ErrCodeConflict                = "CONFLICT"
	ErrCodeTooManyRequests         = "RATE_LIMIT_EXCEEDED"
	ErrCodeBadRequest              = "BAD_REQUEST"
	ErrCodeInternalServer          = "INTERNAL_SERVER_ERROR"
	ErrCodeServiceUnavailable      = "SERVICE_UNAVAILABLE"
	ErrCodeAuthenticationFailed    = "AUTHENTICATION_FAILED"
	ErrCodeInvalidToken            = "INVALID_TOKEN"
	ErrCodeExpiredToken            = "EXPIRED_TOKEN"
	ErrCodeInsufficientPermissions = "INSUFFICIENT_PERMISSIONS"
	ErrCodeResourceExists          = "RESOURCE_ALREADY_EXISTS"
	ErrCodeInvalidInput            = "INVALID_INPUT"
	ErrCodeDatabaseError           = "DATABASE_ERROR"
	ErrCodeExternalService         = "EXTERNAL_SERVICE_ERROR"
)

// Request ID key for context
const RequestIDKey = "request_id"

// APIResponse represents a production-grade API response
type APIResponse struct {
	Success   bool        `json:"success"`
	RequestID string      `json:"request_id"`
	Timestamp string      `json:"timestamp"`
	Code      string      `json:"code,omitempty"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
}

// Meta represents pagination and additional metadata
type Meta struct {
	Page       *int   `json:"page,omitempty"`
	Limit      *int   `json:"limit,omitempty"`
	Total      *int64 `json:"total,omitempty"`
	TotalPages *int   `json:"total_pages,omitempty"`
}

// ValidationError represents a field-level validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// GetRequestID gets or creates a request ID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	// Fallback: create new UUID
	return uuid.New().String()
}

// buildResponse creates a base API response
func buildResponse(c *gin.Context, success bool, code, message string) APIResponse {
	return APIResponse{
		Success:   success,
		RequestID: GetRequestID(c),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Code:      code,
		Message:   message,
	}
}

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	response := buildResponse(c, true, "", message)
	response.Data = data
	c.JSON(statusCode, response)
}

// SuccessWithMetaResponse sends a success response with metadata
func SuccessWithMetaResponse(c *gin.Context, statusCode int, message string, data interface{}, meta *Meta) {
	response := buildResponse(c, true, "", message)
	response.Data = data
	response.Meta = meta
	c.JSON(statusCode, response)
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, code, message string, err error) {
	response := buildResponse(c, false, code, message)

	// In production, don't expose internal error details
	if err != nil && statusCode == http.StatusInternalServerError {
		Logger().WithField("request_id", response.RequestID).
			WithField("error", err.Error()).
			Error("Internal server error")
		// Don't send internal error details to client in production
	} else if err != nil {
		response.Errors = err.Error()
	}

	c.JSON(statusCode, response)
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, message string, errors interface{}) {
	response := buildResponse(c, false, ErrCodeValidation, message)
	response.Errors = errors
	c.JSON(http.StatusBadRequest, response)
}

// NotFoundResponse sends a not found response
func NotFoundResponse(c *gin.Context, message string) {
	response := buildResponse(c, false, ErrCodeNotFound, message)
	c.JSON(http.StatusNotFound, response)
}

// UnauthorizedResponse sends an unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	response := buildResponse(c, false, ErrCodeUnauthorized, message)
	c.JSON(http.StatusUnauthorized, response)
}

// ForbiddenResponse sends a forbidden response
func ForbiddenResponse(c *gin.Context, message string) {
	response := buildResponse(c, false, ErrCodeForbidden, message)
	c.JSON(http.StatusForbidden, response)
}

// InternalServerErrorResponse sends an internal server error response
func InternalServerErrorResponse(c *gin.Context, message string, err error) {
	response := buildResponse(c, false, ErrCodeInternalServer, message)

	// Log the actual error but don't expose to client
	if err != nil {
		Logger().WithField("request_id", response.RequestID).
			WithField("error", err.Error()).
			Error("Internal server error")
	}

	c.JSON(http.StatusInternalServerError, response)
}

// BadRequestResponse sends a bad request response
func BadRequestResponse(c *gin.Context, message string) {
	response := buildResponse(c, false, ErrCodeBadRequest, message)
	c.JSON(http.StatusBadRequest, response)
}

// ConflictResponse sends a conflict response
func ConflictResponse(c *gin.Context, message string) {
	response := buildResponse(c, false, ErrCodeConflict, message)
	c.JSON(http.StatusConflict, response)
}

// TooManyRequestsResponse sends a too many requests response
func TooManyRequestsResponse(c *gin.Context, message string) {
	response := buildResponse(c, false, ErrCodeTooManyRequests, message)
	c.JSON(http.StatusTooManyRequests, response)
}

// ServiceUnavailableResponse sends a service unavailable response
func ServiceUnavailableResponse(c *gin.Context, message string) {
	response := buildResponse(c, false, ErrCodeServiceUnavailable, message)
	c.JSON(http.StatusServiceUnavailable, response)
}

// PaginatedResponse sends a paginated success response
func PaginatedResponse(c *gin.Context, message string, data interface{}, total int64, page, limit int) {
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	meta := &Meta{
		Page:       &page,
		Limit:      &limit,
		Total:      &total,
		TotalPages: &totalPages,
	}

	SuccessWithMetaResponse(c, http.StatusOK, message, data, meta)
}

// Legacy wrapper functions for backward compatibility
func SendSuccessResponse(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusOK, message, data)
}

func SendErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	code := ErrCodeInternalServer
	switch statusCode {
	case http.StatusBadRequest:
		code = ErrCodeBadRequest
	case http.StatusUnauthorized:
		code = ErrCodeUnauthorized
	case http.StatusForbidden:
		code = ErrCodeForbidden
	case http.StatusNotFound:
		code = ErrCodeNotFound
	case http.StatusConflict:
		code = ErrCodeConflict
	case http.StatusTooManyRequests:
		code = ErrCodeTooManyRequests
	}
	ErrorResponse(c, statusCode, code, message, err)
}

func SendValidationErrorResponse(c *gin.Context, message string, errors interface{}) {
	ValidationErrorResponse(c, message, errors)
}

func SendNotFoundResponse(c *gin.Context, message string) {
	NotFoundResponse(c, message)
}

func SendUnauthorizedResponse(c *gin.Context, message string) {
	UnauthorizedResponse(c, message)
}

func SendForbiddenResponse(c *gin.Context, message string) {
	ForbiddenResponse(c, message)
}

func SendInternalServerErrorResponse(c *gin.Context, message string, err error) {
	InternalServerErrorResponse(c, message, err)
}

func SendBadRequestResponse(c *gin.Context, message string) {
	BadRequestResponse(c, message)
}

func SendConflictResponse(c *gin.Context, message string) {
	ConflictResponse(c, message)
}

func SendTooManyRequestsResponse(c *gin.Context, message string) {
	TooManyRequestsResponse(c, message)
}

func SendCreatedResponse(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusCreated, message, data)
}

func SendAcceptedResponse(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusAccepted, message, data)
}

func SendNoContentResponse(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func SendPaginatedResponse(c *gin.Context, data interface{}, total int64, page, limit int) {
	PaginatedResponse(c, "Data retrieved successfully", data, total, page, limit)
}
