package middleware

import (
	"fmt"
	"time"

	"github.com/OderoCeasar/system/utils"
	"github.com/gin-gonic/gin"
)


func LoggingMiddleware() gin.HandlerFunc {
	logger := utils.GetLogger()

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		logger.Info(sprintf("[%s] %s %s - %d - %v", time.Now().Format("2006-01-02"), method, path, statusCode, duration))
	}
}


func CORSMiddleware(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentails", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return 
		}
		c.Next()
	}

}


func sprintf(format string, args ...interface{}) string {
	// helper function
	return fmt.Sprintf(format, args...)
}
