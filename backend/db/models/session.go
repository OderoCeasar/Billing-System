package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)



type SessionStatus string

const (
	SessionStatusActive		SessionStatus = "active"
	SessionStatusExpired	SessionStatus = "expired"
	SessionStatusDisconnected	SessionStatus = "disconnected"
)

type Session struct {
	ID			uuid.UUID 		`gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID		uuid.UUID		`gorm:"type:uuid;not null" json:"user_id"`
	PackageID 	uuid.UUID		`gorm:"type:uuid;not null" json:"package_id"`
	PaymentID 	uuid.UUID		`gorm:"type:uuid" json:"payment_id"`
	SessionID 	string			`gorm:"uniqueIndex" json:"session_id"`
	Username	string			`gorm:"not null" json:"username"`
	NASIPAddress	string		`json:"nas_ip_address"`
	NASPortID	string			`json:"nas_port_id"`
	IPAddress 	string			`json:"ip_address"`
	MACAddress 	string			`json:"mac_address"`
	Status 		SessionStatus	`gorm:"type:varchar(20);default:'active'" json:"status"`
	StartTime 	time.Time		`gorm:"not null" json:"start_time"`
	EndTime		*time.Time		`gorm:"end_time"`
	ExpiresAt	time.Time		`gorm:"not null" json:"expires_at"`
	DataUsedBytes  	int64		`gorm:"default:0" json:"data_used_bytes"`
	DataLimitBytes  int64		`json:"data_limit_bytes"`
	TimeUsedMinutes int			`gorm:"default:0" json:"time_used_minutes"`
	TimeLimitMinutes  	int		`json:"time_limit_minutes"`
	LastUpdateTime		time.Time	`json:"last_update_time"`
	CreatedAt 	time.Time		`json:"created_at"`
	UpdatedAt	time.Time		`json:"updated_at"`
	DeletedAt 	time.Time		`gorm:"index" json:"-"`


	User 		User		`gorm:"foreignKey:UserID" json:"user,omitempty"`
	Package  	Package		`gorm:"foreignKey:PackageID" json:"package,omitempty"`
	Payment		Payment		`gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
}


func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID  == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.StartTime.IsZero() {
		s.StartTime = time.Now()
	}
	if s.LastUpdateTime.IsZero() {
		s.LastUpdateTime = time.Now()
	}
	return nil
}


func (s *Session) isExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) HasExceededDataLimit() bool {
	return s.DataLimitBytes > 0 && s.DataUsedBytes >= s.DataLimitBytes
}

func (s *Session) HasExceededTimeLimit() bool {
	return s.TimeLimitMinutes > 0 && s.TimeUsedMinutes >= s.TimeLimitMinutes
}