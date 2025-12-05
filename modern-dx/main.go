package main

import (
	"github.com/ginjigo/ginji"
)

func main() {
	app := ginji.New()

	// Global middleware
	app.Use(ginji.Logger())
	app.Use(ginji.Recovery())

	// Demonstrate new DX improvements

	// 1. Using Request wrapper for cleaner API
	app.Get("/user/:id", func(c *ginji.Context) error {
		// New way - cleaner namespace
		id := c.Request.Param("id")
		name := c.Request.Query("name")
		role := c.Request.QueryDefault("role", "user")

		// Convenience method - JSONOK instead of JSON(200, ...)
		return c.JSONOK(ginji.H{
			"id":   id,
			"name": name,
			"role": role,
		})
	})

	// 2. Using Var for context storage (cleaner than Set/Get)
	app.Use(func(c *ginji.Context) error {
		// Store user info in context
		c.Var("user_id", 123)
		c.Var("tenant", "acme-corp")
		return c.Next()
	})

	app.Get("/profile", func(c *ginji.Context) error {
		if userID, exists := c.GetVar("user_id"); exists {
			return c.JSONOK(ginji.H{
				"user_id": userID,
				"message": "Profile loaded",
			})
		}
		return c.Fail(404, "Profile not found")
	})

	// 3. Using Fail for quick error responses
	app.Get("/protected", func(c *ginji.Context) error {
		token := c.Request.Query("token")
		if token == "" {
			// Quick error response
			return c.Fail(401, "Unauthorized - token required")
		}

		return c.JSONOK(ginji.H{"data": "protected content"})
	})

	// 4. Using FailWithData for richer error responses
	app.Post("/validate", func(c *ginji.Context) error {
		var input struct {
			Email string `json:"email" ginji:"required,email"`
			Age   int    `json:"age" ginji:"required,min=18"`
		}

		if err := c.BindValidate(&input); err != nil {
			// Rich error response with validation details
			return c.FailWithData(400, "Validation failed", ginji.H{
				"validation_error": err.Error(),
			})
		}

		return c.JSONOK(ginji.H{"status": "valid"})
	})

	// 5. Using convenience methods: TextOK, HTMLOK
	app.Get("/health", func(c *ginji.Context) error {
		return c.TextOK("OK")
	})

	app.Get("/welcome", func(c *ginji.Context) error {
		return c.HTMLOK("<h1>Welcome to Ginji!</h1><p>Clean, fast, and developer-friendly.</p>")
	})

	// 6. Backwards compatibility - old API still works
	app.Get("/old-style/:id", func(c *ginji.Context) error {
		id := c.Param("id")         // Old way still works
		name := c.Query("name")     // Old way still works
		return c.JSON(200, ginji.H{ // Old way still works
			"id":   id,
			"name": name,
		})
	})

	println("🚀 Server with enhanced DX running on :3000")
	println("")
	println("Try these examples:")
	println("  GET  /user/123?name=John&role=admin")
	println("  GET  /profile")
	println("  GET  /protected")
	println("  GET  /protected?token=abc123")
	println("  POST /validate -d '{\"email\":\"test@example.com\",\"age\":25}'")
	println("  GET  /health")
	println("  GET  /welcome")
	println("  GET  /old-style/456?name=Jane")

	_ = app.Listen(":3000")
}
