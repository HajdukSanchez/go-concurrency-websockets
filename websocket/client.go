package websocket

import "github.com/gorilla/websocket"

type Client struct {
	hub      *Hub            // Hub of messages
	id       string          // Client id
	socket   *websocket.Conn // Socket connection for specific client
	outbound chan []byte     // Channel to handle Messages to be send
}

func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
	}
}

func (client *Client) Write() {
	for {
		select {
		case message, ok := <-client.outbound:
			if !ok {
				client.socket.WriteMessage(websocket.CloseMessage, []byte{}) // If something wrong, send an error message
				return
			}
			client.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
