package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	taskrpc "github.com/mtrqq/todo/todo/api/tasks/rpc"
)

var (
	addr = flag.String("addr", "localhost:30030", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := taskrpc.NewTasksAPIClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := client.List(ctx, &taskrpc.ListTasksRequest{})
	if err != nil {
		log.Fatalf("failed to fetch list: %v", err)
	}

	fmt.Println("todo list:")
	for _, task := range resp.GetTasks() {
		fmt.Println(task.GetId(), task.GetName(), task.GetCompleted(), task.GetStarted())
	}
}
