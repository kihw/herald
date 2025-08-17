package auth

import (
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

// SessionConfig holds session configuration
type SessionConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	SessionSecret string
}

// SetupSessions configures session middleware for Gin
func SetupSessions(r *gin.Engine) error {
	config := SessionConfig{
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),
		SessionSecret: getEnv("SESSION_SECRET", "your-super-secret-session-key"),
	}

	store, err := redis.NewStore(10, "tcp", config.RedisAddr, config.RedisPassword, config.SessionSecret)
	if err != nil {
		return err
	}

	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	r.Use(sessions.Sessions("lol-session", store))
	return nil
}

// GetSession returns the current session
func GetSession(c *gin.Context) sessions.Session {
	return sessions.Default(c)
}

// Helper functions for environment variables
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		// Simple conversion, could add error handling
		switch value {
		case "0":
			return 0
		case "1":
			return 1
		case "2":
			return 2
		default:
			return fallback
		}
	}
	return fallback
}
