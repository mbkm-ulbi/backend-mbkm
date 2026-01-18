package models

import (
	"time"

	"gorm.io/gorm"
)

// Job status constants
const (
	JobStatusPending      = "Perlu Ditinjau"
	JobStatusAvailable    = "Tersedia"
	JobStatusRejected     = "Ditolak"
	JobStatusClosed       = "Ditutup"
	JobStatusNotAvailable = "Tidak Tersedia"
)

// Job type constants
const (
	JobTypeFullTime  = "Full-time"
	JobTypePartTime  = "Part-time"
	JobTypeFreelance = "Freelance"
	JobTypeHybrid    = "Hybrid"
)

// Vacancy type constants
const (
	VacancyTypeUmum = "Umum"
	VacancyTypeS1   = "S1"
	VacancyTypeS2   = "S2"
	VacancyTypeS3   = "S3"
)

type Job struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Company     string         `gorm:"size:255" json:"company"`
	Location    string         `gorm:"size:255" json:"location"`
	Duration    *string        `gorm:"size:255" json:"duration,omitempty"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	Benefits    *string        `gorm:"type:text" json:"benefits,omitempty"`
	JobType     *string        `gorm:"size:50" json:"job_type,omitempty"`
	Salary      *string        `gorm:"size:255" json:"salary,omitempty"`
	VacancyType *string        `gorm:"size:50" json:"vacancy_type,omitempty"`
	Status      string         `gorm:"size:50;default:'Perlu Ditinjau'" json:"status"`
	MataKuliah  *string        `gorm:"type:text" json:"mata_kuliah,omitempty"`
	Deadline    *time.Time     `json:"deadline,omitempty"`
	CompanyID   *uint          `json:"company_id,omitempty"`
	CreatedByID uint           `json:"created_by_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	CreatedBy *User `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`

	// Virtual field for media (will be handled separately)
	JobVacancyImage *string `gorm:"-" json:"job_vacancy_image,omitempty"`
}

func (Job) TableName() string {
	return "jobs"
}
