package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mtrqq/todo"
)

func wrapHandler(repo todo.TaskRepository, handler func(todo.TaskRepository, *fiber.Ctx) error) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return handler(repo, ctx)
	}
}

func NewAPIServer(repo todo.TaskRepository) *fiber.App {
	app := fiber.New()

	apiV1 := app.Group("/api/v1")
	{
		apiV1.Get("/tasks", wrapHandler(repo, listTasks)).Name("ListTasks")
		apiV1.Post("/tasks", wrapHandler(repo, newTask)).Name("NewTask")
		apiV1.Get("/tasks/:id", wrapHandler(repo, getTask)).Name("GetTask")
		apiV1.Post("/tasks/:id/complete", wrapHandler(repo, completeTask)).Name("CompleteTask")
	}

	return app
}
