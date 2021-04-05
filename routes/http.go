package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/leomirandadev/websocket-grpc-go/ws"
)

type Router interface {
	SetupRoutes()
	wsEndpoint(w http.ResponseWriter, r *http.Request)
}

type config struct {
	ws       *ws.WS
	upgrader websocket.Upgrader
}

func New(wsInit *ws.WS) Router {

	this := &config{
		ws: wsInit,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}

	this.SetupRoutes()

	return this
}

func (this *config) SetupRoutes() {
	http.HandleFunc("/ws", this.wsEndpoint)
}

func (this *config) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	this.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := this.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	channelWS := "default"
	channelQuery := r.URL.Query()["channel"]

	if channelQuery != nil {
		channelWS = channelQuery[0]
	}

	(*this.ws).SetClient(conn, channelWS)

	log.Println("Client Successfully connected...")

	(*this.ws).Reader(conn)
}
