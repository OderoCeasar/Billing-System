package repositories

import (
	"github.com/OderoCeasar/system/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db:db}
}

func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) FindByID(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("user").Preload("Package").Where("id = ?", id).First(&payment).Error
	return &payment, err
}

func (r *PaymentRepository) FindByCheckouID(checkoutID uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("user").Preload("Package").Where("mpesa_checkout_id = ?", checkoutID).First(&payment).Error
	return &payment, err
}

func (r *PaymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

func (r *PaymentRepository) ListByUser(userID uuid.UUID, limit, offset int) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Preload("package").Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) List(limit, offset int) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Preload("User").Preload("Package").Order("created_at DESC").Limit(limit).Offset(offset).Find(&payments).Error
	return payments, err
}


