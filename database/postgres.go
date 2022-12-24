package database

import (
	"context"
	"database/sql"
	"log"

	"hajduksanchez.com/go/rest-websockets/models"
)

// This repository will be work as a concrete implementation of user repository
type PostgresRepository struct {
	db *sql.DB
}

// Constructor
func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url) // Open SQL connection
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

// Implement User repository
func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	// We use that to create a new SQL statement, passing context to track a debug our flow
	// $ sign tell user which values needs to pass into statement
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)
	return err
}

// Implement User repository
func (repo *PostgresRepository) GetUserById(ctx context.Context, id int64) (*models.User, error) {
	// Query context return rows of data
	rows, _ := repo.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1", id)

	defer func() {
		err := rows.Close() // Close database connection
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user = models.User{}
	for rows.Next() {
		// Try to map values from rows into model
		if err := rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil // Everything ok
		}
	}

	// If there is some error getting data from database
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

// Implement User repository
func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}