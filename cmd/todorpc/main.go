package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/mtrqq/todo"
	taskrpc "github.com/mtrqq/todo/api/tasks/rpc"
	"google.golang.org/grpc"
)

var (
	repositoryUrl = os.Getenv("TODO_REPOSITORY_URL")
	host          = flag.String("host", "", "Server host")
	port          = flag.Int("port", 30030, "Server port")
)

func main() {
	if repositoryUrl == "" {
		log.Fatalf("Consider setting TODO_REPOSITORY_URL environment variable")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	repository, err := todo.NewTaskRepository(context.Background(), repositoryUrl)
	if err != nil {
		log.Fatalf("Failed to establish repository connection: %v", err)
	}

	api := taskrpc.NewTasksAPIServer(repository)
	server := grpc.NewServer()
	taskrpc.RegisterTasksAPIServer(server, api)
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func init() {
	flag.Parse()
}
