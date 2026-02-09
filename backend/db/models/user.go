package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)



type User struct {
	ID		uuid.UUID		`gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PhoneNumber		string	`gorm:"uniqueIndex;not null" json:"phone_number"`
	Username		string	`gorm:"uniqueIndex" json:"username"`
	Password		string	`gorm:"not null" json:"-"`
	IsActive		bool	`gorm:"default:true" json:"is_active"`
	IsAdmin			bool 	`gorm:"default:false" json:"is_admin"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt 		time.Time	`json:"updated_at"`
	DeletedAt		gorm.DeletedAt	`gorm:"index" json:"-"`

 
	Sessions []Session		`gorm:"foreignKey:UserID" json:"sessions,omitempty"`
	Payments []Payment 		`gorm:"foreignKey:UserID" json:"payments,omitempty"` 
}


func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

