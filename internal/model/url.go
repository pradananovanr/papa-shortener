package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Admin struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *Admin) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(hash)
	return nil
}

func (a *Admin) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

type URL struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ShortCode string         `gorm:"uniqueIndex;size:255;not null" json:"short_code"`
	Original  string         `gorm:"size:2048;not null" json:"original_url"`
	Custom    bool           `gorm:"default:false" json:"is_custom"`
	Active    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateURLRequest struct {
	OriginalURL string `json:"original_url" form:"original_url" validate:"required,url"`
	CustomURL   string `json:"custom_url" form:"custom_url" validate:"omitempty,min=3"`
}

type CreateURLResponse struct {
	ShortCode  string `json:"short_code"`
	Original   string `json:"original_url"`
	ShortURL   string `json:"short_url"`
	IsCustom   bool   `json:"is_custom"`
	FullShort  string `json:"full_short_url"`
}

type UpdateURLRequest struct {
	OriginalURL string `json:"original_url,omitempty" form:"original_url,omitempty" validate:"omitempty,url"`
	CustomURL   string `json:"custom_url,omitempty" form:"custom_url,omitempty" validate:"omitempty,min=3"`
	Active      *bool  `json:"active,omitempty" form:"active,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	BaseURL string `json:"base_url,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type AdminListResponse struct {
	URLs      []URL `json:"urls"`
	Total     int64 `json:"total"`
	Page      int   `json:"page"`
	PerPage   int   `json:"per_page"`
}