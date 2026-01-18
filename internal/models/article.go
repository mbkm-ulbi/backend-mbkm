package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	Content   *string        `gorm:"type:text" json:"content,omitempty"`
	Views     int            `gorm:"default:0" json:"views"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Virtual field for media
	Picture *string `gorm:"-" json:"picture,omitempty"`
}

func (Article) TableName() string {
	return "articles"
}
