package models

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Title        string         `json:"title"`
	Company      string         `json:"company"`
	Location     string         `json:"location"`
	Description  string         `json:"description"`
	Requirements string         `json:"requirements"`
	Benefits     string         `json:"benefits"`
	JobType      string         `json:"job_type"`
	Salary       string         `json:"salary"`
	ImageURL     string         `json:"image_url"` // Menyimpan URL gambar
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
