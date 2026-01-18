package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Name               string         `gorm:"size:255;not null" json:"name"`
	Email              string         `gorm:"size:255;uniqueIndex" json:"email"`
	Username           string         `gorm:"size:255;uniqueIndex" json:"username"`
	Password           string         `gorm:"size:255" json:"-"`
	NIM                *string        `gorm:"size:50" json:"nim,omitempty"`
	ProgramStudy       *string        `gorm:"size:255" json:"program_study,omitempty"`
	Faculty            *string        `gorm:"size:255" json:"faculty,omitempty"`
	Semester           *string        `gorm:"column:semester;size:255" json:"semester,omitempty"`
	PhoneNumber        *string        `gorm:"size:50" json:"phone_number,omitempty"`
	Address            *string        `gorm:"type:text" json:"address,omitempty"`
	SocialMedia        *string        `gorm:"size:255" json:"social_media,omitempty"`
	EmergencyContact   *string        `gorm:"size:255" json:"emergency_contact,omitempty"`
	ProfileDescription *string        `gorm:"type:text" json:"profile_description,omitempty"`
	Position           *string        `gorm:"size:255" json:"position,omitempty"`
	IPK                *string        `gorm:"column:ipk;size:255" json:"ipk,omitempty"`
	Birthdate          *time.Time     `json:"birthdate,omitempty"`
	Role               string         `gorm:"size:50;default:'student'" json:"role"`
	Status             *string        `gorm:"size:50" json:"status,omitempty"`
	Verified           bool           `gorm:"default:false" json:"verified"`
	VerifiedAt         *time.Time     `json:"verified_at,omitempty"`
	VerificationToken  *string        `gorm:"size:255" json:"verification_token,omitempty"`
	Approved           bool           `gorm:"default:true" json:"approved"`
	TwoFactor          bool           `gorm:"default:false" json:"two_factor"`
	TwoFactorCode      *string        `gorm:"size:10" json:"-"`
	TwoFactorExpiresAt *time.Time     `json:"-"`
	TeamID             *uint          `json:"team_id,omitempty"`
	IdProgramStudi     *string        `gorm:"column:id_program_studi;size:40" json:"id_program_studi,omitempty"`
	EmailVerifiedAt    *time.Time     `json:"email_verified_at,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Roles []Role `gorm:"many2many:role_user" json:"roles,omitempty"`
	Team  *Team  `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}

func (User) TableName() string {
	return "users"
}
