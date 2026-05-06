package repository

import (
	"papa-shortener/internal/model"

	"gorm.io/gorm"
)

type URLRepository interface {
	Create(url *model.URL) error
	FindByShortCode(code string) (*model.URL, error)
	ExistsByShortCode(code string) (bool, error)
}

type urlRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) URLRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Create(url *model.URL) error {
	return r.db.Create(url).Error
}

func (r *urlRepository) FindByShortCode(code string) (*model.URL, error) {
	var url model.URL
	err := r.db.Where("short_code = ?", code).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) ExistsByShortCode(code string) (bool, error) {
	var count int64
	err := r.db.Model(&model.URL{}).Where("short_code = ?", code).Count(&count).Error
	return count > 0, err
}