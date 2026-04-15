package main

import (
	"fmt"
	"log"
	"strings"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// run initializes all dependencies and starts the HTTP server.
// Returning an error instead of calling log.Fatalf inside ensures that all
// deferred cleanup functions (e.g. db.Close) are executed before the process exits.
func run() error {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize services
	participantSvc := services.NewParticipantService(db)
	responseSvc := services.NewResponseService(db)
	scoreSvc := services.NewScoreService(db)
	scoringSvc := services.NewScoringService()

	// Initialize handlers
	handler := handlers.NewHandler(participantSvc, responseSvc, scoreSvc, scoringSvc)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10 MB
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.AllowedOrigins, ", "),
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Routes
	api := app.Group("/api")

	// Participant routes
	api.Post("/participants", handler.CreateParticipant)
	api.Get("/participants/:id", handler.GetParticipant)

	// Response routes
	api.Post("/responses", handler.SubmitResponse)

	// Scoring routes
	api.Post("/scoring", handler.CalculateScore)
	api.Get("/scores/:participantId", handler.GetScores)

	// Export routes
	api.Get("/export/:participantId", handler.ExportCSV)
	api.Get("/export", handler.ExportAllCSV)

	// Dashboard routing
	api.Get("/dashboard", handler.GetDashboardData)

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	return app.Listen(":" + cfg.ServerPort)
}
