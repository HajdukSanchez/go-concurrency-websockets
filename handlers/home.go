package handlers

import (
	"encoding/json"
	"net/http"

	"hajduksanchez.com/go/rest-websockets/server"
)

// Response to return to client
type HomeResponse struct {
	Message string `json:"message"` // This is a serialization to handle the value on JSON converting
	Status  bool   `json:"status"`
}

func HomeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json") // Specify the response is a json response
		w.WriteHeader(http.StatusOK)                       // 200

		json.NewEncoder(w).Encode(HomeResponse{
			Message: "Welcome to new server",
			Status:  true,
		}) // Encoder to create new response
	}
}
