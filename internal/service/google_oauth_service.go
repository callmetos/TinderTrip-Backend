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
func (s *GoogleOAuthService) CreateOrUpdateUser(ctx context.Context, userInfo *GoogleUserInfo) (*models.User, error) {
	// Check if user already exists by Google ID
	var user models.User
	err := database.GetDB().Where("google_id = ? AND provider = ?", userInfo.ID, models.AuthProviderGoogle).First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// User doesn't exist, create new user
			user = models.User{
				Email:       &userInfo.Email,
				Provider:    models.AuthProviderGoogle,
				GoogleID:    &userInfo.ID,
				DisplayName: &userInfo.Name,
				LastLoginAt: &[]time.Time{time.Now()}[0],
			}

			err = database.GetDB().Create(&user).Error
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to check user existence: %w", err)
		}
	} else {
		// User exists, update last login
		now := time.Now()
		user.LastLoginAt = &now
		user.DisplayName = &userInfo.Name

		err = database.GetDB().Save(&user).Error
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	return &user, nil
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
