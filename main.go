package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]string)

var broadcast = make(chan Message)

type Message struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

func main() {
	log.Println("Initialization Websocket")

	setupRoutes()
	go handleMessages()

	log.Fatal(http.ListenAndServe(":0504", nil))
}

func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/ws/", wsEndpoint)
}

func reader(conn *websocket.Conn) {
	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("reader.error: %v", err)
			delete(clients, conn)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("handleMessages.error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	clients[ws] = "default"
	log.Println("Client Successfully connected...")

	reader(ws)
}
