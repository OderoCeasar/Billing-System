package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)



type PackageType string 

const (
	PackageTypeTime		PackageType = "time"
	PackageTypeData		PackageType = "data"
)

type Package struct {
	ID				uuid.UUID  			`gorm:"type:uuid;primary_key;default:gen_random_uuid" json:"id"`
	Name			string				`gorm:"not null" json:"name"`
	Description 	string				`json:"description"`
	PackageType 	PackageType			`gorm:"type:varchar(20);not null" json:"package_type"`
	Price			float64				`gorm:"not null" json:"price"`
	DurationMinutes	int					`json:"duration_minutes"`
	DataLimitMB 	int64 				`json:"data_limit_md"`
	SpeedLimitUp	int					`json:"speed_limit_up"`
	SpeedLimitDown  int					`json:"speed_limit_down"`
	ValidityDays	int 				`gorm:"default:30" json:"validity_days"`
	IsActive		bool 				`gorm:"default:true" json:"is_active"`
	CreatedAt		time.Time			`json:"created_at"`
	UpdatedAt      	time.Time			`json:"updated_at"`
	DeletedAt		gorm.DeletedAt		`gorm:"index" json:"-"`
}


func (p *Package) BeforeCreate(tx *gorm.DB)	error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()

	}
	return nil
}