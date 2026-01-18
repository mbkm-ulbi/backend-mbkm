package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Users       []User       `gorm:"many2many:role_user" json:"users,omitempty"`
	Permissions []Permission `gorm:"many2many:permission_role" json:"permissions,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}

type Permission struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Roles []Role `gorm:"many2many:permission_role" json:"roles,omitempty"`
}

func (Permission) TableName() string {
	return "permissions"
}

type Team struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255" json:"name"`
	OwnerID   uint           `json:"owner_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Owner User `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

func (Team) TableName() string {
	return "teams"
}
