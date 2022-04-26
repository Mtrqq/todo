package repo

import (
	"context"
	"errors"

	"github.com/mtrqq/todo/todo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const databaseName = "todo"
const collectionName = "tasks"

type TaskRepository struct {
	client *mongo.Client
	tasks  *mongo.Collection
}

func Connect(ctx context.Context, dbUrl string) (*TaskRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUrl))
	if err != nil {
		return nil, err
	}

	tasks := client.Database(databaseName).Collection(collectionName)
	return &TaskRepository{client: client, tasks: tasks}, nil
}

func (repo *TaskRepository) Disconnect(ctx context.Context) error {
	return repo.client.Disconnect(ctx)
}

func (repo *TaskRepository) GetByID(ctx context.Context, id primitive.ObjectID) (todo.Task, error) {
	singleResult := repo.tasks.FindOne(ctx, bson.D{{Key: "_id", Value: id}})
	if singleResult.Err() != nil {
		return todo.Task{}, singleResult.Err()
	}

	var task todo.Task
	if err := singleResult.Decode(&task); err != nil {
		return todo.Task{}, err
	} else {
		return task, nil
	}
}

func (repo *TaskRepository) Add(ctx context.Context, task todo.Task) (todo.Task, error) {
	insertResult, err := repo.tasks.InsertOne(ctx, task)
	if err != nil {
		return todo.Task{}, err
	}

	return todo.Task{
		ID:        insertResult.InsertedID.(primitive.ObjectID),
		Name:      task.Name,
		Completed: task.Completed,
		Started:   task.Started,
	}, nil
}

func (repo *TaskRepository) List(ctx context.Context) ([]todo.Task, error) {
	cursor, err := repo.tasks.Find(ctx, bson.D{})
	if err != nil {
		return []todo.Task{}, nil
	}
	defer cursor.Close(ctx)

	tasks := make([]todo.Task, 0)
	if err := cursor.All(ctx, &tasks); err != nil {
		return []todo.Task{}, nil
	}

	return tasks, nil
}

func (repo *TaskRepository) Complete(ctx context.Context, task todo.Task) (todo.Task, error) {
	if task.Completed {
		return todo.Task{}, errors.New("task is already completed")
	}

	updates := bson.D{{Key: "$set", Value: bson.D{{Key: "completed", Value: true}}}}
	_, err := repo.tasks.UpdateByID(ctx, task.ID, updates)
	if err != nil {
		return todo.Task{}, err
	}

	return todo.Task{
		ID:        task.ID,
		Name:      task.Name,
		Completed: true,
		Started:   task.Started,
	}, nil
}

func (repo *TaskRepository) CompleteByID(ctx context.Context, id primitive.ObjectID) (primitive.ObjectID, error) {
	task, err := repo.GetByID(ctx, id)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	if task, err = repo.Complete(ctx, task); err != nil {
		return primitive.ObjectID{}, err
	} else {
		return task.ID, nil
	}
}
