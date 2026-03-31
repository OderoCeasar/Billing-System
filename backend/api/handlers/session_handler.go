package handlers

import (
	"net/http"
	"github.com/OderoCeasar/system/services"
	"github.com/OderoCeasar/system/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SessionHandler struct {
	sessionService *services.SessionService
}

func NewSessionHandler(sessionService *services.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

type CreateSessionRequest struct {
	PackageID string `json:"package_id" binding:"required,uuid"`
	PaymentID string `json:"payment_id" binding:"required,uuid"`
}

// CreateSession creates a new session after successful payment
func (h *SessionHandler) CreateSession(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	packageID, err := uuid.Parse(req.PackageID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid package ID")
		return
	}

	paymentID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	session, err := h.sessionService.CreateSession(userID.(uuid.UUID), packageID, paymentID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Session created successfully", gin.H{
		"session": session,
	})
}

// GetActiveSession gets the user's active session
func (h *SessionHandler) GetActiveSession(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	session, err := h.sessionService.GetActiveSession(userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "No active session found")
		return
	}

	// Get session stats
	stats, err := h.sessionService.GetSessionStats(session.ID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get session stats")
		return
	}

	utils.DataResponse(c, http.StatusOK, gin.H{
		"session": session,
		"stats":   stats,
	})
}

// GetSession gets a specific session by ID
func (h *SessionHandler) GetSession(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	stats, err := h.sessionService.GetSessionStats(sessionID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Session not found")
		return
	}

	utils.DataResponse(c, http.StatusOK, stats)
}

// GetSessionStats gets statistics for a session
func (h *SessionHandler) GetSessionStats(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	stats, err := h.sessionService.GetSessionStats(sessionID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Session not found")
		return
	}

	utils.DataResponse(c, http.StatusOK, stats)
}

// DisconnectSession manually disconnects a session
func (h *SessionHandler) DisconnectSession(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID")
		return
	}

	if err := h.sessionService.DisconnectSession(sessionID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Session disconnected successfully", nil)
}

// ListUserSessions lists all sessions for the authenticated user
func (h *SessionHandler) ListUserSessions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	sessions, err := h.sessionService.ListUserSessions(userID.(uuid.UUID), 50, 0)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch sessions")
		return
	}

	utils.DataResponse(c, http.StatusOK, gin.H{
		"sessions": sessions,
		"count":    len(sessions),
	})
}