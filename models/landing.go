package models

type Landing struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Header    string `json:"header" gorm:"type:text"`  // JSON string untuk header
	Banner1   string `json:"banner1" gorm:"type:text"` // Path URL banner 1
	Banner2   string `json:"banner2" gorm:"type:text"` // Path URL banner 2
	Filters   string `json:"filters" gorm:"type:text"` // JSON string untuk kategori dan program
	Footer    string `json:"footer" gorm:"type:text"`  // JSON string untuk footer
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
