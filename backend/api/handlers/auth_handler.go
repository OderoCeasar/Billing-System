package handlers

import (
	"net/http"

	"github.com/OderoCeasar/system/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	PhoneNumber string  `json:"phone_number" binding:"required"`
	Password    string	`json:"password" binding:"required"`
}

type LoginRequest struct {
	PhoneNumber  string		`json:"phone_number" binding:"required"`
	Password 	 string		`json:"password" binding:"required"`
}

type QuickRegisterRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}


func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	user, err := h.authService.Register(req.PhoneNumber, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"token":	token,
		"user":		user,
	})
}


func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	token, user, err := h.authService.Login(req.PhoneNumber, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":	"login successfully",
		"token":	token,
		"user":		user,
	})
}


func (h *AuthHandler) QuickRegister(c *gin.Context) {
	var req QuickRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	user, token, err := h.authService.QuickRegiter(req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":	"registration successfully",
		"token":	token,
		"user":		user,
	})
}


func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"unathourized"})
		return
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"unathourized"})
		return
	}

	user, err := h.authService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"unathourized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
