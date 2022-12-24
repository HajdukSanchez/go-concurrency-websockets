package repository

import (
	"context"

	"hajduksanchez.com/go/rest-websockets/models"
)

// Repository for handle user process
type UserRepository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id int64) (*models.User, error)
	Close() error
}

// Implementation for this abstract interface
var implementation UserRepository

// Function to handle dependency injection for this repository abstraction
func SetRepository(repository UserRepository) {
	implementation = repository
}

// Function handle by the abstraction
func InsertUser(ctx context.Context, user *models.User) error {
	return implementation.InsertUser(ctx, user)
}

// Function handle by the abstraction
func GetUserById(ctx context.Context, id int64) (*models.User, error) {
	return implementation.GetUserById(ctx, id)
}

// Function handle by the abstraction
func Close() error {
	return implementation.Close()
}
