package main

import (
	"log"
	"time"

	"github.com/ginjigo/ginji"
)

func main() {
	app := ginji.New()

	// Register some routes
	app.Get("/", func(c *ginji.Context) error {
		if err := c.JSON(200, ginji.H{"message": "Hello, World!"}); err != nil {
			log.Printf("Error sending JSON response: %v", err)
		}
		return nil
	})

	app.Get("/ping", func(c *ginji.Context) error {
		if err := c.JSON(200, ginji.H{"status": "pong"}); err != nil {
			log.Printf("Error sending JSON response: %v", err)
		}
		return nil
	})

	// Simulate long-running requests
	app.Get("/long", func(c *ginji.Context) error {
		log.Println("Starting long request...")
		time.Sleep(5 * time.Second)
		log.Println("Long request completed")
		return c.JSON(200, ginji.H{"message": "Completed after 5 seconds"})
	})

	// Start server with graceful shutdown
	//
	// The server will:
	// 1. Listen for SIGINT (Ctrl+C) or SIGTERM signals
	// 2. When received, stop accepting new connections
	// 3. Wait up to 10 seconds for in-flight requests to complete
	// 4. Cleanly shut down all plugins
	// 5. Exit gracefully
	//
	// To test:
	// 1. Run: go run examples/graceful-shutdown/main.go
	// 2. In another terminal: curl http://localhost:3000/long
	// 3. While the request is running, press Ctrl+C
	// 4. Observe that the long request completes before shutdown
	log.Println("Starting server on :3000 with graceful shutdown...")
	log.Println("Press Ctrl+C to trigger graceful shutdown")

	if err := app.ListenWithShutdown(":3000", 10*time.Second); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
