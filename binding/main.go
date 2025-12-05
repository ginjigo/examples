package main

import (
	"fmt"

	"github.com/ginjigo/ginji"
)

type User struct {
	Name  string `json:"name" ginji:"required,min=3"`
	Email string `json:"email" ginji:"required,email"`
	Age   int    `json:"age" ginji:"min=18"`
}

type SearchQuery struct {
	Q    string `query:"q" ginji:"required"`
	Page int    `query:"page"`
}

func main() {
	app := ginji.New()

	// POST /bind-json - Bind and validate JSON
	app.Post("/bind-json", func(c *ginji.Context) error {
		var input struct {
			Name string `json:"name" ginji:"required,min=3"`
		}
		if err := c.BindValidate(&input); err != nil {
			return c.JSON(400, ginji.H{"error": err.Error()})
		}
		return c.JSONOK(input)
	})
	// GET /query-params - Test query parameter binding
	app.Get("/query-params", func(c *ginji.Context) error {
		var query struct {
			Page int `query:"page"`
		}
		if err := c.BindValidate(&query); err != nil {
			return c.JSON(400, ginji.H{"error": err.Error()})
		}
		return c.JSONOK(query)
	})

	fmt.Println("Server running on :8082")
	if err := app.Listen(":8082"); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
