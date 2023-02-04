package models

type WebsocketMessage struct {
	Type    string      `json:"type"`    // Type of message to send on websocket (direct, broadcast, etc)
	Payload interface{} `json:"payload"` // Message payload to send
}
