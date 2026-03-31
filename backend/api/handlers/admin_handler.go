package handlers

import (
	"net/http"

	"github.com/OderoCeasar/system/services"
	"github.com/OderoCeasar/system/utils"
	"github.com/gin-gonic/gin"
)


type AdminHandler struct {
	sessionService		*services.SessionService
}

func NewAdminHandler(sessionService *services.SessionService) *AdminHandler {
	return &AdminHandler{
		sessionService: sessionService,
	}
}


func (h *AdminHandler) ListActiveSessions(c *gin.Context) {
	sessions, err := h.sessionService.ListActiveSessions()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch sessions")
		return
	}

	utils.DataResponse(c, http.StatusOK, gin.H{
		"sessions": sessions,
		"count": len(sessions),
	})
}