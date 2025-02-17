package main

import "github.com/gofiber/fiber/v2"

// routes initializes new fiber server with routes
func (app *application) routes() *fiber.App {
	f := fiber.New()

	f.Use(app.logRequest)

	f.Post("/tasks", app.PostTask)
	f.Get("/tasks", app.GetTask)
	f.Put("/tasks/:id", app.PutTask)
	f.Delete("/tasks/:id", app.DeleteTask)

	return f
}
