package handler

import (
	"github.com/prdnnvnrnt/papa-shortener/internal/model"
	"github.com/prdnnvnrnt/papa-shortener/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service service.AuthService
	validate *validator.Validate
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{
		service:  svc,
		validate: validator.New(),
	}
}

// Login godoc
// @Summary Admin login
// @Description Login with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "validation_error",
			Message: "Username and password are required",
		})
	}

	token, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{
			Error:   "invalid_credentials",
			Message: "Invalid username or password",
		})
	}

	return c.JSON(model.LoginResponse{Token: token})
}

// RegisterAdmin godoc
// @Summary Register initial admin
// @Description Register the first admin account (only works if no admin exists)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Admin registration"
// @Success 201 {object} model.ErrorResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) RegisterAdmin(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "validation_error",
			Message: "Username and password are required",
		})
	}

	err := h.service.RegisterAdmin(req.Username, req.Password)
	if err != nil {
		if err == service.ErrUsernameExists {
			return c.Status(fiber.StatusConflict).JSON(model.ErrorResponse{
				Error:   "username_exists",
				Message: "Username already taken",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to register admin",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.ErrorResponse{
		Error:   "",
		Message: "Admin registered successfully",
	})
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

type AdminHandler struct {
	service service.AdminService
}

func NewAdminHandler(svc service.AdminService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// ListURLs godoc
// @Summary List all URLs
// @Description Get paginated list of all URLs (admin only)
// @Tags admin
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Security BearerAuth
// @Success 200 {object} model.AdminListResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/admin/urls [get]
func (h *AdminHandler) ListURLs(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	result, err := h.service.ListURLs(page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to fetch URLs",
		})
	}

	return c.JSON(result)
}

// GetURL godoc
// @Summary Get URL details
// @Description Get details of a specific URL by short code
// @Tags admin
// @Produce json
// @Param code path string true "Short code"
// @Security BearerAuth
// @Success 200 {object} model.URL
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/admin/urls/{code} [get]
func (h *AdminHandler) GetURL(c *fiber.Ctx) error {
	code := c.Params("code")

	url, err := h.service.GetURL(code)
	if err != nil {
		if err == service.ErrURLNotFound {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
				Error:   "not_found",
				Message: "URL not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to fetch URL",
		})
	}

	return c.JSON(url)
}

// UpdateURL godoc
// @Summary Update URL
// @Description Update original URL, custom code, or active status
// @Tags admin
// @Accept json
// @Produce json
// @Param code path string true "Short code"
// @Param request body model.UpdateURLRequest true "Update request"
// @Security BearerAuth
// @Success 200 {object} model.ErrorResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/admin/urls/{code} [put]
func (h *AdminHandler) UpdateURL(c *fiber.Ctx) error {
	code := c.Params("code")

	var req model.UpdateURLRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	err := h.service.UpdateURL(code, &req)
	if err != nil {
		if err == service.ErrURLNotFound {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
				Error:   "not_found",
				Message: "URL not found",
			})
		}
		if err == service.ErrCustomURLExists {
			return c.Status(fiber.StatusConflict).JSON(model.ErrorResponse{
				Error:   "custom_url_exists",
				Message: "Custom URL already taken",
			})
		}
		if err == service.ErrInvalidCustomURL {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Error:   "invalid_custom_url",
				Message: "Invalid custom URL format",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update URL",
		})
	}

	return c.JSON(model.ErrorResponse{
		Error:   "",
		Message: "URL updated successfully",
	})
}

// DeleteURL godoc
// @Summary Delete URL
// @Description Permanently delete a URL
// @Tags admin
// @Produce json
// @Param code path string true "Short code"
// @Security BearerAuth
// @Success 200 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/admin/urls/{code} [delete]
func (h *AdminHandler) DeleteURL(c *fiber.Ctx) error {
	code := c.Params("code")

	err := h.service.DeleteURL(code)
	if err != nil {
		if err == service.ErrURLNotFound {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
				Error:   "not_found",
				Message: "URL not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete URL",
		})
	}

	return c.JSON(model.ErrorResponse{
		Error:   "",
		Message: "URL deleted successfully",
	})
}

// ToggleURLActive godoc
// @Summary Toggle URL active status
// @Description Enable or disable a URL without deleting it
// @Tags admin
// @Produce json
// @Param code path string true "Short code"
// @Security BearerAuth
// @Success 200 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/admin/urls/{code}/toggle [post]
func (h *AdminHandler) ToggleURLActive(c *fiber.Ctx) error {
	code := c.Params("code")

	err := h.service.ToggleURLActive(code)
	if err != nil {
		if err == service.ErrURLNotFound {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
				Error:   "not_found",
				Message: "URL not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to toggle URL status",
		})
	}

	return c.JSON(model.ErrorResponse{
		Error:   "",
		Message: "URL status toggled successfully",
	})
}