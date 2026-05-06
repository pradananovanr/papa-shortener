package handler

import (
	"errors"
	"fmt"

	"github.com/prdnnvnrnt/papa-shortener/internal/config"
	"github.com/prdnnvnrnt/papa-shortener/internal/model"
	"github.com/prdnnvnrnt/papa-shortener/internal/service"

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
		// Check if HTMX request - return HTML error
		if c.Get("HX-Request") == "true" {
			var msg string
			status := fiber.StatusBadRequest
			if errors.Is(err, service.ErrCustomURLExists) {
				msg = `<div class="mt-4 p-4 bg-red-500/20 border border-red-500/50 rounded-xl text-red-200">
					<p class="font-semibold">❌ Custom URL sudah dipakai!</p>
					<p class="text-sm mt-1">Coba pakai kombinasi lain.</p>
				</div>`
				status = fiber.StatusConflict
			} else if errors.Is(err, service.ErrInvalidCustomURL) {
				msg = `<div class="mt-4 p-4 bg-red-500/20 border border-red-500/50 rounded-xl text-red-200">
					<p class="font-semibold">❌ Custom URL tidak valid</p>
					<p class="text-sm mt-1">3-20 karakter huruf kecil & angka saja.</p>
				</div>`
			} else {
				msg = `<div class="mt-4 p-4 bg-red-500/20 border border-red-500/50 rounded-xl text-red-200">
					<p class="font-semibold">❌ Gagal membuat short URL</p>
				</div>`
			}
			return c.Status(status).SendString(msg)
		}

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

	// Check if HTMX request - return HTML success
	if c.Get("HX-Request") == "true" {
		badge := "🎲 Random"
		if result.IsCustom {
			badge = "✨ Custom"
		}
		html := fmt.Sprintf(`<div class="mt-4 p-6 bg-green-500/20 border border-green-500/50 rounded-xl">
			<div class="flex items-center justify-between mb-3">
				<span class="text-green-200 font-semibold">%s</span>
				<button onclick="copyToClipboard('%s')" class="text-2xl hover:scale-110 transition">📋</button>
			</div>
			<div class="bg-white/10 rounded-lg p-3 mb-3">
				<a href="%s" class="text-primary font-mono text-lg break-all hover:underline">%s</a>
			</div>
			<div class="text-sm text-green-200">
				<p>Original: <span class="opacity-75">%s</span></p>
			</div>
		</div>
		<script>hljs.highlightAll();</script>`, badge, result.FullShort, result.FullShort, result.FullShort, result.Original)
		return c.Status(fiber.StatusCreated).SendString(html)
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

// RenderIndex renders the index HTML page
func (h *URLHandler) RenderIndex(c *fiber.Ctx) error {
	data := fiber.Map{
		"BaseURL": config.AppCfg.App.BaseURL,
	}
	return c.Render("index", data)
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