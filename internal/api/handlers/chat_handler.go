package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/service"
	"TinderTrip-Backend/internal/utils"

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
// @Description Send a message to a specific chat room (supports JSON and multipart/form-data for images/files)
// @Tags chat
// @Security BearerAuth
// @Accept json,mpfd
// @Produce json
// @Param id path string true "Room ID"
// @Param request body dto.SendMessageRequest true "Message data (JSON)"
// @Param room_id formData string false "Room ID (multipart)"
// @Param body formData string false "Message body (multipart)"
// @Param message_type formData string true "Message type: text, image, or file (multipart)"
// @Param file formData file false "Image or file to upload (multipart)"
// @Success 201 {object} dto.ChatMessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /chat/rooms/{id}/messages [post]
func (h *ChatHandler) SendMessage(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		utils.BadRequestResponse(c, "Room ID is required")
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Check content type to determine if it's multipart or JSON
	contentType := c.GetHeader("Content-Type")
	var req dto.SendMessageRequest
	var imageURL, fileURL *string

	// Check if request has file upload (multipart form data)
	_, fileHeader, _ := c.Request.FormFile("file")
	hasFile := c.Request.MultipartForm != nil || fileHeader != nil

	if hasFile || strings.Contains(contentType, "multipart/form-data") {
		// Handle multipart form data (for image/file uploads)
		req.RoomID = roomID
		req.Body = c.PostForm("body")
		req.MessageType = c.PostForm("message_type")

		if req.MessageType == "" {
			utils.BadRequestResponse(c, "message_type is required")
			return
		}

		// Handle file upload
		file, err := c.FormFile("file")
		if err == nil && file != nil {
			// Initialize file service
			fs, err := service.NewFileService()
			if err != nil {
				utils.InternalServerErrorResponse(c, "Storage initialization failed", err)
				return
			}

			// Open file
			src, err := file.Open()
			if err != nil {
				utils.BadRequestResponse(c, "Invalid file")
				return
			}
			defer src.Close()

			// Determine upload folder based on message type
			folder := "chat_images"
			if req.MessageType == "file" {
				folder = "chat_files"
			}

			// Upload file
			_, url, _, _, _, err := fs.UploadImage(c, folder, file.Filename, src)
			if err != nil {
				utils.BadRequestResponse(c, "Failed to upload file: "+err.Error())
				return
			}

			// Set URL based on message type
			if req.MessageType == "image" {
				imageURL = &url
			} else if req.MessageType == "file" {
				fileURL = &url
			}
		} else if req.MessageType == "image" || req.MessageType == "file" {
			// File is required for image/file message types
			utils.BadRequestResponse(c, "File is required for "+req.MessageType+" message type")
			return
		}
	} else {
		// Handle JSON request - only bind if Content-Type is explicitly JSON
		// If Content-Type is empty or not multipart, try to bind JSON
		if strings.Contains(contentType, "application/json") || contentType == "" {
			if err := c.ShouldBindJSON(&req); err != nil {
				// If binding fails and no Content-Type, try to parse as form data
				if contentType == "" {
					// Fallback: try to get from form
					req.RoomID = roomID
					req.Body = c.PostForm("body")
					req.MessageType = c.PostForm("message_type")

					if req.MessageType == "" {
						utils.BadRequestResponse(c, "message_type is required")
						return
					}
				} else {
					utils.ValidationErrorResponse(c, "Invalid request", err.Error())
					return
				}
			} else {
				req.RoomID = roomID
			}
		} else {
			// Try form data parsing as fallback
			req.RoomID = roomID
			req.Body = c.PostForm("body")
			req.MessageType = c.PostForm("message_type")

			if req.MessageType == "" {
				utils.BadRequestResponse(c, "message_type is required")
				return
			}
		}
	}

	// Validate message type
	if req.MessageType != "text" && req.MessageType != "image" && req.MessageType != "file" {
		utils.BadRequestResponse(c, "Invalid message_type. Must be: text, image, or file")
		return
	}

	// Send message
	var message *dto.ChatMessageResponse
	var err error
	if imageURL != nil || fileURL != nil {
		message, err = h.chatService.SendMessageWithMedia(roomID, userID, req, imageURL, fileURL)
	} else {
		message, err = h.chatService.SendMessage(roomID, userID, req)
	}

	if err != nil {
		if err.Error() == "room not found" {
			utils.NotFoundResponse(c, "Room not found")
			return
		}
		if err.Error() == "unauthorized" {
			utils.ForbiddenResponse(c, "You don't have access to this room")
			return
		}
		if strings.Contains(err.Error(), "required for") {
			utils.BadRequestResponse(c, err.Error())
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to send message", err)
		return
	}

	c.JSON(http.StatusCreated, message)
}
