package models

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	ID                        uint           `gorm:"primaryKey" json:"id"`
	CompanyName               string         `gorm:"size:255" json:"company_name"`
	BusinessFields            *string        `gorm:"size:255" json:"business_fields,omitempty"`
	CompanySize               *string        `gorm:"size:100" json:"company_size,omitempty"`
	CompanyWebsite            *string        `gorm:"size:255" json:"company_website,omitempty"`
	CompanyProfileDescription *string        `gorm:"type:text" json:"company_profile_description,omitempty"`
	CompanyPhoneNumber        *string        `gorm:"size:50" json:"company_phone_number,omitempty"`
	CompanyAddress            *string        `gorm:"type:text" json:"company_address,omitempty"`
	UserID                    *uint          `json:"user_id,omitempty"`
	CreatedByID               *uint          `json:"created_by_id,omitempty"`
	CreatedAt                 time.Time      `json:"created_at"`
	UpdatedAt                 time.Time      `json:"updated_at"`
	DeletedAt                 gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User      *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedBy *User `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`

	// Virtual field for media
	CompanyLogo *string `gorm:"-" json:"company_logo,omitempty"`
}

func (Company) TableName() string {
	return "companies"
}
