package models

import (
	"time"

	"gorm.io/gorm"
)

// Evaluation status constants
const (
	EvalStatusInProgress       = "Sedang Dibuat"
	EvalStatusCreatedByStudent = "Selesai Dibuat Oleh Mahasiswa"
	EvalStatusApprovedCompany  = "Disetujui Oleh Perusahaan"
	EvalStatusApprovedLecturer = "Disetujui Oleh Dosen Wali"
	EvalStatusApprovedProdi    = "Disetujui Oleh Prodi"
)

type Evaluation struct {
	ID                       uint           `gorm:"primaryKey" json:"id"`
	ApplyJobID               *uint          `json:"apply_job_id,omitempty"`
	Status                   *string        `gorm:"size:100" json:"status,omitempty"`
	Grade                    *string        `gorm:"size:10" json:"grade,omitempty"`
	CompanyPersonnelID       *uint          `json:"company_personnel_id,omitempty"`
	CompanyGrade             *string        `gorm:"size:10" json:"company_grade,omitempty"`
	CompanyGradeScore        *float64       `json:"company_grade_score,omitempty"`
	CompanyGradeDescription  *string        `gorm:"type:text" json:"company_grade_description,omitempty"`
	CompanyGradeDate         *time.Time     `json:"company_grade_date,omitempty"`
	LecturerID               *uint          `json:"lecturer_id,omitempty"`
	LecturerGrade            *string        `gorm:"size:10" json:"lecturer_grade,omitempty"`
	LecturerGradeScore       *float64       `json:"lecturer_grade_score,omitempty"`
	LecturerGradeDescription *string        `gorm:"type:text" json:"lecturer_grade_description,omitempty"`
	LecturerGradeDate        *time.Time     `json:"lecturer_grade_date,omitempty"`
	ProdiID                  *uint          `json:"prodi_id,omitempty"`
	ProdiGrade               *string        `gorm:"size:10" json:"prodi_grade,omitempty"`
	ProdiGradeScore          *float64       `json:"prodi_grade_score,omitempty"`
	ProdiGradeDescription    *string        `gorm:"type:text" json:"prodi_grade_description,omitempty"`
	ProdiGradeDate           *time.Time     `json:"prodi_grade_date,omitempty"`
	ExaminerID               *uint          `json:"examiner_id,omitempty"`
	ExaminerGrade            *string        `gorm:"size:10" json:"examiner_grade,omitempty"`
	ExaminerGradeScore       *float64       `json:"examiner_grade_score,omitempty"`
	ExaminerGradeDescription *string        `gorm:"type:text" json:"examiner_grade_description,omitempty"`
	ExaminerGradeDate        *time.Time     `json:"examiner_grade_date,omitempty"`
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	ApplyJob         *ApplyJob `gorm:"foreignKey:ApplyJobID" json:"apply_job,omitempty"`
	CompanyPersonnel *User     `gorm:"foreignKey:CompanyPersonnelID" json:"company_personnel,omitempty"`
	Lecturer         *User     `gorm:"foreignKey:LecturerID" json:"lecturer,omitempty"`
	Prodi            *User     `gorm:"foreignKey:ProdiID" json:"prodi,omitempty"`
	Examiner         *User     `gorm:"foreignKey:ExaminerID" json:"examiner,omitempty"`
}

func (Evaluation) TableName() string {
	return "evaluations"
}

type BobotNilai struct {
	ID                   uint    `gorm:"primaryKey" json:"id"`
	IdProgramStudi       *uint   `json:"id_program_studi,omitempty"`
	BobotNilaiPerusahaan float64 `json:"bobot_nilai_perusahaan"`
	BobotNilaiPembimbing float64 `json:"bobot_nilai_pembimbing"`
	BobotNilaiPenguji    float64 `json:"bobot_nilai_penguji"`
}

func (BobotNilai) TableName() string {
	return "bobot_nilais"
}

// CalculateFinalGrade calculates the final grade based on weighted scores
func (e *Evaluation) CalculateFinalGrade(bobot *BobotNilai) map[string]interface{} {
	if bobot == nil {
		return nil
	}

	totalBobot := bobot.BobotNilaiPerusahaan + bobot.BobotNilaiPembimbing + bobot.BobotNilaiPenguji
	if totalBobot == 0 {
		return nil
	}

	var companyScore, lecturerScore, examinerScore float64

	if e.CompanyGradeScore != nil {
		companyScore = *e.CompanyGradeScore * (bobot.BobotNilaiPerusahaan / totalBobot)
	}
	if e.LecturerGradeScore != nil {
		lecturerScore = *e.LecturerGradeScore * (bobot.BobotNilaiPembimbing / totalBobot)
	}
	if e.ExaminerGradeScore != nil {
		examinerScore = *e.ExaminerGradeScore * (bobot.BobotNilaiPenguji / totalBobot)
	}

	totalScore := companyScore + lecturerScore + examinerScore
	hasGrade := companyScore > 0 && lecturerScore > 0 && examinerScore > 0

	grade := calculateGradeLetter(totalScore, hasGrade)

	return map[string]interface{}{
		"company":     companyScore,
		"lecturer":    lecturerScore,
		"examiner":    examinerScore,
		"total_score": totalScore,
		"grade":       grade,
	}
}

func calculateGradeLetter(score float64, hasGrade bool) string {
	if !hasGrade {
		return "-"
	}
	if score >= 85 {
		return "A"
	} else if score >= 70 {
		return "B"
	} else if score >= 55 {
		return "C"
	} else if score >= 40 {
		return "D"
	}
	return "E"
}
