package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Used to allow HTTP connection to use websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow everyone to connect
}

type Hub struct {
	clients    []*Client    // Clients to handle
	register   chan *Client // Channel to handle new client connection
	unregister chan *Client // Channel to handle client disconnect
	mutex      *sync.Mutex  // To avoid race conditions in our Hub
}

// Create a new HUB
func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

func (hub *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil) // Update socket connection

	if err != nil {
		log.Println(err)
		http.Error(w, "Couldn't upgrade socket", http.StatusBadRequest)
	}

	client := NewClient(hub, socket)
	hub.register <- client // Send Client to register channel

	go client.Write() // New routine in charge of sending messages to client
}
