package main

import (
	"context"
	"log"
	"time"

	services "github.com/leomirandadev/websocket-grpc-go/types/message"
	"google.golang.org/grpc"
)

var (
	serverURL = "localhost:10000"
)

func getGRPCClient() *grpc.ClientConn {
	var opts = []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	conn, err := grpc.Dial(serverURL, opts...)

	if err != nil {
		log.Fatal("Fail to dial %v", err)
	}

	return conn
}

func main() {
	conn := getGRPCClient()
	defer conn.Close()

	client := services.NewMessageServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	_, err := client.ReceiveMsg(ctx, &services.Message{User: "Leonardo GRPC", Message: "I'm here by grpc", Channel: "test"})

	if err != nil {
		log.Fatal(err)
	}

}
