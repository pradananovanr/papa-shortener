package service

import (
	"errors"
	"strings"

	"github.com/prdnnvnrnt/papa-shortener/internal/config"
	"github.com/prdnnvnrnt/papa-shortener/internal/model"
	"github.com/prdnnvnrnt/papa-shortener/internal/repository"
	"github.com/prdnnvnrnt/papa-shortener/pkg/validator"

	"github.com/google/uuid"
)

var (
	ErrInvalidCustomURL   = errors.New("invalid custom URL: must be 3-20 lowercase alphanumeric characters")
	ErrCustomURLExists    = errors.New("custom URL already taken")
	ErrURLNotFound        = errors.New("URL not found")
	ErrInvalidOriginalURL = errors.New("invalid original URL format")
)

type URLService interface {
	CreateShortURL(req *model.CreateURLRequest) (*model.CreateURLResponse, error)
	ResolveShortCode(code string) (string, error)
}

type urlService struct {
	repo   repository.URLRepository
	config *config.Config
}

func NewURLService(repo repository.URLRepository, cfg *config.Config) URLService {
	return &urlService{
		repo:   repo,
		config: cfg,
	}
}

func (s *urlService) generateShortCode() string {
	return strings.ToLower(uuid.New().String()[:8])
}

func (s *urlService) CreateShortURL(req *model.CreateURLRequest) (*model.CreateURLResponse, error) {
	if req.OriginalURL == "" {
		return nil, ErrInvalidOriginalURL
	}

	original := req.OriginalURL
	if !strings.HasPrefix(original, "http://") && !strings.HasPrefix(original, "https://") {
		original = "https://" + original
	}

	isCustom := false
	shortCode := s.generateShortCode()

	if req.CustomURL != "" {
		custom := strings.ToLower(req.CustomURL)

		if !validator.IsValidCustomURL(custom) {
			return nil, ErrInvalidCustomURL
		}

		exists, err := s.repo.ExistsByShortCode(custom)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrCustomURLExists
		}

		shortCode = custom
		isCustom = true
	}

	url := &model.URL{
		ShortCode: shortCode,
		Original:  original,
		Custom:    isCustom,
		Active:    true,
	}

	if err := s.repo.Create(url); err != nil {
		return nil, err
	}

	fullShort := s.config.App.BaseURL + "/" + shortCode

	return &model.CreateURLResponse{
		ShortCode: shortCode,
		Original:  original,
		ShortURL:  shortCode,
		IsCustom:  isCustom,
		FullShort: fullShort,
	}, nil
}

func (s *urlService) ResolveShortCode(code string) (string, error) {
	url, err := s.repo.FindByShortCode(code)
	if err != nil {
		return "", ErrURLNotFound
	}

	// Check if URL is active
	if !url.Active {
		return "", ErrURLNotFound
	}

	return url.Original, nil
}