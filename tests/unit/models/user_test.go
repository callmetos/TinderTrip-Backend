package models_test

import (
	"testing"

	"TinderTrip-Backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUser_IsPasswordAuth(t *testing.T) {
	tests := []struct {
		name     string
		provider models.AuthProvider
		want     bool
	}{
		{
			name:     "Password provider",
			provider: models.AuthProviderPassword,
			want:     true,
		},
		{
			name:     "Google provider",
			provider: models.AuthProviderGoogle,
			want:     false,
		},
		{
			name:     "Apple provider",
			provider: models.AuthProviderApple,
			want:     false,
		},
		{
			name:     "Facebook provider",
			provider: models.AuthProviderFacebook,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &models.User{Provider: tt.provider}
			assert.Equal(t, tt.want, user.IsPasswordAuth())
		})
	}
}

func TestUser_IsGoogleAuth(t *testing.T) {
	tests := []struct {
		name     string
		provider models.AuthProvider
		want     bool
	}{
		{
			name:     "Google provider",
			provider: models.AuthProviderGoogle,
			want:     true,
		},
		{
			name:     "Password provider",
			provider: models.AuthProviderPassword,
			want:     false,
		},
		{
			name:     "Apple provider",
			provider: models.AuthProviderApple,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &models.User{Provider: tt.provider}
			assert.Equal(t, tt.want, user.IsGoogleAuth())
		})
	}
}

func TestUser_HasPassword(t *testing.T) {
	password := "hashed_password"
	emptyPassword := ""

	tests := []struct {
		name         string
		passwordHash *string
		want         bool
	}{
		{
			name:         "Has password",
			passwordHash: &password,
			want:         true,
		},
		{
			name:         "No password",
			passwordHash: nil,
			want:         false,
		},
		{
			name:         "Empty password",
			passwordHash: &emptyPassword,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &models.User{PasswordHash: tt.passwordHash}
			assert.Equal(t, tt.want, user.HasPassword())
		})
	}
}

func TestUser_GetDisplayName(t *testing.T) {
	displayName := "John Doe"
	email := "john@example.com"

	tests := []struct {
		name        string
		displayName *string
		email       *string
		want        string
	}{
		{
			name:        "Has display name",
			displayName: &displayName,
			email:       &email,
			want:        "John Doe",
		},
		{
			name:        "No display name, has email",
			displayName: nil,
			email:       &email,
			want:        "john@example.com",
		},
		{
			name:        "No display name, no email",
			displayName: nil,
			email:       nil,
			want:        "Unknown User",
		},
		{
			name:        "Empty display name, has email",
			displayName: new(string),
			email:       &email,
			want:        "john@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &models.User{
				DisplayName: tt.displayName,
				Email:       tt.email,
			}
			assert.Equal(t, tt.want, user.GetDisplayName())
		})
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	t.Run("Generates UUID if not set", func(t *testing.T) {
		user := &models.User{}
		err := user.BeforeCreate(nil)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, user.ID)
	})

	t.Run("Does not override existing UUID", func(t *testing.T) {
		existingID := uuid.New()
		user := &models.User{ID: existingID}
		err := user.BeforeCreate(nil)

		assert.NoError(t, err)
		assert.Equal(t, existingID, user.ID)
	})
}

func TestUser_TableName(t *testing.T) {
	user := models.User{}
	assert.Equal(t, "users", user.TableName())
}

func TestAuthProvider_Constants(t *testing.T) {
	assert.Equal(t, models.AuthProvider("password"), models.AuthProviderPassword)
	assert.Equal(t, models.AuthProvider("google"), models.AuthProviderGoogle)
	assert.Equal(t, models.AuthProvider("apple"), models.AuthProviderApple)
	assert.Equal(t, models.AuthProvider("facebook"), models.AuthProviderFacebook)
}

func TestUser_CompleteScenario(t *testing.T) {
	t.Run("Password user complete scenario", func(t *testing.T) {
		email := "test@example.com"
		displayName := "Test User"
		password := "hashed_password"

		user := &models.User{
			Email:         &email,
			DisplayName:   &displayName,
			Provider:      models.AuthProviderPassword,
			PasswordHash:  &password,
			EmailVerified: true,
		}

		assert.True(t, user.IsPasswordAuth())
		assert.False(t, user.IsGoogleAuth())
		assert.True(t, user.HasPassword())
		assert.Equal(t, "Test User", user.GetDisplayName())
		assert.Equal(t, "users", user.TableName())
	})

	t.Run("Google user complete scenario", func(t *testing.T) {
		email := "google@example.com"
		googleID := "google-id-123"

		user := &models.User{
			Email:         &email,
			Provider:      models.AuthProviderGoogle,
			GoogleID:      &googleID,
			EmailVerified: true,
		}

		assert.False(t, user.IsPasswordAuth())
		assert.True(t, user.IsGoogleAuth())
		assert.False(t, user.HasPassword())
		assert.Equal(t, "google@example.com", user.GetDisplayName())
	})
}
