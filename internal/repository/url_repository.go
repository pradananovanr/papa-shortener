package repository

import (
	"github.com/prdnnvnrnt/papa-shortener/internal/model"

	"gorm.io/gorm"
)

type AdminRepository interface {
	Create(admin *model.Admin) error
	FindByUsername(username string) (*model.Admin, error)
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) Create(admin *model.Admin) error {
	return r.db.Create(admin).Error
}

func (r *adminRepository) FindByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// URLRepository handles URL CRUD operations
type URLRepository interface {
	Create(url *model.URL) error
	FindByShortCode(code string) (*model.URL, error)
	FindAll(page, perPage int) ([]model.URL, int64, error)
	Update(url *model.URL) error
	Delete(code string) error
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

func (r *urlRepository) FindAll(page, perPage int) ([]model.URL, int64, error) {
	var urls []model.URL
	var total int64

	offset := (page - 1) * perPage

	err := r.db.Model(&model.URL{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&urls).Error
	if err != nil {
		return nil, 0, err
	}

	return urls, total, nil
}

func (r *urlRepository) Update(url *model.URL) error {
	return r.db.Save(url).Error
}

func (r *urlRepository) Delete(code string) error {
	return r.db.Where("short_code = ?", code).Delete(&model.URL{}).Error
}

func (r *urlRepository) ExistsByShortCode(code string) (bool, error) {
	var count int64
	err := r.db.Model(&model.URL{}).Where("short_code = ?", code).Count(&count).Error
	return count > 0, err
}