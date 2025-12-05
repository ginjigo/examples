package main

import (
	"github.com/ginjigo/ginji"
	"github.com/ginjigo/schema"
)

func main() {
	app := ginji.New()

	// Middleware
	app.Use(ginji.Logger())
	app.Use(ginji.Recovery())
	app.Use(ginji.DefaultErrorHandler())

	// Define schemas
	userSchema := schema.NewSchema(map[string]schema.Field{
		"email": *schema.String().Required().IsEmail().Describe("User email address"),
		"name":  *schema.String().Required().MinLength(2).MaxLength(50),
		"age":   *schema.Integer().Min(18).Max(120),
		"role":  *schema.String().Enum("admin", "user", "guest").Default("user"),
	}).Describe("User registration schema")

	updateUserSchema := schema.NewSchema(map[string]schema.Field{
		"name": *schema.String().MinLength(2).MaxLength(50),
		"age":  *schema.Integer().Min(18).Max(120),
	})

	// Routes with  schema validation
	app.Post("/users", func(c *ginji.Context) error {
		// The schema validation happens automatically before this handler
		var user map[string]any
		if err := c.BindJSON(&user); err != nil {
			return c.Fail(400, "Invalid JSON")
		}

		// Validate against schema
		var errors []schema.ValidationError
		for field, value := range user {
			if fieldSchema, ok := userSchema.Properties[field]; ok {
				errors = append(errors, fieldSchema.Validate(value, field)...)
			}
		}

		if len(errors) > 0 {
			return c.FailWithData(422, "Validation failed", ginji.H{
				"errors": errors,
			})
		}

		// User is valid, process registration
		return c.JSONOK(ginji.H{
			"message": "User created successfully",
			"user":    user,
		})
	}).
		Body(userSchema).
		Summary("Create a new user").
		Tags("users")

	app.Get("/users/:id", func(c *ginji.Context) error {
		id := c.Request.Param("id")
		role := c.Request.QueryDefault("role", "all")

		return c.JSONOK(ginji.H{
			"id":      id,
			"role":    role,
			"message": "User retrieved",
		})
	}).
		Summary("Get user by ID").
		Tags("users")

	app.Put("/users/:id", func(c *ginji.Context) error {
		id := c.Request.Param("id")

		var updates map[string]any
		if err := c.BindJSON(&updates); err != nil {
			return c.Fail(400, "Invalid JSON")
		}

		// Validate updates
		var errors []schema.ValidationError
		for field, value := range updates {
			if fieldSchema, ok := updateUserSchema.Properties[field]; ok {
				errors = append(errors, fieldSchema.Validate(value, field)...)
			}
		}

		if len(errors) > 0 {
			return c.FailWithData(422, "Validation failed", ginji.H{
				"errors": errors,
			})
		}

		return c.JSONOK(ginji.H{
			"id":      id,
			"updates": updates,
			"message": "User updated successfully",
		})
	}).
		Body(updateUserSchema).
		Summary("Update user").
		Tags("users")

	// Example with nested object schema
	addressSchema := schema.NewSchema(map[string]schema.Field{
		"street":  *schema.String().Required(),
		"city":    *schema.String().Required(),
		"state":   *schema.String().MinLength(2).MaxLength(2),
		"zip":     *schema.String().Pattern(`^\d{5}$`),
		"country": *schema.String().Default("US"),
	})

	profileSchema := schema.NewSchema(map[string]schema.Field{
		"name":    *schema.String().Required(),
		"email":   *schema.String().Required().IsEmail(),
		"address": *schema.Object(addressSchema.Properties),
		"tags":    *schema.Array(schema.String().MinLength(2)),
	})

	app.Post("/profiles", func(c *ginji.Context) error {
		var profile map[string]any
		if err := c.BindJSON(&profile); err != nil {
			return c.Fail(400, "Invalid JSON")
		}

		// Validate full profile including nested address
		var errors []schema.ValidationError
		for field, value := range profile {
			if fieldSchema, ok := profileSchema.Properties[field]; ok {
				errors = append(errors, fieldSchema.Validate(value, field)...)
			}
		}

		if len(errors) > 0 {
			return c.FailWithData(422, "Validation failed", ginji.H{
				"errors": errors,
			})
		}

		return c.JSONOK(ginji.H{
			"message": "Profile created",
			"profile": profile,
		})
	}).
		Body(profileSchema).
		Summary("Create user profile with nested address").
		Tags("profiles")

	// Health check
	app.Get("/health", func(c *ginji.Context) error {
		return c.TextOK("OK")
	})

	println("🚀 Schema Validation Example running on :3000")
	println("")
	println("Try these examples:")
	println("")
	println("  Valid user:")
	println(`  curl -X POST http://localhost:3000/users -H "Content-Type: application/json" \`)
	println(`    -d '{"email":"test@example.com","name":"John Doe","age":25,"role":"admin"}'`)
	println("")
	println("  Invalid email:")
	println(`  curl -X POST http://localhost:3000/users -H "Content-Type: application/json" \`)
	println(`    -d '{"email":"not-an-email","name":"John","age":25}'`)
	println("")
	println("  Missing required field:")
	println(`  curl -X POST http://localhost:3000/users -H "Content-Type: application/json" \`)
	println(`    -d '{"name":"John","age":25}'`)
	println("")
	println("  Invalid age (too young):")
	println(`  curl -X POST http://localhost:3000/users -H "Content-Type: application/json" \`)
	println(`    -d '{"email":"test@example.com","name":"John","age":15}'`)
	println("")
	println("  Update user:")
	println(`  curl -X PUT http://localhost:3000/users/123 -H "Content-Type: application/json" \`)
	println(`    -d '{"name":"Jane Doe","age":30}'`)
	println("")
	println("  Get user:")
	println("  curl http://localhost:3000/users/123?role=admin")

	_ = app.Listen(":3000")
}
