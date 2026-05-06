package main

import (
	"fmt"
	"log"

	"papa-shortener/internal/config"
	"papa-shortener/internal/handler"
	"papa-shortener/internal/middleware"
	"papa-shortener/internal/model"
	"papa-shortener/internal/repository"
	"papa-shortener/internal/service"
	"papa-shortener/pkg/validator"

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
	if err := db.AutoMigrate(&model.URL{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize layers
	repo := repository.NewURLRepository(db)
	svc := service.NewURLService(repo, config.AppCfg)
	h := handler.NewURLHandler(svc)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal_error",
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(middleware.Logger())

	// Routes
	app.Get("/health", h.HealthCheck)

	// API routes
	api := app.Group("/api")
	api.Post("/shorten", h.CreateShortURL)

	// Redirect route (must be last)
	app.Get("/:code", h.ResolveURL)

	// Start server
	addr := fmt.Sprintf("%s:%d", config.AppCfg.App.Host, config.AppCfg.App.Port)
	log.Printf("Starting server on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}