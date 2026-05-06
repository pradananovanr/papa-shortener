package main

import (
	"fmt"
	"html/template"
	"io"
	"log"

	"github.com/prdnnvnrnt/papa-shortener/internal/config"
	"github.com/prdnnvnrnt/papa-shortener/internal/handler"
	"github.com/prdnnvnrnt/papa-shortener/internal/middleware"
	"github.com/prdnnvnrnt/papa-shortener/internal/model"
	"github.com/prdnnvnrnt/papa-shortener/internal/repository"
	"github.com/prdnnvnrnt/papa-shortener/internal/service"
	"github.com/prdnnvnrnt/papa-shortener/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load config
	if err := config.Load("config.yaml"); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize validator
	validator.Init()

	// Connect to database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.AppCfg.Database.Host,
		config.AppCfg.Database.Port,
		config.AppCfg.Database.User,
		config.AppCfg.Database.Password,
		config.AppCfg.Database.Name,
		config.AppCfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(&model.URL{}, &model.Admin{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize layers
	urlRepo := repository.NewURLRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	urlSvc := service.NewURLService(urlRepo, config.AppCfg)
	authSvc := service.NewAuthService(adminRepo, config.AppCfg)
	adminSvc := service.NewAdminService(urlRepo, config.AppCfg)

	urlHandler := handler.NewURLHandler(urlSvc)
	authHandler := handler.NewAuthHandler(authSvc)
	adminHandler := handler.NewAdminHandler(adminSvc)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "internal_error",
				"message": err.Error(),
			})
		},
		Views: NewTemplateEngine(),
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New())

	// Routes
	app.Get("/health", urlHandler.HealthCheck)

	// HTML pages (order matters! /admin must be before /:code)
	app.Get("/", urlHandler.RenderIndex)
	app.Get("/admin/login", authHandler.RenderLogin)
	app.Get("/admin", adminHandler.RenderDashboard)

	// Auth routes (public)
	api := app.Group("/api")
	api.Post("/auth/login", authHandler.Login)
	api.Post("/auth/register", authHandler.RegisterAdmin)

	// Shorten URL (public)
	api.Post("/shorten", urlHandler.CreateShortURL)

	// Admin API routes (protected)
	admin := api.Group("/admin", middleware.AuthRequired())
	admin.Get("/urls", adminHandler.ListURLs)
	admin.Get("/urls/:code", adminHandler.GetURL)
	admin.Put("/urls/:code", adminHandler.UpdateURL)
	admin.Delete("/urls/:code", adminHandler.DeleteURL)
	admin.Post("/urls/:code/toggle", adminHandler.ToggleURLActive)

	// Redirect route (must be LAST - catches all remaining paths)
	app.Get("/:code", urlHandler.ResolveURL)

	// Start server
	addr := fmt.Sprintf("%s:%d", config.AppCfg.App.Host, config.AppCfg.App.Port)
	log.Printf("Starting server on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// TemplateEngine implements fiber's ViewEngine for Go templates
type TemplateEngine struct {
	templates *template.Template
}

func NewTemplateEngine() *TemplateEngine {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	return &TemplateEngine{templates: tmpl}
}

func (e *TemplateEngine) Render(w io.Writer, name string, bind interface{}, layout ...string) error {
	return e.templates.ExecuteTemplate(w, name+".html", bind)
}

func (e *TemplateEngine) Load() error {
	var err error
	e.templates, err = template.ParseGlob("templates/*.html")
	return err
}
