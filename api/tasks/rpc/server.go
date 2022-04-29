package rpc

import (
	"context"

	"github.com/mtrqq/todo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func idFromProto(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

func idToProto(id primitive.ObjectID) string {
	return id.Hex()
}

func taskToProto(task todo.Task) *Task {
	return &Task{
		Id:        idToProto(task.ID),
		Name:      task.Name,
		Completed: task.Completed,
		Started:   timestamppb.New(task.Started),
	}
}

type tasksAPIServer struct {
	UnimplementedTasksAPIServer
	repository todo.TaskRepository
}

func NewTasksAPIServer(repository todo.TaskRepository) TasksAPIServer {
	return tasksAPIServer{repository: repository}
}

func (api tasksAPIServer) New(ctx context.Context, task *NewTaskRequest) (*ID, error) {
	repoTask := todo.NewTask(task.GetName())
	repoTask, err := api.repository.Add(ctx, repoTask)
	if err != nil {
		return nil, err
	}

	return &ID{Id: idToProto(repoTask.ID)}, nil
}
func (api tasksAPIServer) Get(ctx context.Context, id *ID) (*Task, error) {
	repoId, err := idFromProto(id.GetId())
	if err != nil {
		return nil, err
	}

	repoTask, err := api.repository.GetByID(ctx, repoId)
	if err != nil {
		return nil, err
	}

	return taskToProto(repoTask), nil
}
func (api tasksAPIServer) Complete(ctx context.Context, id *ID) (*ID, error) {
	repoId, err := idFromProto(id.GetId())
	if err != nil {
		return nil, err
	}

	repoId, err = api.repository.CompleteByID(ctx, repoId)
	if err != nil {
		return nil, err
	}

	return &ID{Id: idToProto(repoId)}, nil
}
func (api tasksAPIServer) List(ctx context.Context, _ *ListTasksRequest) (*ListTasksResponse, error) {
	repoTasks, err := api.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	tasks := make([]*Task, 0, len(repoTasks))
	for _, repoTask := range repoTasks {
		tasks = append(tasks, taskToProto(repoTask))
	}

	return &ListTasksResponse{Tasks: tasks}, nil
}
