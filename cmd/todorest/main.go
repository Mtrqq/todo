package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mtrqq/todo"
	taskrest "github.com/mtrqq/todo/api/tasks/rest"
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

	repository, err := todo.NewTaskRepository(context.Background(), repositoryUrl)
	if err != nil {
		log.Fatalf("Failed to establish repository connection: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", *host, *port)
	app := taskrest.NewAPIServer(repository)
	log.Printf("server listening at %v", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func init() {
	flag.Parse()
}
