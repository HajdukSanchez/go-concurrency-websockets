package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"hajduksanchez.com/go/rest-websockets/handlers"
	"hajduksanchez.com/go/rest-websockets/server"
)

func main() {
	err := godotenv.Load(".env") // Load the environments

	if err != nil {
		log.Fatal("Error loading environments")
	}

	// Get all environments
	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATA_BASE_URL := os.Getenv("DATA_BASE_URL")

	// Create the new server
	server, err := server.NewServer(context.Background(), &server.Config{
		JWTSecret: JWT_SECRET,
		Port:      PORT,
		DBUrl:     DATA_BASE_URL,
	})

	if err != nil {
		log.Fatal("Error creating server")
	}

	server.Start(BindRoutes) // Start the server
}

// Function to handle routes and start server
func BindRoutes(server server.Server, router *mux.Router) {
	// Define endpoints and methods for endpoints
	router.HandleFunc("/", handlers.HomeHandler(server)).Methods(http.MethodGet)
	router.HandleFunc("/sign_up", handlers.SignUpHandler(server)).Methods(http.MethodPost)
	router.HandleFunc("/login", handlers.LoginHandler(server)).Methods(http.MethodPost)
}
