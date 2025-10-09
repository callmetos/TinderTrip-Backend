package handlers

import (
	"net/http"
	"strconv"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/service"

	"github.com/gin-gonic/gin"
)

// ChatHandler handles chat-related requests
type ChatHandler struct {
	chatService *service.ChatService
}

// NewChatHandler creates a new chat handler
func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		chatService: service.NewChatService(),
	}
}

// GetRooms gets chat rooms for the current user
// @Summary Get chat rooms
// @Description Get chat rooms for the current user
// @Tags chat
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.ChatRoomListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /chat/rooms [get]
func (h *ChatHandler) GetRooms(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get chat rooms
	rooms, err := h.chatService.GetChatRooms(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to get chat rooms",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ChatRoomListResponse{
		Rooms: rooms,
	})
}

// GetMessages gets messages from a chat room
// @Summary Get chat messages
// @Description Get messages from a specific chat room
// @Tags chat
// @Security BearerAuth
// @Produce json
// @Param id path string true "Room ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} dto.ChatMessageListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /chat/rooms/{id}/messages [get]
func (h *ChatHandler) GetMessages(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid room ID",
			Message: "Room ID is required",
		})
		return
	}

	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Get messages
	messages, total, err := h.chatService.GetMessages(roomID, userID, page, limit)
	if err != nil {
		if err.Error() == "room not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Room not found",
				Message: err.Error(),
			})
			return
		}
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "Unauthorized",
				Message: "You don't have access to this room",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to get messages",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ChatMessageListResponse{
		Messages: messages,
		Total:    total,
		Page:     page,
		Limit:    limit,
	})
}

// SendMessage sends a message to a chat room
// @Summary Send chat message
// @Description Send a message to a specific chat room
// @Tags chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Room ID"
// @Param request body dto.SendMessageRequest true "Message data"
// @Success 201 {object} dto.ChatMessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /chat/rooms/{id}/messages [post]
func (h *ChatHandler) SendMessage(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid room ID",
			Message: "Room ID is required",
		})
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Send message
	message, err := h.chatService.SendMessage(roomID, userID, req)
	if err != nil {
		if err.Error() == "room not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Room not found",
				Message: err.Error(),
			})
			return
		}
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "Unauthorized",
				Message: "You don't have access to this room",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to send message",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, message)
}
