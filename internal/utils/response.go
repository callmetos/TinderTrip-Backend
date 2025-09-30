package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	response := Response{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	c.JSON(statusCode, response)
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, message string, errors map[string]string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
		Data:    errors,
	})
}

// NotFoundResponse sends a not found response
func NotFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: message,
	})
}

// UnauthorizedResponse sends an unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Message: message,
	})
}

// ForbiddenResponse sends a forbidden response
func ForbiddenResponse(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Message: message,
	})
}

// InternalServerErrorResponse sends an internal server error response
func InternalServerErrorResponse(c *gin.Context, message string, err error) {
	response := Response{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	c.JSON(http.StatusInternalServerError, response)
}

// BadRequestResponse sends a bad request response
func BadRequestResponse(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
	})
}

// ConflictResponse sends a conflict response
func ConflictResponse(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Success: false,
		Message: message,
	})
}

// TooManyRequestsResponse sends a too many requests response
func TooManyRequestsResponse(c *gin.Context, message string) {
	c.JSON(http.StatusTooManyRequests, Response{
		Success: false,
		Message: message,
	})
}

// PaginationResponse represents a paginated response
type PaginationResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// SendPaginatedResponse sends a paginated response
func SendPaginatedResponse(c *gin.Context, data interface{}, total int64, page, limit int) {
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, PaginationResponse{
		Success:    true,
		Message:    "Data retrieved successfully",
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// SendSuccessResponse sends a success response
func SendSuccessResponse(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusOK, message, data)
}

// SendErrorResponse sends an error response
func SendErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	ErrorResponse(c, statusCode, message, err)
}

// SendValidationErrorResponse sends a validation error response
func SendValidationErrorResponse(c *gin.Context, message string, errors map[string]string) {
	ValidationErrorResponse(c, message, errors)
}

// SendNotFoundResponse sends a not found response
func SendNotFoundResponse(c *gin.Context, message string) {
	NotFoundResponse(c, message)
}

// SendUnauthorizedResponse sends an unauthorized response
func SendUnauthorizedResponse(c *gin.Context, message string) {
	UnauthorizedResponse(c, message)
}

// SendForbiddenResponse sends a forbidden response
func SendForbiddenResponse(c *gin.Context, message string) {
	ForbiddenResponse(c, message)
}

// SendInternalServerErrorResponse sends an internal server error response
func SendInternalServerErrorResponse(c *gin.Context, message string, err error) {
	InternalServerErrorResponse(c, message, err)
}

// SendBadRequestResponse sends a bad request response
func SendBadRequestResponse(c *gin.Context, message string) {
	BadRequestResponse(c, message)
}

// SendConflictResponse sends a conflict response
func SendConflictResponse(c *gin.Context, message string) {
	ConflictResponse(c, message)
}

// SendTooManyRequestsResponse sends a too many requests response
func SendTooManyRequestsResponse(c *gin.Context, message string) {
	TooManyRequestsResponse(c, message)
}

// SendCreatedResponse sends a created response
func SendCreatedResponse(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusCreated, message, data)
}

// SendAcceptedResponse sends an accepted response
func SendAcceptedResponse(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusAccepted, message, data)
}

// SendNoContentResponse sends a no content response
func SendNoContentResponse(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
