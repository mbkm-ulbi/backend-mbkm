package models

import (
	"time"

	"gorm.io/gorm"
)

// Report status constants
const (
	ReportStatusCreatedByStudent  = "Selesai Dibuat Oleh Mahasiswa"
	ReportStatusCheckedByCompany  = "Selesai Diperiksa Oleh Perusahaan"
	ReportStatusCheckedByLecturer = "Selesai Diperiksa Oleh Dosen Wali"
	ReportStatusCheckedByProdi    = "Selesai Diperiksa Oleh Prodi"
)

type Report struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	StartDate         *time.Time     `json:"start_date,omitempty"`
	EndDate           *time.Time     `json:"end_date,omitempty"`
	Status            *string        `gorm:"size:100" json:"status,omitempty"`
	ApplyJobID        *uint          `json:"apply_job_id,omitempty"`
	ReportJobUser     *string        `gorm:"size:255" json:"report_job_user,omitempty"`
	CompanyCheckedID  *uint          `json:"company_checked_id,omitempty"`
	LecturerCheckedID *uint          `json:"lecturer_checked_id,omitempty"`
	ProdiCheckedID    *uint          `json:"prodi_checked_id,omitempty"`
	ExaminerCheckedID *uint          `json:"examiner_checked_id,omitempty"`
	CompanyCheckedAt  *time.Time     `json:"company_checked_at,omitempty"`
	LecturerCheckedAt *time.Time     `json:"lecturer_checked_at,omitempty"`
	ProdiCheckedAt    *time.Time     `json:"prodi_checked_at,omitempty"`
	ExaminerCheckedAt *time.Time     `json:"examiner_checked_at,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	ApplyJob        *ApplyJob        `gorm:"foreignKey:ApplyJobID" json:"apply_job,omitempty"`
	CompanyChecked  *User            `gorm:"foreignKey:CompanyCheckedID" json:"company_checked,omitempty"`
	LecturerChecked *User            `gorm:"foreignKey:LecturerCheckedID" json:"lecturer_checked,omitempty"`
	ProdiChecked    *User            `gorm:"foreignKey:ProdiCheckedID" json:"prodi_checked,omitempty"`
	ExaminerChecked *User            `gorm:"foreignKey:ExaminerCheckedID" json:"examiner_checked,omitempty"`
	ActivityDetails []ActivityDetail `gorm:"foreignKey:ReportJobID" json:"activity_details,omitempty"`

	// Virtual field for media
	FileLaporan *string `gorm:"-" json:"file_laporan,omitempty"`
}

func (Report) TableName() string {
	return "reports"
}

type ActivityDetail struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ReportJobID uint           `json:"report_job_id"`
	Date        *time.Time     `json:"date,omitempty"`
	Activity    *string        `gorm:"type:text" json:"activity,omitempty"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ActivityDetail) TableName() string {
	return "activity_details"
}
