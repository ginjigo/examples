package main

import (
	"log/slog"
	"os"

	"github.com/ginjigo/ginji"
	"github.com/ginjigo/middleware"
)

func main() {
	// Create app in debug mode (default)
	app := ginji.New()

	// The logger is already initialized automatically based on mode
	// Debug mode: JSON logs with debug level
	// Release mode: JSON logs with info level

	// Optional: Customize logger
	// app.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 	Level: slog.LevelDebug,
	// }))

	// Add logger middleware to automatically log all requests
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		// Skip health check endpoint to reduce noise
		SkipPaths: []string{"/health"},
		// Could also skip based on custom logic:
		// SkipFunc: func(c *ginji.Context) bool {
		//     return c.Header("X-Skip-Log") == "true"
		// },
	}))

	// Routes
	app.Get("/", func(c *ginji.Context) error {
		// Log custom messages using the engine's logger
		app.Logger.Info("Homepage accessed")
		return c.JSON(200, ginji.H{"message": "Hello, World!"})
	})

	app.Get("/user/:id", func(c *ginji.Context) error {
		userID := c.Param("id")

		// Structured logging with context
		app.Logger.Info("User profile viewed",
			slog.String("user_id", userID),
			slog.String("ip", c.Req.RemoteAddr),
		)

		return c.JSON(200, ginji.H{
			"user_id": userID,
			"name":    "John Doe",
		})
	})

	app.Post("/users", func(c *ginji.Context) error {
		var user struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		if err := c.BindJSON(&user); err != nil {
			app.Logger.Error("Failed to bind user data",
				slog.String("error", err.Error()),
			)
			return c.JSON(400, ginji.H{"error": "Invalid request"})
		}

		// Log successful operations
		app.Logger.Info("New user created",
			slog.String("name", user.Name),
			slog.String("email", user.Email),
		)

		return c.JSON(201, ginji.H{"message": "User created", "user": user})
	})

	// Error endpoint to test error-level logging
	app.Get("/error", func(c *ginji.Context) error {
		app.Logger.Error("Intentional error triggered",
			slog.String("path", c.Req.URL.Path),
		)
		return c.JSON(500, ginji.H{"error": "Internal server error"})
	})

	// Health check (won't be logged due to SkipPaths)
	app.Get("/health", func(c *ginji.Context) error {
		return c.JSON(200, ginji.H{"status": "healthy"})
	})

	// Start server
	app.Logger.Info("Starting server", slog.String("port", "3000"))

	if err := app.Listen(":3000"); err != nil {
		app.Logger.Error("Server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
