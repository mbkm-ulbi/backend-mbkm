package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Picture   string         `json:"picture"`   // URL gambar artikel
	Author    string         `json:"author"`    // Penulis artikel
	Title     string         `json:"title"`     // Judul artikel
	Content   string         `json:"content"`   // Isi artikel
	PostDate  time.Time      `json:"post_date"` // Tanggal posting
	Views     int            `json:"views"`     // Jumlah views
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
