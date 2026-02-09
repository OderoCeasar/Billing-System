package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)



type PaymentStatus	string 

const (
	PaymentStatusPending		PaymentStatus = "pending"
	PaymentStatusComplete		PaymentStatus = "Completed"
	PaymentStatusFailed 		PaymentStatus = "failed"
	PaymentStatusCancelled		PaymentStatus = "cancelled"
)

type Payment struct {
	ID 			uuid.UUID 			`gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID			`gorm:"type:uuid;not null" json:"user_id"`
	PackageID	uuid.UUID			`gorm:"type:uuid;not null" json:"package_id"`
	Amount		float64				`gorm:"not null" json:"amount"`
	PhoneNumber string				`gorm:"not null" json:"phone_number"`
	Status 		PaymentStatus		`gorm:"type:varchar(20);default:'pending'" json:"status"`
	MpesaCheckoutID string			`gorm:"uniqueIndex" json:"mpesa_checkout_id"`
	MpesaReceiptNumber	string		`json:"mpesa_receipt_number"`
	MpesaTransactionID 	string		`json:"mpesa_transaction_id"`
	TransactionDate		*time.Time	`json:"transaction_date"`
	ResultCode 		int 			`json:"result_code"`
	ResultDesc    	string			`json:"result_desc"`
	CallbackReceived	bool 		`gorm:"default:false" json:"callback_received"`
	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt   	time.Time		`json:"updated_at"`
	DeletedAt  		gorm.DeletedAt	`gotm:"index" json:"-"`


	User	User			`gorm:"foreignKey:UserID" json:"user, omitempty"`
	Package	Package			`gorm:"foreignKey:PackageID" json:"package, omitempty"`

}


func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}