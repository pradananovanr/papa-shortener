package handler

import (
	"errors"

	"papa-shortener/internal/model"
	"papa-shortener/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type URLHandler struct {
	service service.URLService
	validate *validator.Validate
}

func NewURLHandler(svc service.URLService) *URLHandler {
	return &URLHandler{
		service:  svc,
		validate: validator.New(),
	}
}

// CreateShortURL godoc
// @Summary Create a shortened URL
// @Description Create a new short URL, optionally with a custom code
// @Tags url
// @Accept json
// @Produce json
// @Param request body model.CreateURLRequest true "URL creation request"
// @Success 201 {object} model.CreateURLResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/shorten [post]
func (h *URLHandler) CreateShortURL(c *fiber.Ctx) error {
	var req model.CreateURLRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "validation_error",
			Message: formatValidationErrors(err),
		})
	}

	result, err := h.service.CreateShortURL(&req)
	if err != nil {
		if errors.Is(err, service.ErrCustomURLExists) {
			return c.Status(fiber.StatusConflict).JSON(model.ErrorResponse{
				Error:   "custom_url_exists",
				Message: "Custom URL is already taken",
			})
		}
		if errors.Is(err, service.ErrInvalidCustomURL) {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Error:   "invalid_custom_url",
				Message: "Custom URL must be 3-20 lowercase alphanumeric characters only",
			})
		}
		if errors.Is(err, service.ErrInvalidOriginalURL) {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Error:   "invalid_original_url",
				Message: "Invalid original URL format",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create short URL",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// ResolveURL godoc
// @Summary Redirect to original URL
// @Description Resolve a short code and redirect to the original URL
// @Tags url
// @Produce redirect
// @Param code path string true "Short code"
// @Success 302 {string} string "Redirect to original URL"
// @Failure 404 {string} string "URL not found"
// @Router /{code} [get]
func (h *URLHandler) ResolveURL(c *fiber.Ctx) error {
	code := c.Params("code")

	original, err := h.service.ResolveShortCode(code)
	if err != nil {
		if errors.Is(err, service.ErrURLNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
				Error:   "not_found",
				Message: "Short URL not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to resolve URL",
		})
	}

	return c.Redirect(original, fiber.StatusTemporaryRedirect)
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the service is running
// @Tags health
// @Produce json
// @Success 200 {object} model.HealthResponse
// @Router /health [get]
func (h *URLHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(model.HealthResponse{
		Status: "ok",
	})
}

func formatValidationErrors(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				return e.Field() + " is required"
			case "url":
				return e.Field() + " must be a valid URL"
			case "min":
				return e.Field() + " is too short"
			case "max":
				return e.Field() + " is too long"
			case "alphanum":
				return e.Field() + " must be alphanumeric"
			}
		}
	}
	return "Validation failed"
}