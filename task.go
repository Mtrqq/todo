package todo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Completed bool               `bson:"completed" json:"completed"`
	Started   time.Time          `bson:"started" json:"started"`
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
