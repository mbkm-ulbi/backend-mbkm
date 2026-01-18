package models

import (
	"time"

	"gorm.io/gorm"
)

type Report struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	ApplyJobID        uint           `json:"apply_job_id"`
	ApplyJob          *ApplyJob      `gorm:"foreignKey:ApplyJobID" json:"apply_job,omitempty"`
	ReportJobUser     string         `gorm:"size:255" json:"report_job_user"`
	StartDate         *time.Time     `json:"start_date"`
	EndDate           *time.Time     `json:"end_date"`
	Status            string         `gorm:"size:50;default:'Draft'" json:"status"`
	FileLaporan       string         `gorm:"type:text" json:"file_laporan"` // URL or path
	CompanyCheckedID  *uint          `json:"company_checked_id"`
	CompanyCheckedAt  *time.Time     `json:"company_checked_at"`
	LecturerCheckedID *uint          `json:"lecturer_checked_id"`
	LecturerCheckedAt *time.Time     `json:"lecturer_checked_at"`
	ExaminerCheckedID *uint          `json:"examiner_checked_id"`
	ExaminerCheckedAt *time.Time     `json:"examiner_checked_at"`
	ProdiCheckedID    *uint          `json:"prodi_checked_id"`
	ProdiCheckedAt    *time.Time     `json:"prodi_checked_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	ActivityDetails []ActivityDetail `gorm:"foreignKey:ReportJobID" json:"activity_details,omitempty"`
	CompanyChecked  *User            `gorm:"foreignKey:CompanyCheckedID" json:"company_checked,omitempty"`
	LecturerChecked *User            `gorm:"foreignKey:LecturerCheckedID" json:"lecturer_checked,omitempty"`
	ExaminerChecked *User            `gorm:"foreignKey:ExaminerCheckedID" json:"examiner_checked,omitempty"`
	ProdiChecked    *User            `gorm:"foreignKey:ProdiCheckedID" json:"prodi_checked,omitempty"`
}

func (Report) TableName() string {
	return "reports"
}

type ActivityDetail struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	ReportJobID     uint           `json:"report_job_id"`
	Date            *time.Time     `json:"date"`
	ActivityDetails string         `gorm:"type:text" json:"activity_details"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (ActivityDetail) TableName() string {
	return "activity_details"
}

type BobotNilai struct {
	ID                   uint           `gorm:"primarykey" json:"id"`
	IDProgramStudi       uint           `json:"id_program_studi"`
	BobotNilaiPerusahaan float64        `json:"bobot_nilai_perusahaan"`
	BobotNilaiPembimbing float64        `json:"bobot_nilai_pembimbing"`
	BobotNilaiPenguji    float64        `json:"bobot_nilai_penguji"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (BobotNilai) TableName() string {
	return "bobot_nilai"
}

type Evaluation struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	ApplyJobID uint      `json:"apply_job_id"`
	ApplyJob   *ApplyJob `gorm:"foreignKey:ApplyJobID" json:"apply_job,omitempty"`
	Status     string    `gorm:"size:50" json:"status"`
	Grade      string    `gorm:"size:5" json:"grade"`

	CompanyPersonnelID      *uint      `json:"company_personnel_id"`
	CompanyGrade            string     `gorm:"size:5" json:"company_grade"`
	CompanyGradeScore       float64    `json:"company_grade_score"`
	CompanyGradeDescription string     `gorm:"type:text" json:"company_grade_description"`
	CompanyGradeDate        *time.Time `json:"company_grade_date"`

	LecturerID               *uint      `json:"lecturer_id"` // Note: Laravel uses 'lecture_id' sometimes, verify column
	LecturerGrade            string     `gorm:"size:5" json:"lecturer_grade"`
	LecturerGradeScore       float64    `json:"lecturer_grade_score"`
	LecturerGradeDescription string     `gorm:"type:text" json:"lecturer_grade_description"`
	LecturerGradeDate        *time.Time `json:"lecturer_grade_date"`

	ExaminerID               *uint      `json:"examiner_id"`
	ExaminerGrade            string     `gorm:"size:5" json:"examiner_grade"`
	ExaminerGradeScore       float64    `json:"examiner_grade_score"`
	ExaminerGradeDescription string     `gorm:"type:text" json:"examiner_grade_description"`
	ExaminerGradeDate        *time.Time `json:"examiner_grade_date"`

	ProdiID               *uint      `json:"prodi_id"`
	ProdiGrade            string     `gorm:"size:5" json:"prodi_grade"`
	ProdiGradeScore       float64    `json:"prodi_grade_score"`
	ProdiGradeDescription string     `gorm:"type:text" json:"prodi_grade_description"`
	ProdiGradeDate        *time.Time `json:"prodi_grade_date"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	CompanyPersonnel *User `gorm:"foreignKey:CompanyPersonnelID" json:"company_personnel,omitempty"`
	Lecturer         *User `gorm:"foreignKey:LecturerID" json:"lecturer,omitempty"`
	Examiner         *User `gorm:"foreignKey:ExaminerID" json:"examiner,omitempty"`
	Prodi            *User `gorm:"foreignKey:ProdiID" json:"prodi,omitempty"`
}

func (Evaluation) TableName() string {
	return "evaluations"
}

type KonversiNilai struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	ApplyJobID uint           `json:"apply_job_id"`
	MatkulID   uint           `json:"matkul_id"` // or mata_kuliah_id
	Grade      string         `gorm:"size:5" json:"grade"`
	Score      float64        `json:"score"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	MataKuliah *MataKuliah `gorm:"foreignKey:MatkulID" json:"mata_kuliah,omitempty"`
	ApplyJob   *ApplyJob   `gorm:"foreignKey:ApplyJobID" json:"apply_job,omitempty"`
}

func (KonversiNilai) TableName() string {
	return "konversi_nilai"
}
