package handlers

import (
	"net/http"

	"github.com/OderoCeasar/system/mpesa"
	"github.com/OderoCeasar/system/services"
	"github.com/OderoCeasar/system/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	mpesaService	*mpesa.Service
	sessionService  *services.SessionService
}

func NewPaymentHandler(mpesaService *mpesa.Service, sessionService *services.SessionService) *PaymentHandler {
	return &PaymentHandler{
		mpesaService: mpesaService,
		sessionService: sessionService,
	}
}


type InitiatePaymentRequest struct {
	PackageID string	`json:"package_id" binding:"required"`
	PhoneNumber string	`json:"phone_number" binding:"required"`
}


func (h *PaymentHandler) InitiatePayment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req InitiatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	packageID, err := uuid.Parse(req.PackageID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid package ID")
		return
	}

	payment, err := h.mpesaService.InitiatePayment(userID.(uuid.UUID), packageID, req.PhoneNumber)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "STK push sent succesfully", gin.H{
		"payment":  payment,
		"payment_id": payment.ID,
	})
}


func (h *PaymentHandler) HandleCallback(c *gin.Context) {
	var callback mpesa.MpesaCallback
	if err := c.ShouldBindJSON(&callback); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.mpesaService.ProcessCallback(&callback); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ResultCode": 0,
		"ResultDesc": "callback processed successfully",
	})

}


func (h *PaymentHandler) CheckPaymentStatus(c *gin.Context) {
	paymentIDStr := c.Param("id")
	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	payment, err := h.mpesaService.GetPaymentStatus(paymentID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Payment not found")
		return
	}

	utils.DataResponse(c, http.StatusOK, gin.H{
		"payment": payment,
	})

}
