package todo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Completed bool               `bson:"completed"`
	Started   time.Time          `bson:"started"`
}

func NewTask(name string) Task {
	timestamp := time.Now().UTC()
	return Task{
		ID:        primitive.NewObjectIDFromTimestamp(timestamp),
		Name:      name,
		Completed: false,
		Started:   timestamp,
	}
}

func NewCompletedTask(name string) Task {
	timestamp := time.Now().UTC()
	return Task{
		ID:        primitive.NewObjectIDFromTimestamp(timestamp),
		Name:      name,
		Completed: true,
		Started:   timestamp,
	}
}
