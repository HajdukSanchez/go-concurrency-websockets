package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/segmentio/ksuid"
	"hajduksanchez.com/go/rest-websockets/models"
	"hajduksanchez.com/go/rest-websockets/repository"
	"hajduksanchez.com/go/rest-websockets/server"
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpRequest{}
		err := json.NewDecoder(r.Body).Decode(&request) // Try to decode request body into struct
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // Bad request from client
			return
		}
		id, err := ksuid.NewRandom() // UID random generation from external package
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // Error from server
			return
		}
		var user = models.User{
			Email:    request.Email,
			Password: request.Password,
			Id:       id.String(),
		}
		err = repository.InsertUser(r.Context(), &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // Error from server
			return
		}
		// Return correct SignUp
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SignUpResponse{
			Id:    user.Id,
			Email: user.Email,
		})
	}
}
