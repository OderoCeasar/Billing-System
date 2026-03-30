package middleware

import (
	"net/http"
	"strings"

	"github.com/OderoCeasar/system/services"
	"github.com/OderoCeasar/system/utils"
	"github.com/gin-gonic/gin"
)


func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return 
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization format")
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := authService.ValidateToken(token)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return 
		}

		c.Set("user_id", claims.UserID)
		c.Set("phone_number", claims.PhoneNumber)
		c.Set("is_admin", claims.IsAdmin)
		c.Next()
	}
}


func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			utils.ErrorResponse(c, http.StatusForbidden, "Admin access required")
			c.Abort()
			return 
		}
		c.Next()
	}
}