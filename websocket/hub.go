package websocket

import (
	"encoding/json"
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

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

// Show client connects and his Address
func (hub *Hub) onConnect(client *Client) {
	log.Println("Client connected", client.socket.RemoteAddr())

	// Lock hub to handle user connection before accept another connection
	hub.mutex.Lock()
	// Unlock hub at the end of connection
	defer hub.mutex.Unlock()

	client.id = client.socket.RemoteAddr().String() // Client ID is his address connection
	hub.clients = append(hub.clients, client)       // Add new client to slice
}

func (hub *Hub) onDisconnect(client *Client) {
	log.Println("Client disconnect", client.socket.RemoteAddr())

	// Close client connection
	client.socket.Close()

	// Lock channel to disconnect client before other clients disconnect
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	i := -1
	for index, clientData := range hub.clients {
		if clientData.id == client.id {
			i = index // Client index on slice
		}
	}

	copy(hub.clients[i:], hub.clients[i+1:])       // Copy without this specific entry (i)
	hub.clients[len(hub.clients)-1] = nil          // Last position on slice will be set to nil
	hub.clients = hub.clients[:len(hub.clients)-1] // New slice without last position
}

// Message send to every client except for ignoreClient specified
func (hub *Hub) Broadcast(message interface{}, ignoreClient *Client) {
	data, _ := json.Marshal(message)
	for _, client := range hub.clients {
		if client != ignoreClient {
			client.outbound <- data // Send message to outbound channel to send message to each client
		}
	}
}
