package service

import (
	"fmt"

	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChatService handles chat business logic
type ChatService struct {
}

// NewChatService creates a new chat service
func NewChatService() *ChatService {
	return &ChatService{}
}

// GetChatRooms gets chat rooms for a user
func (s *ChatService) GetChatRooms(userID string) ([]dto.ChatRoomResponse, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Get chat rooms where user is a member
	var rooms []models.ChatRoom
	err = database.GetDB().Preload("Event").Preload("Event.Creator").
		Joins("JOIN event_members ON chat_rooms.event_id = event_members.event_id").
		Where("event_members.user_id = ? AND event_members.status = ?", userUUID, models.MemberStatusConfirmed).
		Find(&rooms).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get chat rooms: %w", err)
	}

	// Convert to response DTOs
	responses := make([]dto.ChatRoomResponse, len(rooms))
	for i, room := range rooms {
		responses[i] = s.convertChatRoomToResponse(room)
	}

	return responses, nil
}

// GetMessages gets messages in a chat room
func (s *ChatService) GetMessages(roomID, userID string, page, limit int) ([]dto.ChatMessageResponse, int64, error) {
	// Parse IDs
	roomUUID, err := uuid.Parse(roomID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid room ID: %w", err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid user ID: %w", err)
	}

	// Check if user is a member of the room
	var member models.EventMember
	err = database.GetDB().Joins("JOIN chat_rooms ON event_members.event_id = chat_rooms.event_id").
		Where("chat_rooms.id = ? AND event_members.user_id = ? AND event_members.status = ?", roomUUID, userUUID, models.MemberStatusConfirmed).
		First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf("unauthorized")
		}
		return nil, 0, fmt.Errorf("database error: %w", err)
	}

	// Get total count
	var total int64
	err = database.GetDB().Model(&models.ChatMessage{}).Where("room_id = ?", roomUUID).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count messages: %w", err)
	}

	// Get messages with pagination
	var messages []models.ChatMessage
	offset := (page - 1) * limit
	err = database.GetDB().Preload("Sender").Where("room_id = ?", roomUUID).
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&messages).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get messages: %w", err)
	}

	// Convert to response DTOs
	responses := make([]dto.ChatMessageResponse, len(messages))
	for i, message := range messages {
		responses[i] = s.convertChatMessageToResponse(message)
	}

	return responses, total, nil
}

// SendMessage sends a message in a chat room
func (s *ChatService) SendMessage(roomID, userID string, req dto.SendMessageRequest) (*dto.ChatMessageResponse, error) {
	// Parse IDs
	roomUUID, err := uuid.Parse(roomID)
	if err != nil {
		return nil, fmt.Errorf("invalid room ID: %w", err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Check if user is a member of the room
	var member models.EventMember
	err = database.GetDB().Joins("JOIN chat_rooms ON event_members.event_id = chat_rooms.event_id").
		Where("chat_rooms.id = ? AND event_members.user_id = ? AND event_members.status = ?", roomUUID, userUUID, models.MemberStatusConfirmed).
		First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("unauthorized")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Create message
	message := &models.ChatMessage{
		RoomID:      roomUUID,
		SenderID:    userUUID,
		Body:        &req.Body,
		MessageType: &req.MessageType,
	}

	err = database.GetDB().Create(message).Error
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Load message with sender
	err = database.GetDB().Preload("Sender").Where("id = ?", message.ID).First(message).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load message: %w", err)
	}

	response := s.convertChatMessageToResponse(*message)
	return &response, nil
}

// GetRoomByEventID gets chat room by event ID
func (s *ChatService) GetRoomByEventID(eventID string) (*dto.ChatRoomResponse, error) {
	// Parse event ID
	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return nil, fmt.Errorf("invalid event ID: %w", err)
	}

	// Get chat room
	var room models.ChatRoom
	err = database.GetDB().Preload("Event").Preload("Event.Creator").
		Where("event_id = ?", eventUUID).First(&room).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("chat room not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	response := s.convertChatRoomToResponse(room)
	return &response, nil
}

// GetRoomMembers gets members of a chat room
func (s *ChatService) GetRoomMembers(roomID string) ([]dto.UserResponse, error) {
	// Parse room ID
	roomUUID, err := uuid.Parse(roomID)
	if err != nil {
		return nil, fmt.Errorf("invalid room ID: %w", err)
	}

	// Get members
	var members []models.EventMember
	err = database.GetDB().Preload("User").Preload("User.Profile").
		Joins("JOIN chat_rooms ON event_members.event_id = chat_rooms.event_id").
		Where("chat_rooms.id = ? AND event_members.status = ?", roomUUID, models.MemberStatusConfirmed).
		Find(&members).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get room members: %w", err)
	}

	// Convert to response DTOs
	responses := make([]dto.UserResponse, len(members))
	for i, member := range members {
		if member.User != nil {
			responses[i] = dto.UserResponse{
				ID:          member.User.ID.String(),
				Email:       *member.User.Email,
				DisplayName: member.User.GetDisplayName(),
				Provider:    string(member.User.Provider),
				CreatedAt:   member.User.CreatedAt,
			}

			// Add profile info
			if member.User.Profile != nil {
				responses[i].Profile = &dto.UserProfileResponse{
					ID:            member.User.Profile.ID.String(),
					UserID:        member.User.Profile.UserID.String(),
					Bio:           member.User.Profile.Bio,
					Languages:     member.User.Profile.Languages,
					DateOfBirth:   member.User.Profile.DateOfBirth,
					Gender:        string(*member.User.Profile.Gender),
					JobTitle:      member.User.Profile.JobTitle,
					Smoking:       string(*member.User.Profile.Smoking),
					InterestsNote: member.User.Profile.InterestsNote,
					AvatarURL:     member.User.Profile.AvatarURL,
					HomeLocation:  member.User.Profile.HomeLocation,
					CreatedAt:     member.User.Profile.CreatedAt,
					UpdatedAt:     member.User.Profile.UpdatedAt,
				}
			}
		}
	}

	return responses, nil
}

// Helper function to convert chat room to response DTO
func (s *ChatService) convertChatRoomToResponse(room models.ChatRoom) dto.ChatRoomResponse {
	response := dto.ChatRoomResponse{
		ID:        room.ID.String(),
		EventID:   room.EventID.String(),
		CreatedAt: room.CreatedAt,
	}

	// Add event info
	if room.Event != nil {
		response.Event = &dto.EventResponse{
			ID:            room.Event.ID.String(),
			CreatorID:     room.Event.CreatorID.String(),
			Title:         room.Event.Title,
			Description:   room.Event.Description,
			EventType:     string(room.Event.EventType),
			AddressText:   room.Event.AddressText,
			Lat:           room.Event.Lat,
			Lng:           room.Event.Lng,
			StartAt:       room.Event.StartAt,
			EndAt:         room.Event.EndAt,
			Capacity:      room.Event.Capacity,
			Status:        string(room.Event.Status),
			CoverImageURL: room.Event.CoverImageURL,
			CreatedAt:     room.Event.CreatedAt,
			UpdatedAt:     room.Event.UpdatedAt,
		}

		// Add creator info
		if room.Event.Creator != nil {
			response.Event.Creator = &dto.UserResponse{
				ID:          room.Event.Creator.ID.String(),
				Email:       *room.Event.Creator.Email,
				DisplayName: room.Event.Creator.GetDisplayName(),
				Provider:    string(room.Event.Creator.Provider),
				CreatedAt:   room.Event.Creator.CreatedAt,
			}
		}
	}

	return response
}

// Helper function to convert chat message to response DTO
func (s *ChatService) convertChatMessageToResponse(message models.ChatMessage) dto.ChatMessageResponse {
	response := dto.ChatMessageResponse{
		ID:          message.ID.String(),
		RoomID:      message.RoomID.String(),
		SenderID:    message.SenderID.String(),
		Body:        message.Body,
		MessageType: *message.MessageType,
		CreatedAt:   message.CreatedAt,
	}

	// Add sender info
	if message.Sender != nil {
		response.Sender = &dto.UserResponse{
			ID:          message.Sender.ID.String(),
			Email:       *message.Sender.Email,
			DisplayName: message.Sender.GetDisplayName(),
			Provider:    string(message.Sender.Provider),
			CreatedAt:   message.Sender.CreatedAt,
		}
	}

	return response
}

// GetRooms gets chat rooms for a user
func (s *ChatService) GetRooms(userID string) ([]dto.ChatRoomResponse, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Get user's chat rooms
	var rooms []models.ChatRoom
	err = database.GetDB().
		Preload("Event").
		Preload("Event.Creator").
		Joins("JOIN event_members ON chat_rooms.event_id = event_members.event_id").
		Where("event_members.user_id = ? AND event_members.status = ?", userUUID, models.MemberStatusConfirmed).
		Find(&rooms).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get chat rooms: %w", err)
	}

	// Convert to response DTOs
	responses := make([]dto.ChatRoomResponse, len(rooms))
	for i, room := range rooms {
		responses[i] = s.convertChatRoomToResponse(room)
	}

	return responses, nil
}
