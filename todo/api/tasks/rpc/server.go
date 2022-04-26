package rpc

import (
	"context"

	"github.com/mtrqq/todo/todo"
	"github.com/mtrqq/todo/todo/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func idFromProto(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

func idToProto(id primitive.ObjectID) string {
	return id.Hex()
}

func taskFromProto(task *Task) (todo.Task, error) {
	id, err := idFromProto(task.GetId())
	if err != nil {
		return todo.Task{}, err
	}

	return todo.Task{
		ID:        id,
		Name:      task.GetName(),
		Completed: task.GetCompleted(),
		Started:   task.GetStarted().AsTime(),
	}, nil
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
	repository *repo.TaskRepository
}

func NewTasksAPIServer(ctx context.Context, repositoryUrl string) (TasksAPIServer, error) {
	repository, err := repo.Connect(ctx, repositoryUrl)
	if err != nil {
		return tasksAPIServer{}, err
	}

	return tasksAPIServer{repository: repository}, nil
}

func (api tasksAPIServer) New(ctx context.Context, task *Task) (*ID, error) {
	repoTask, err := taskFromProto(task)
	if err != nil {
		return nil, err
	}

	repoTask, err = api.repository.Add(ctx, repoTask)
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
