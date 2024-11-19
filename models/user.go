package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Name             string         `json:"name"`
	Email            string         `json:"email" gorm:"unique"`
	Username         string         `json:"username" gorm:"unique"`
	Password         string         `json:"-"`    // Disimpan dalam hash
	Role             string         `json:"role"` // Misalnya "admin" atau "user"
	NIM              string         `json:"nim"`  // Nomor Induk Mahasiswa
	ProgramStudy     string         `json:"program_study"`
	Faculty          string         `json:"faculty"`
	Semester         string         `json:"semester"`
	ProfilePicture   string         `json:"profile_picture"` // Path foto profil
	SocialMedia      string         `json:"social_media"`
	EmergencyContact string         `json:"emergency_contact"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
