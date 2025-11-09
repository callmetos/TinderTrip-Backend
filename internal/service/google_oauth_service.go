package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/pkg/config"
	"TinderTrip-Backend/pkg/database"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// GoogleOAuthService handles Google OAuth authentication
type GoogleOAuthService struct {
	config *oauth2.Config
}

// GoogleUserInfo represents the user info from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// NewGoogleOAuthService creates a new Google OAuth service
func NewGoogleOAuthService() *GoogleOAuthService {
	cfg := config.AppConfig

	// Use Google OAuth credentials from environment
	clientID := cfg.Google.ClientID
	clientSecret := cfg.Google.ClientSecret
	redirectURL := cfg.Google.RedirectURL

	// Debug log
	fmt.Printf("Google OAuth Config - ClientID: %s, ClientSecret: %s, RedirectURL: %s\n",
		clientID, clientSecret, redirectURL)

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleOAuthService{
		config: config,
	}
}

// GetAuthURL returns the Google OAuth authorization URL
func (s *GoogleOAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCodeForToken exchanges authorization code for access token
func (s *GoogleOAuthService) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	return token, nil
}

// GetUserInfo fetches user information from Google
func (s *GoogleOAuthService) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := s.config.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// CreateOrUpdateUser creates or updates a user from Google user info
func (s *GoogleOAuthService) CreateOrUpdateUser(ctx context.Context, userInfo *GoogleUserInfo) (*models.User, bool, error) {
	// Check if user already exists by Google ID
	var user models.User
	err := database.GetDB().Where("google_id = ? AND provider = ?", userInfo.ID, models.AuthProviderGoogle).First(&user).Error

	isNewUser := false

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// User doesn't exist, create new user
			isNewUser = true

			// Check if display_name already exists
			var existingDisplayName models.User
			err = database.GetDB().Where("display_name = ? AND deleted_at IS NULL", userInfo.Name).First(&existingDisplayName).Error
			if err == nil {
				// Display name already taken, append Google ID suffix to make it unique
				uniqueDisplayName := userInfo.Name + "_" + userInfo.ID[:8]
				userInfo.Name = uniqueDisplayName
			} else if err != gorm.ErrRecordNotFound {
				return nil, false, fmt.Errorf("failed to check display name: %w", err)
			}

			user = models.User{
				Email:         &userInfo.Email,
				Provider:      models.AuthProviderGoogle,
				GoogleID:      &userInfo.ID,
				DisplayName:   &userInfo.Name,
				EmailVerified: true, // Google OAuth users are automatically verified
				LastLoginAt:   &[]time.Time{time.Now()}[0],
			}

			err = database.GetDB().Create(&user).Error
			if err != nil {
				return nil, false, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return nil, false, fmt.Errorf("failed to check user existence: %w", err)
		}
	} else {
		// User exists, update last login and ensure email is verified
		now := time.Now()
		user.LastLoginAt = &now

		// Check if display_name needs to be updated and if it's unique
		if user.DisplayName == nil || *user.DisplayName != userInfo.Name {
			// Check if new display_name already exists (excluding current user)
			var existingDisplayName models.User
			err = database.GetDB().Where("display_name = ? AND id != ? AND deleted_at IS NULL", userInfo.Name, user.ID).First(&existingDisplayName).Error
			if err == nil {
				// Display name already taken, append Google ID suffix to make it unique
				uniqueDisplayName := userInfo.Name + "_" + userInfo.ID[:8]
				user.DisplayName = &uniqueDisplayName
			} else if err != gorm.ErrRecordNotFound {
				return nil, false, fmt.Errorf("failed to check display name: %w", err)
			} else {
				// Display name is available, update it
				user.DisplayName = &userInfo.Name
			}
		}

		user.EmailVerified = true // Ensure Google OAuth users are always verified

		err = database.GetDB().Save(&user).Error
		if err != nil {
			return nil, false, fmt.Errorf("failed to update user: %w", err)
		}
	}

	return &user, isNewUser, nil
}

// ValidateToken validates a Google access token
func (s *GoogleOAuthService) ValidateToken(ctx context.Context, token *oauth2.Token) (bool, error) {
	client := s.config.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return false, fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// RefreshToken refreshes a Google access token
func (s *GoogleOAuthService) RefreshToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	tokenSource := s.config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	return newToken, nil
}

// RevokeToken revokes a Google access token
func (s *GoogleOAuthService) RevokeToken(ctx context.Context, token *oauth2.Token) error {
	client := s.config.Client(ctx, token)

	resp, err := client.Get(fmt.Sprintf("https://oauth2.googleapis.com/revoke?token=%s", token.AccessToken))
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to revoke token: status %d", resp.StatusCode)
	}

	return nil
}
