package service

import (
	"errors"
	"strings"

	"github.com/prdnnvnrnt/papa-shortener/internal/config"
	"github.com/prdnnvnrnt/papa-shortener/internal/middleware"
	"github.com/prdnnvnrnt/papa-shortener/internal/model"
	"github.com/prdnnvnrnt/papa-shortener/internal/repository"
	"github.com/prdnnvnrnt/papa-shortener/pkg/validator"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUsernameExists     = errors.New("username already exists")
)

type AuthService interface {
	Login(username, password string) (string, error)
	RegisterAdmin(username, password string) error
}

type authService struct {
	repo   repository.AdminRepository
	config *config.Config
}

func NewAuthService(repo repository.AdminRepository, cfg *config.Config) AuthService {
	return &authService{
		repo:   repo,
		config: cfg,
	}
}

func (s *authService) Login(username, password string) (string, error) {
	admin, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !admin.CheckPassword(password) {
		return "", ErrInvalidCredentials
	}

	return middleware.GenerateToken(admin.ID, admin.Username)
}

func (s *authService) RegisterAdmin(username, password string) error {
	// Check if username already exists
	existing, _ := s.repo.FindByUsername(username)
	if existing != nil {
		return ErrUsernameExists
	}

	admin := &model.Admin{
		Username: username,
	}

	if err := admin.SetPassword(password); err != nil {
		return err
	}

	return s.repo.Create(admin)
}

type AdminService interface {
	ListURLs(page, perPage int) (*model.AdminListResponse, error)
	GetURL(code string) (*model.URL, error)
	UpdateURL(code string, req *model.UpdateURLRequest) error
	DeleteURL(code string) error
	ToggleURLActive(code string) error
}

type adminService struct {
	repo   repository.URLRepository
	config *config.Config
}

func NewAdminService(repo repository.URLRepository, cfg *config.Config) AdminService {
	return &adminService{
		repo:   repo,
		config: cfg,
	}
}

func (s *adminService) ListURLs(page, perPage int) (*model.AdminListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	urls, total, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, err
	}

	return &model.AdminListResponse{
		URLs:    urls,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}, nil
}

func (s *adminService) GetURL(code string) (*model.URL, error) {
	url, err := s.repo.FindByShortCode(code)
	if err != nil {
		return nil, ErrURLNotFound
	}
	return url, nil
}

func (s *adminService) UpdateURL(code string, req *model.UpdateURLRequest) error {
	url, err := s.repo.FindByShortCode(code)
	if err != nil {
		return ErrURLNotFound
	}

	// Update fields if provided
	if req.OriginalURL != "" {
		original := req.OriginalURL
		if !strings.HasPrefix(original, "http://") && !strings.HasPrefix(original, "https://") {
			original = "https://" + original
		}
		url.Original = original
	}

	// Handle custom URL change
	if req.CustomURL != "" {
		newCode := strings.ToLower(req.CustomURL)
		exists, err := s.repo.ExistsByShortCode(newCode)
		if err != nil {
			return err
		}
		if exists && newCode != code {
			return ErrCustomURLExists
		}
		if !validator.IsValidCustomURL(newCode) {
			return ErrInvalidCustomURL
		}
		url.ShortCode = newCode
		url.Custom = true
	}

	if req.Active != nil {
		url.Active = *req.Active
	}

	return s.repo.Update(url)
}

func (s *adminService) DeleteURL(code string) error {
	_, err := s.repo.FindByShortCode(code)
	if err != nil {
		return ErrURLNotFound
	}
	return s.repo.Delete(code)
}

func (s *adminService) ToggleURLActive(code string) error {
	url, err := s.repo.FindByShortCode(code)
	if err != nil {
		return ErrURLNotFound
	}
	url.Active = !url.Active
	return s.repo.Update(url)
}