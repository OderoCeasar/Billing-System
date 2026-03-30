package handlers

import (
	"net/http"
	"github.com/OderoCeasar/system/services"
	"github.com/OderoCeasar/system/utils"

	"github.com/gin-gonic/gin"
)

type RADIUSHandler struct {
	radiusService *services.RADIUSService
}

func NewRADIUSHandler(radiusService *services.RADIUSService) *RADIUSHandler {
	return &RADIUSHandler{radiusService: radiusService}
}

type AuthenticateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AccountingStartRequest struct {
	SessionID   string `json:"session_id" binding:"required"`
	Username    string `json:"username" binding:"required"`
	NASIPAddress string `json:"nas_ip_address"`
	NASPortID   string `json:"nas_port_id"`
	IPAddress   string `json:"ip_address"`
	MACAddress  string `json:"mac_address"`
}

type AccountingUpdateRequest struct {
	SessionID     string `json:"session_id" binding:"required"`
	SessionTime   int64  `json:"session_time"`
	InputOctets   int64  `json:"input_octets"`
	OutputOctets  int64  `json:"output_octets"`
}

type AccountingStopRequest struct {
	SessionID     string `json:"session_id" binding:"required"`
	SessionTime   int64  `json:"session_time"`
	InputOctets   int64  `json:"input_octets"`
	OutputOctets  int64  `json:"output_octets"`
}

// Authenticate handles RADIUS authentication requests
func (h *RADIUSHandler) Authenticate(c *gin.Context) {
	var req AuthenticateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	authorized, user, err := h.radiusService.AuthenticateUser(req.Username, req.Password)
	if err != nil || !authorized {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication failed")
		return
	}

	// Get authorization attributes
	attributes, err := h.radiusService.CheckAuthorization(req.Username)
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Authentication successful", gin.H{
		"user":       user,
		"authorized": true,
		"attributes": attributes,
	})
}

// AccountingStart handles RADIUS accounting start
func (h *RADIUSHandler) AccountingStart(c *gin.Context) {
	var req AccountingStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.radiusService.AccountingStart(
		req.SessionID,
		req.Username,
		req.NASIPAddress,
		req.NASPortID,
		req.IPAddress,
		req.MACAddress,
	)

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Accounting start recorded", nil)
}

// AccountingUpdate handles RADIUS accounting interim updates
func (h *RADIUSHandler) AccountingUpdate(c *gin.Context) {
	var req AccountingUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.radiusService.AccountingUpdate(
		req.SessionID,
		req.SessionTime,
		req.InputOctets,
		req.OutputOctets,
	)

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Accounting update recorded", nil)
}

// AccountingStop handles RADIUS accounting stop
func (h *RADIUSHandler) AccountingStop(c *gin.Context) {
	var req AccountingStopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.radiusService.AccountingStop(
		req.SessionID,
		req.SessionTime,
		req.InputOctets,
		req.OutputOctets,
	)

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Accounting stop recorded", nil)
}