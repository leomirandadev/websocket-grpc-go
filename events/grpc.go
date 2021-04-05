package events

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/leomirandadev/websocket-grpc-go/types"
	"github.com/leomirandadev/websocket-grpc-go/types/message"
	"github.com/leomirandadev/websocket-grpc-go/ws"
	"google.golang.org/grpc"
)

var (
	serverURL  = "localhost:10000"
	portServer = 10000
)

type Server struct {
	ws *ws.WS
}

func New(wsInit *ws.WS) {

	log.Println("GRPC running on", serverURL)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", portServer))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := Server{ws: wsInit}

	grpcServer := grpc.NewServer()
	message.RegisterMessageServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("GRPC failed to serve: %s", err)
	}
}

func (s *Server) ReceiveMsg(ctx context.Context, in *message.Message) (*message.Result, error) {

	(*s.ws).SendMessage(types.Message{
		User:    in.User,
		Message: in.Message,
		Channel: in.Channel,
	})

	return &message.Result{
		Ok: true,
	}, nil
}
