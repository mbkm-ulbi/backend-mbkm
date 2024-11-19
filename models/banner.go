package models

type Banner struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Info     string `json:"info"`
	Link     string `json:"link"`
	ImageURL string `json:"image_url"`
}
