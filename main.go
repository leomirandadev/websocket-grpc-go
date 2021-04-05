package main

import (
	"log"
	"net/http"

	"github.com/leomirandadev/websocket-grpc-go/events"
	"github.com/leomirandadev/websocket-grpc-go/routes"
	"github.com/leomirandadev/websocket-grpc-go/ws"
)

func main() {
	log.Println("Initialization Websocket")

	wsInit := ws.New()
	go wsInit.HandleMessages()

	routes.New(&wsInit)

	events.New(&wsInit)

	log.Fatal(http.ListenAndServe(":0504", nil))
}
