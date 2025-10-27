package utils_test

import (
	"strings"
	"testing"

	"TinderTrip-Backend/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Hash valid password",
			password: "SecurePass123!",
			wantErr:  false,
		},
		{
			name:     "Hash short password",
			password: "abc",
			wantErr:  false,
		},
		{
			name:     "Hash long password",
			password: strings.Repeat("a", 200),
			wantErr:  false,
		},
		{
			name:     "Hash password with special characters",
			password: "P@ssw0rd!#$%^&*()",
			wantErr:  false,
		},
		{
			name:     "Hash empty password",
			password: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := utils.HashPassword(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.Contains(t, hash, "$argon2id$")

				// Verify that the hash can verify the original password
				valid, err := utils.VerifyPassword(tt.password, hash)
				assert.NoError(t, err)
				assert.True(t, valid)
			}
		})
	}
}

func TestHashPasswordDifferentSalts(t *testing.T) {
	password := "SamePassword123!"

	hash1, err := utils.HashPassword(password)
	require.NoError(t, err)

	hash2, err := utils.HashPassword(password)
	require.NoError(t, err)

	// Hashes should be different due to random salts
	assert.NotEqual(t, hash1, hash2, "Same password should produce different hashes")

	// But both should verify correctly
	valid1, err := utils.VerifyPassword(password, hash1)
	require.NoError(t, err)
	assert.True(t, valid1)

	valid2, err := utils.VerifyPassword(password, hash2)
	require.NoError(t, err)
	assert.True(t, valid2)
}

func TestVerifyPassword(t *testing.T) {
	correctPassword := "CorrectPassword123!"
	hash, err := utils.HashPassword(correctPassword)
	require.NoError(t, err)

	tests := []struct {
		name        string
		password    string
		hash        string
		wantValid   bool
		wantErr     bool
		errContains string
	}{
		{
			name:      "Correct password",
			password:  correctPassword,
			hash:      hash,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "Wrong password",
			password:  "WrongPassword123!",
			hash:      hash,
			wantValid: false,
			wantErr:   false,
		},
		{
			name:        "Invalid hash format",
			password:    correctPassword,
			hash:        "invalid-hash",
			wantValid:   false,
			wantErr:     true,
			errContains: "invalid hash format",
		},
		{
			name:        "Empty hash",
			password:    correctPassword,
			hash:        "",
			wantValid:   false,
			wantErr:     true,
			errContains: "invalid hash format",
		},
		{
			name:        "Wrong algorithm",
			password:    correctPassword,
			hash:        "$bcrypt$v=19$m=65536,t=1,p=4$salt$hash",
			wantValid:   false,
			wantErr:     true,
			errContains: "unsupported algorithm",
		},
		{
			name:      "Empty password with valid hash",
			password:  "",
			hash:      hash,
			wantValid: false,
			wantErr:   false,
		},
		{
			name:        "Malformed parameters",
			password:    correctPassword,
			hash:        "$argon2id$v=19$invalid$salt$hash",
			wantValid:   false,
			wantErr:     true,
			errContains: "invalid parameters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := utils.VerifyPassword(tt.password, tt.hash)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValid, valid)
			}
		})
	}
}

func TestGenerateRandomPassword(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{
			name:    "Generate 8 character password",
			length:  8,
			wantErr: false,
		},
		{
			name:    "Generate 16 character password",
			length:  16,
			wantErr: false,
		},
		{
			name:    "Generate 32 character password",
			length:  32,
			wantErr: false,
		},
		{
			name:    "Generate 1 character password",
			length:  1,
			wantErr: false,
		},
		{
			name:    "Generate 100 character password",
			length:  100,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := utils.GenerateRandomPassword(tt.length)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, password)
			} else {
				assert.NoError(t, err)
				assert.Len(t, password, tt.length)

				// Check that password contains only valid characters
				validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
				for _, char := range password {
					assert.Contains(t, validChars, string(char))
				}
			}
		})
	}
}

func TestGenerateRandomPasswordUniqueness(t *testing.T) {
	length := 16
	password1, err := utils.GenerateRandomPassword(length)
	require.NoError(t, err)

	password2, err := utils.GenerateRandomPassword(length)
	require.NoError(t, err)

	// Passwords should be different (with very high probability)
	assert.NotEqual(t, password1, password2, "Generated passwords should be unique")
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "Strong password",
			password: "SecurePass123!",
			wantErr:  false,
		},
		{
			name:     "Another strong password",
			password: "MyP@ssw0rd!",
			wantErr:  false,
		},
		{
			name:        "Too short",
			password:    "Short1!",
			wantErr:     true,
			errContains: "at least 8 characters",
		},
		{
			name:        "Too long",
			password:    "A1!" + strings.Repeat("a", 130),
			wantErr:     true,
			errContains: "no more than 128 characters",
		},
		{
			name:        "No uppercase",
			password:    "password123!",
			wantErr:     true,
			errContains: "at least one uppercase letter",
		},
		{
			name:        "No lowercase",
			password:    "PASSWORD123!",
			wantErr:     true,
			errContains: "at least one lowercase letter",
		},
		{
			name:        "No digit",
			password:    "Password!",
			wantErr:     true,
			errContains: "at least one digit",
		},
		{
			name:        "No special character",
			password:    "Password123",
			wantErr:     true,
			errContains: "at least one special character",
		},
		{
			name:        "Empty password",
			password:    "",
			wantErr:     true,
			errContains: "at least 8 characters",
		},
		{
			name:     "All special characters supported",
			password: "Pass123!@#$%^&*()_+-=[]{}|;:,.<>?",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidatePasswordStrength(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "TestPassword123!"
	hash, err := utils.HashPassword(password)
	require.NoError(t, err)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "Correct password and hash",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "Wrong password",
			password: "WrongPassword123!",
			hash:     hash,
			want:     false,
		},
		{
			name:     "Invalid hash",
			password: password,
			hash:     "invalid-hash",
			want:     false,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "Empty hash",
			password: password,
			hash:     "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CheckPasswordHash(tt.password, tt.hash)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestPasswordHashAndVerifyIntegration(t *testing.T) {
	passwords := []string{
		"SimplePass123!",
		"C0mpl3x!P@ssw0rd",
		"Long" + strings.Repeat("a", 50) + "Pass123!",
		"Special!@#$%^&*()123Abc",
	}

	for _, password := range passwords {
		t.Run("Integration test for: "+password[:10], func(t *testing.T) {
			// Hash the password
			hash, err := utils.HashPassword(password)
			require.NoError(t, err)
			require.NotEmpty(t, hash)

			// Verify correct password
			valid, err := utils.VerifyPassword(password, hash)
			require.NoError(t, err)
			assert.True(t, valid, "Should verify correct password")

			// Verify wrong password
			valid, err = utils.VerifyPassword(password+"wrong", hash)
			require.NoError(t, err)
			assert.False(t, valid, "Should not verify wrong password")
		})
	}
}

func TestDefaultPasswordConfig(t *testing.T) {
	config := utils.DefaultPasswordConfig()

	assert.NotNil(t, config)
	assert.Equal(t, uint32(1), config.Time)
	assert.Equal(t, uint32(64*1024), config.Memory)
	assert.Equal(t, uint8(4), config.Threads)
	assert.Equal(t, uint32(32), config.KeyLen)
}

func BenchmarkHashPassword(b *testing.B) {
	password := "BenchmarkPassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = utils.HashPassword(password)
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "BenchmarkPassword123!"
	hash, _ := utils.HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = utils.VerifyPassword(password, hash)
	}
}
