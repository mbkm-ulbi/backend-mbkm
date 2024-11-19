package landing

import (
	"github.com/mongmx/fiber-cms/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetLandingPage() (*models.Landing, error) {
	var landing models.Landing
	err := r.db.First(&landing).Error
	if err != nil {
		return nil, err
	}
	return &landing, nil
}

func (r *Repository) UpdateLandingPage(landing *models.Landing) error {
	return r.db.Save(landing).Error
}
