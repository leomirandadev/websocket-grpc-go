package ws

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/leomirandadev/websocket-grpc-go/types"
)

type WS interface {
	Reader(conn *websocket.Conn)
	HandleMessages()
	SetClient(conn *websocket.Conn, channel string)
}

type config struct {
	clients   map[*websocket.Conn]string
	broadcast chan types.Message
}

func New() WS {
	return &config{
		clients:   make(map[*websocket.Conn]string),
		broadcast: make(chan types.Message),
	}
}

func (this *config) SetClient(conn *websocket.Conn, channel string) {
	this.clients[conn] = channel
}

func (this *config) Reader(conn *websocket.Conn) {
	for {
		var msg types.Message
		// Read in a new message as JSON and map it to a Message object
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("reader.error: %v", err)

			delete(this.clients, conn)

			break
		}
		// Send the newly received message to the broadcast channel
		this.broadcast <- msg
	}
}

func (this *config) HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-this.broadcast

		// Send it out to every client that is currently connected
		for client := range this.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("handleMessages.error: %v", err)
				client.Close()
				delete(this.clients, client)
			}
		}
	}
}
