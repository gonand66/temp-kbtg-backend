package main

import (
	"log"
	"temp-kbtg-backend/database"
	"temp-kbtg-backend/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "KBTG Backend API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// User routes
	app.Get("/users", handlers.GetAllUsers)
	app.Get("/users/:id", handlers.GetUserByID)
	app.Post("/users", handlers.CreateUser)
	app.Put("/users/:id", handlers.UpdateUser)
	app.Delete("/users/:id", handlers.DeleteUser)

	// Start server
	log.Println("🚀 Server starting on port 3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
