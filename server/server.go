package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"hajduksanchez.com/go/rest-websockets/database"
	"hajduksanchez.com/go/rest-websockets/repository"
)

// Configuration to connect our server
type Config struct {
	Port      string // Port to connect to
	JWTSecret string // JWTSecret to connect to
	DBUrl     string // DB URL to connect to
}

type Server interface {
	Config() *Config // Server configuration
}

// / Broker is going to handle servers
type Broker struct {
	config *Config     // Properties to configure
	router *mux.Router // Router to define API routes
}

// Broker is no a server implementation
func (b *Broker) Config() *Config {
	return b.config
}

// Create a new server
// [ctx] allow us to identify where is the problem (for example if we work in routines)
func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is not specified")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("JWTSecret is not specified")
	}
	if config.DBUrl == "" {
		return nil, errors.New("DBUrl is not specified")
	}

	// If there is no error we create and return a new broker (server)
	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
	}
	return broker, nil
}

// Start a new server instance
func (b *Broker) Start(binder func(server Server, router *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)

	// Start DB connection
	repo, err := database.NewPostgresRepository(b.config.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	repository.SetRepository(repo)

	// Start server
	log.Println("Starting server on port", b.Config().Port)
	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatal("ListenAndServe: ", err) // If something goes wrong on HTTP initialization
	}
}
