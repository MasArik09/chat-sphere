package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"chatsphere/pkg/config"
)

// AuthMiddleware authenticates requests containing Bearer tokens.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		} else {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token is required"})
			c.Abort()
			return
		}

		userID, err := ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
