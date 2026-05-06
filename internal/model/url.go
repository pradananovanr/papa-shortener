package model

import (
	"time"

	"gorm.io/gorm"
)

type URL struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ShortCode string         `gorm:"uniqueIndex;size:10;not null" json:"short_code"`
	Original  string         `gorm:"size:2048;not null" json:"original_url"`
	Custom    bool           `gorm:"default:false" json:"is_custom"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateURLRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
	CustomURL   string `json:"custom_url" validate:"omitempty,alphanum,min=3,max=20"`
}

type CreateURLResponse struct {
	ShortCode  string `json:"short_code"`
	Original   string `json:"original_url"`
	ShortURL   string `json:"short_url"`
	IsCustom   bool   `json:"is_custom"`
	FullShort  string `json:"full_short_url"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	BaseURL string `json:"base_url,omitempty"`
}