package models

import (
	"time"

	"gorm.io/gorm"
)

type Fakultas struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Nama      string         `gorm:"size:255" json:"nama"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Fakultas) TableName() string {
	return "fakultas"
}

type ProgramStudi struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Nama         string         `gorm:"size:255" json:"nama"`
	IDUnitParent uint           `json:"id_unit_parent"`
	Fakultas     *Fakultas      `gorm:"foreignKey:IDUnitParent" json:"fakultas,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (ProgramStudi) TableName() string {
	return "program_studi"
}

type MataKuliah struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	KodeMatkul     string         `gorm:"size:50" json:"kode_matkul"`
	NamaMatkul     string         `gorm:"size:255" json:"nama_matkul"`
	Sks            int            `json:"sks"`
	IDProgramStudi uint           `json:"id_program_studi"`
	ProgramStudi   *ProgramStudi  `gorm:"foreignKey:IDProgramStudi" json:"program_studi,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (MataKuliah) TableName() string {
	return "mata_kuliah"
}

type Perusahaan struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	NamaPerusahaan    string         `gorm:"size:255" json:"nama_perusahaan"`
	AlamatPerusahaan  string         `gorm:"type:text" json:"alamat_perusahaan"`
	EmailPerusahaan   string         `gorm:"size:255" json:"email_perusahaan"`
	WebsitePerusahaan string         `gorm:"size:255" json:"website_perusahaan"`
	Facebook          string         `gorm:"size:255" json:"facebook"`
	Instagram         string         `gorm:"size:255" json:"instagram"`
	Tiktok            string         `gorm:"size:255" json:"tiktok"`
	Linkedin          string         `gorm:"size:255" json:"linkedin"`
	UserID            uint           `json:"user_id"`
	User              *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Perusahaan) TableName() string {
	return "perusahaans"
}
