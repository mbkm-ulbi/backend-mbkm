package models

import (
	"time"

	"gorm.io/gorm"
)

type ApplyJob struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	JobUser               *string        `gorm:"size:255" json:"job_user,omitempty"`
	Status                *string        `gorm:"size:100" json:"status,omitempty"`
	ResponsibleLecturerID *uint          `json:"responsible_lecturer_id,omitempty"`
	ExaminerLecturerID    *uint          `json:"examiner_lecturer_id,omitempty"`
	CreatedByID           *uint          `json:"created_by_id,omitempty"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Users               []User          `gorm:"many2many:apply_job_user" json:"users,omitempty"`
	Jobs                []Job           `gorm:"many2many:apply_job_job" json:"jobs,omitempty"`
	CreatedBy           *User           `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	ResponsibleLecturer *User           `gorm:"foreignKey:ResponsibleLecturerID" json:"responsible_lecturer,omitempty"`
	ExaminerLecturer    *User           `gorm:"foreignKey:ExaminerLecturerID" json:"examiner_lecturer,omitempty"`
	Report              *Report         `gorm:"foreignKey:ApplyJobID" json:"report,omitempty"`
	Evaluation          *Evaluation     `gorm:"foreignKey:ApplyJobID" json:"evaluation,omitempty"`
	KonversiNilai       []KonversiNilai `gorm:"foreignKey:ApplyJobID" json:"konversi_nilai,omitempty"`

	// Virtual fields for media
	DHS                   *string `gorm:"-" json:"dhs,omitempty"`
	KTM                   *string `gorm:"-" json:"ktm,omitempty"`
	CV                    *string `gorm:"-" json:"cv,omitempty"`
	SuratLamaran          *string `gorm:"-" json:"surat_lamaran,omitempty"`
	SuratRekomendasiProdi *string `gorm:"-" json:"surat_rekomendasi_prodi,omitempty"`
}

func (ApplyJob) TableName() string {
	return "apply_jobs"
}

type ApplyJobMonthlyLog struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ApplyJobID uint           `json:"apply_job_id"`
	Month      int            `json:"month"`
	Year       int            `json:"year"`
	Status     *string        `gorm:"size:100" json:"status,omitempty"`
	Content    *string        `gorm:"type:text" json:"content,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	ApplyJob *ApplyJob `gorm:"foreignKey:ApplyJobID" json:"apply_job,omitempty"`
}

func (ApplyJobMonthlyLog) TableName() string {
	return "apply_job_monthly_logs"
}

type KonversiNilai struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ApplyJobID uint           `json:"apply_job_id"`
	MataKuliah *string        `gorm:"size:255" json:"mata_kuliah,omitempty"`
	SKS        *int           `json:"sks,omitempty"`
	Nilai      *string        `gorm:"size:10" json:"nilai,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (KonversiNilai) TableName() string {
	return "konversi_nilais"
}
