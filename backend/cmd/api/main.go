package main

import (
	"log"
	"os"

	"github.com/dejavu/backend/internal/handler"
	"github.com/dejavu/backend/pkg/cache"
	"github.com/dejavu/backend/pkg/database"
	"github.com/dejavu/backend/pkg/queue"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize connections
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	redis := cache.Connect()
	defer redis.Close()

	nats, err := queue.Connect()
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nats.Close()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Dejavu API",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"service": "dejavu-api",
		})
	})

	// Initialize handlers
	authHandler := handler.NewAuthHandler(db, redis)
	projectHandler := handler.NewProjectHandler(db)
	deployHandler := handler.NewDeployHandler(db, nats)

	// Routes
	api := app.Group("/api")
	
	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected routes
	api.Use(handler.AuthMiddleware(redis))

	// Project routes
	projects := api.Group("/projects")
	projects.Get("/", projectHandler.List)
	projects.Post("/", projectHandler.Create)
	projects.Get("/:id", projectHandler.Get)
	projects.Put("/:id", projectHandler.Update)
	projects.Delete("/:id", projectHandler.Delete)

	// Deployment routes
	deploy := api.Group("/deploy")
	deploy.Post("/", deployHandler.Trigger)
	deploy.Get("/:id", deployHandler.GetStatus)
	deploy.Get("/:id/logs", deployHandler.StreamLogs)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

