package utils

import "github.com/gin-gonic/gin"


type Response struct {
	Success bool	`json:"success"`
	Message string	`json:"message,omitempty"`
	Data interface{} `json:"data,omitempty"`
	Error string	`json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data: data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, error string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: error,
	})
}

func DataResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data: data,
	})
}