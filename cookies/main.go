package main

import (
	"fmt"
	"net/http"

	"github.com/ginjigo/ginji"
)

func main() {
	app := ginji.New()
	// Set a cookie
	app.Get("/set", func(c *ginji.Context) error {
		c.SetCookie(&http.Cookie{
			Name:  "session_token",
			Value: "abc123",
			Path:  "/",
		})
		return c.JSONOK(ginji.H{"message": "Cookie set"})
	})
	// Read a cookie
	app.Get("/get", func(c *ginji.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			return c.JSON(404, ginji.H{"error": "Cookie not found"})
		}
		return c.JSONOK(ginji.H{"cookie_value": cookie.Value})
	})

	if err := app.Run(":8083"); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
