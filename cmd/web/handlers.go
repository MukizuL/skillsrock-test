package main

import (
	"github.com/MukizuL/skillsrock-test/internal/validator"
	"github.com/gofiber/fiber/v2"
)

// CreateTaskRequest is used for deserializing data
type CreateTaskRequest struct {
	Title               string `json:"title"`
	Description         string `json:"description"`
	validator.Validator `json:"-"`
}

// UpdateTaskRequest is used for deserializing data
type UpdateTaskRequest struct {
	Title               string `json:"title"`
	Description         string `json:"description"`
	Status              string `json:"status"`
	validator.Validator `json:"-"`
}

// PostTask handler creates a task with given title and description and default status = "new"
func (app *application) PostTask(c *fiber.Ctx) error {
	var task CreateTaskRequest
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	task.CheckField(validator.NotBlank(task.Title), "title", "This field should not be blank")
	task.CheckField(validator.NotBlank(task.Description), "description", "This field should not be blank")

	if !task.Valid() {
		return c.Status(fiber.StatusBadRequest).JSON(task.FieldErrors)
	}

	err := app.tasks.Insert(task.Title, task.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert task",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task created successfully",
	})
}

// GetTask selects all tasks and serializes them as JSON
func (app *application) GetTask(c *fiber.Ctx) error {
	tasks, err := app.tasks.SelectAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get tasks",
		})
	}

	return c.Status(fiber.StatusOK).JSON(tasks)
}

// PutTask updates a specific task
func (app *application) PutTask(c *fiber.Ctx) error {
	var task UpdateTaskRequest
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	task.CheckField(validator.NotBlank(task.Title), "title", "This field should not be blank")
	task.CheckField(validator.NotBlank(task.Description), "description", "This field should not be blank")
	task.CheckField(validator.NotBlank(task.Status), "status", "This field should not be blank")

	task.CheckField(validator.PermittedValue(task.Status, "new", "in_progress", "done"), "status", "This value is not permitted")

	if !task.Valid() {
		return c.Status(fiber.StatusBadRequest).JSON(task.FieldErrors)
	}

	taskNew, err := app.tasks.Update(id, task.Title, task.Description, task.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update a task",
		})
	}

	return c.Status(fiber.StatusOK).JSON(taskNew)
}

// DeleteTask deletes a specific task
func (app *application) DeleteTask(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	err = app.tasks.Delete(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete a task",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task deleted successfully",
	})
}
