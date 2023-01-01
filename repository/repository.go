package repository

import (
	"context"

	"hajduksanchez.com/go/rest-websockets/models"
)

// Repository for handle user process
type Repository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	InsertPost(ctx context.Context, user *models.Post) error
	Close() error
}

// Implementation for this abstract interface
var implementation Repository

// Function to handle dependency injection for this repository abstraction
func SetRepository(repository Repository) {
	implementation = repository
}

// Function handle by the abstraction
func InsertUser(ctx context.Context, user *models.User) error {
	return implementation.InsertUser(ctx, user)
}

// Function handle by the abstraction
func GetUserById(ctx context.Context, id string) (*models.User, error) {
	return implementation.GetUserById(ctx, id)
}

// Function handle by the abstraction
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.GetUserByEmail(ctx, email)
}

// Function handle by the abstraction
func InsertPost(ctx context.Context, post *models.Post) error {
	return implementation.InsertPost(ctx, post)
}

// Function handle by the abstraction
func Close() error {
	return implementation.Close()
}
