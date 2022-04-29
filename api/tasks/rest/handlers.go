package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mtrqq/todo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func listTasks(repo todo.TaskRepository, ctx *fiber.Ctx) error {
	tasks, err := repo.List(ctx.Context())
	if err != nil {
		return fiber.NewError(500, "failed to fetch tasks list")
	}

	return ctx.JSON(tasks)
}

func newTask(repo todo.TaskRepository, ctx *fiber.Ctx) error {
	var newTaskRequest struct {
		Name string `json:"name"`
	}
	if err := ctx.BodyParser(&newTaskRequest); err != nil {
		return fiber.NewError(403, "invalid request body")
	}

	task := todo.NewTask(newTaskRequest.Name)
	if task, err := repo.Add(ctx.Context(), task); err != nil {
		return fiber.NewError(500, "failed to create task")
	} else {
		return ctx.JSON(task)
	}
}

func completeTask(repo todo.TaskRepository, ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(403, "id not specified")
	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fiber.NewError(403, "malformed id specified")
	}

	objId, err = repo.CompleteByID(ctx.Context(), objId)
	if err != nil {
		return fiber.NewError(500, "failed to complete task")
	}

	return ctx.JSON(fiber.Map{"id": objId.Hex()})
}

func getTask(repo todo.TaskRepository, ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(403, "id not specified")
	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fiber.NewError(403, "malformed id specified")
	}

	task, err := repo.GetByID(ctx.Context(), objId)
	if err != nil {
		return fiber.NewError(500, "failed to lookup task")
	}

	return ctx.JSON(task)
}
