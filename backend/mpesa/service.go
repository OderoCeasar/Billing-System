package mpesa

import (
	"fmt"
	"time"

	"github.com/OderoCeasar/system/config"
	"github.com/OderoCeasar/system/db/models"
	"github.com/OderoCeasar/system/db/repositories"
	"github.com/google/uuid"
)

type Service struct {
	client 		*Client
	config 		*config.MpesaConfig
	paymentRepo	*repositories.PaymentRepository
	packageRepo *repositories.PackageRepository
	userRepo 	*repositories.UserRepository
}

func NewService(
	cfg 			*config.MpesaConfig,
	paymentRepo 	*repositories.PaymentRepository,
	packageRepo		*repositories.PackageRepository,
	userRepo		*repositories.UserRepository,
) *Service {
	return &Service{
		client: 	NewClient(cfg),
		config:		cfg,
		paymentRepo: paymentRepo,
		packageRepo: packageRepo,
		userRepo:    userRepo,
	}
}


// Initiate MPesa STK Push payment
func (s *Service) InitiatePayment(userID, packageID uuid.UUID, phoneNumber string) (*models.Payment, error) {
	pkg, err := s.packageRepo.FindByID(packageID)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}

	_, err = s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	formattedPhone := s.client.FormatPhoneNumber(phoneNumber)

	payment := &models.Payment{
		UserID: 	userID,
		PackageID:  packageID,
		Amount: 	pkg.Price,
		PhoneNumber: formattedPhone,
		Status: 	models.PaymentStatusPending,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	timestamp := s.client.GenerateTimestamp()
	password := s.client.GeneratePassword(timestamp)

	stkRequest := &STKPushRequest{
		BusinessShortCode: 		s.config.ShortCode,
		Password: 				password,
		Timestamp: 				timestamp,
		TransactionType: 		"CustomerPayBillOnline",
		Amount: 				fmt.Sprintf("%.0f", pkg.Price),
		PartyA: 				formattedPhone,
		PartyB: 				s.config.ShortCode,
		PhoneNumber: 			formattedPhone,
		CallbackURL: 			s.config.CallbackURL,
		AccountReference: 		payment.ID.String(),
		TransactionDesc: 		fmt.Sprintf("wifi Package: %s", pkg.Name),
	}

	stkResp, err := s.client.InitiateSTKPush(stkRequest)
	if err != nil {
		payment.Status = models.PaymentStatusFailed
		payment.ResultDesc = err.Error()
		s.paymentRepo.Update(payment)
		return nil, fmt.Errorf("STK Push failed: %w", err)
	}

	payment.MpesaCheckoutID = stkResp.CheckoutRequestID
	if err := s.paymentRepo.Update(payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}
	return payment, nil
}

func (s *Service) ProcessCallback(callback *MpesaCallback) error {
	checkoutID := callback.Body.StkCallback.CheckoutRequestID
	resultCode := callback.Body.StkCallback.ResultCode
	resultDesc := callback.Body.StkCallback.ResultDesc

	payment, err := s.paymentRepo.FindByCheckoutID(checkoutID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	payment.CallbackReceived = true
	payment.ResultCode = resultCode
	payment.ResultDesc = resultDesc

	if resultCode == 0 {
		payment.Status = models.PaymentStatusCompleted

		for _, item := range callback.Body.StkCallback.CallbackMetadata.Item {
			switch item.Name {
			case "MpesaReceiptNumber":
				if val, ok := item.Value.(string); ok {
					payment.MpesaReceiptNumber = val
				}
			case "TransactionDate":
				if val, ok := item.Value.(float64); ok {
					timestamp := time.Unix(int64(val), 0)
					payment.TransactionDate = &timestamp
				}	
			}
		}
	} else {
		payment.Status = models.PaymentStatusFailed
	}

	return s.paymentRepo.Update(payment)
}

func (s *Service) GetPaymentStatus(paymentID uuid.UUID) (*models.Payment, error) {
	return s.paymentRepo.FindByID(paymentID)
}

func (s *Service) ListUserPayments(userID uuid.UUID, limit, offset int) ([]models.Payment, error) {
	return s.paymentRepo.ListByUser(userID, limit, offset)
}
