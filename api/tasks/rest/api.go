package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mtrqq/todo"
)

func NewAPIServer(repo todo.TaskRepository) *fiber.App {
	app := fiber.New()

	repoProxy := &repositoryProxy{repo: repo}
	apiV1 := app.Group("/api/v1")
	{
		apiV1.Get("/tasks", repoProxy.ListTasks).Name("ListTasks")
		apiV1.Post("/tasks", repoProxy.NewTask).Name("NewTask")
		apiV1.Get("/tasks/:id", repoProxy.GetTask).Name("GetTask")
		apiV1.Post("/tasks/:id/complete", repoProxy.CompleteTask).Name("CompleteTask")
	}

	return app
}
