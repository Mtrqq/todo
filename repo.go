package todo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const databaseName = "todo"
const collectionName = "tasks"

type TaskRepository interface {
	GetByID(context.Context, primitive.ObjectID) (Task, error)
	Add(context.Context, Task) (Task, error)
	List(context.Context) ([]Task, error)
	Complete(context.Context, Task) (Task, error)
	CompleteByID(context.Context, primitive.ObjectID) (primitive.ObjectID, error)
}

type repository struct {
	client *mongo.Client
	tasks  *mongo.Collection
}

func NewTaskRepository(ctx context.Context, dbUrl string) (TaskRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUrl))
	if err != nil {
		return nil, err
	}

	tasks := client.Database(databaseName).Collection(collectionName)
	return &repository{client: client, tasks: tasks}, nil
}

func (repo *repository) Disconnect(ctx context.Context) error {
	return repo.client.Disconnect(ctx)
}

func (repo *repository) GetByID(ctx context.Context, id primitive.ObjectID) (Task, error) {
	singleResult := repo.tasks.FindOne(ctx, bson.D{{Key: "_id", Value: id}})
	if singleResult.Err() != nil {
		return Task{}, singleResult.Err()
	}

	var task Task
	if err := singleResult.Decode(&task); err != nil {
		return Task{}, err
	} else {
		return task, nil
	}
}

func (repo *repository) Add(ctx context.Context, task Task) (Task, error) {
	insertResult, err := repo.tasks.InsertOne(ctx, task)
	if err != nil {
		return Task{}, err
	}

	return Task{
		ID:        insertResult.InsertedID.(primitive.ObjectID),
		Name:      task.Name,
		Completed: task.Completed,
		Started:   task.Started,
	}, nil
}

func (repo *repository) List(ctx context.Context) ([]Task, error) {
	cursor, err := repo.tasks.Find(ctx, bson.D{})
	if err != nil {
		return []Task{}, nil
	}
	defer cursor.Close(ctx)

	tasks := make([]Task, 0)
	if err := cursor.All(ctx, &tasks); err != nil {
		return []Task{}, nil
	}

	return tasks, nil
}

func (repo *repository) Complete(ctx context.Context, task Task) (Task, error) {
	if task.Completed {
		return Task{}, errors.New("task is already completed")
	}

	updates := bson.D{{Key: "$set", Value: bson.D{{Key: "completed", Value: true}}}}
	_, err := repo.tasks.UpdateByID(ctx, task.ID, updates)
	if err != nil {
		return Task{}, err
	}

	return Task{
		ID:        task.ID,
		Name:      task.Name,
		Completed: true,
		Started:   task.Started,
	}, nil
}

func (repo *repository) CompleteByID(ctx context.Context, id primitive.ObjectID) (primitive.ObjectID, error) {
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
