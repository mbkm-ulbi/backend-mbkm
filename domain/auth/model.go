package auth

import (
	"time"

	"github.com/google/uuid"
)

type model struct {
	ID        int64
	UUID      uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type User struct {
	model
	Email            string `json:"email"`
	Username         string `json:"username"`
	Auth             *Auth  `json:"-"`
	Name             string `json:"name"`
	Role             string `json:"role"`
	NIM              string `json:"nim"`
	ProgramStudy     string `json:"program_study"`
	Faculty          string `json:"faculty"`
	Semester         string `json:"semester"`
	ProfilePicture   string `json:"profile_picture"`
	SocialMedia      string `json:"social_media"`
	EmergencyContact string `json:"emergency_contact"`
	Password         string `json:"-"`
}

type Auth struct {
	UserID int64
	Type   string
	Secret string
}

type Topic struct {
	model
	UserID int64  `json:"-"`
	Title  string `json:"title"`
}

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Type     string `json:"type"`
}
