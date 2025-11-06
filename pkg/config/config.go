package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	JWT        JWTConfig
	Email      EmailConfig
	AWS        AWSConfig
	Firebase   FirebaseConfig
	Google     GoogleConfig
	RateLimit  RateLimitConfig
	CORS       CORSConfig
	Nextcloud  NextcloudConfig
	Monitoring MonitoringConfig
}

type ServerConfig struct {
	Port        string
	Host        string
	Mode        string
	FrontendURL string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFromName string
}

type AWSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	S3Bucket        string
}

type FirebaseConfig struct {
	ProjectID   string
	PrivateKey  string
	ClientEmail string
}

type RateLimitConfig struct {
	Requests int
	Window   string
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type NextcloudConfig struct {
	URL      string
	Username string
	Password string
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type MonitoringConfig struct {
	Enabled     bool
	MetricsPort string
	HealthPort  string
}

var AppConfig *Config

func LoadConfig() {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("No .env file found: %v, using system environment variables", err)
	} else {
		log.Println("Successfully loaded .env file")
	}

	AppConfig = &Config{
		Server: ServerConfig{

			Port:        getEnv("SERVER_PORT", ""),
			Host:        getEnv("SERVER_HOST", ""),
			Mode:        getEnv("GIN_MODE", ""),
			FrontendURL: "http://192.168.1.120:8081",
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnv("DB_PORT", ""),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", ""),
			SSLMode:  getEnv("DB_SSLMODE", ""),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", ""),
			Port:     getEnv("REDIS_PORT", ""),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", -1),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", ""),
			ExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", -1),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", ""),
			SMTPPort:     getEnvAsInt("SMTP_PORT", -1),
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			SMTPFromName: getEnv("SMTP_FROM_NAME", ""),
		},
		AWS: AWSConfig{
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			Region:          getEnv("AWS_REGION", ""),
			S3Bucket:        getEnv("AWS_S3_BUCKET", ""),
		},
		Firebase: FirebaseConfig{
			ProjectID:   getEnv("FIREBASE_PROJECT_ID", ""),
			PrivateKey:  getEnv("FIREBASE_PRIVATE_KEY", ""),
			ClientEmail: getEnv("FIREBASE_CLIENT_EMAIL", ""),
		},
		Google: GoogleConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
		},
		RateLimit: RateLimitConfig{
			Requests: getEnvAsInt("RATE_LIMIT_REQUESTS", -1),
			Window:   getEnv("RATE_LIMIT_WINDOW", ""),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{}),
			AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{}),
			AllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{}),
		},
		Nextcloud: NextcloudConfig{
			URL:      getEnv("NEXTCLOUD_URL", ""),
			Username: getEnv("NEXTCLOUD_USERNAME", ""),
			Password: getEnv("NEXTCLOUD_PASSWORD", ""),
		},
		Monitoring: MonitoringConfig{
			Enabled:     getEnvAsBool("MONITORING_ENABLED", true),
			MetricsPort: getEnv("METRICS_PORT", "9090"),
			HealthPort:  getEnv("HEALTH_PORT", "8080"),
		},
	}

	// Validate required configuration
	validateConfig()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func validateConfig() {
	required := map[string]string{
		"SERVER_PORT":          AppConfig.Server.Port,
		"SERVER_HOST":          AppConfig.Server.Host,
		"DB_HOST":              AppConfig.Database.Host,
		"DB_PORT":              AppConfig.Database.Port,
		"DB_USER":              AppConfig.Database.User,
		"DB_PASSWORD":          AppConfig.Database.Password,
		"DB_NAME":              AppConfig.Database.Name,
		"JWT_SECRET":           AppConfig.JWT.Secret,
		"GOOGLE_CLIENT_ID":     AppConfig.Google.ClientID,
		"GOOGLE_CLIENT_SECRET": AppConfig.Google.ClientSecret,
	}

	var missing []string
	for key, value := range required {
		if value == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		log.Fatalf("Missing required environment variables: %v", missing)
	}

	// Validate numeric values
	if AppConfig.Database.Port == "" {
		log.Fatal("DB_PORT is required")
	}
	if AppConfig.Redis.Port == "" {
		log.Fatal("REDIS_PORT is required")
	}
	if AppConfig.JWT.ExpireHours <= 0 {
		log.Fatal("JWT_EXPIRE_HOURS must be greater than 0")
	}
	if AppConfig.Redis.DB < 0 {
		log.Fatal("REDIS_DB must be a valid number")
	}
	// SMTP_PORT is optional - only validate if email is configured
	if AppConfig.Email.SMTPUsername != "" && AppConfig.Email.SMTPPort <= 0 {
		log.Fatal("SMTP_PORT must be greater than 0 when SMTP_USERNAME is set")
	}
	// RATE_LIMIT_REQUESTS is optional - set default if not provided
	if AppConfig.RateLimit.Requests <= 0 {
		AppConfig.RateLimit.Requests = 100
		log.Println("Using default RATE_LIMIT_REQUESTS: 100")
	}

	// Set default values for optional fields
	if AppConfig.RateLimit.Window == "" {
		AppConfig.RateLimit.Window = "1h"
		log.Println("Using default RATE_LIMIT_WINDOW: 1h")
	}

	// Set default CORS values if not provided
	if len(AppConfig.CORS.AllowedOrigins) == 0 {
		AppConfig.CORS.AllowedOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
		log.Println("Using default CORS_ALLOWED_ORIGINS")
	}
	if len(AppConfig.CORS.AllowedMethods) == 0 {
		AppConfig.CORS.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		log.Println("Using default CORS_ALLOWED_METHODS")
	}
	if len(AppConfig.CORS.AllowedHeaders) == 0 {
		AppConfig.CORS.AllowedHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
		log.Println("Using default CORS_ALLOWED_HEADERS")
	}

	log.Println("Configuration validation passed")
}
