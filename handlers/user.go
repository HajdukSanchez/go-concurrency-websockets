package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
	"hajduksanchez.com/go/rest-websockets/models"
	"hajduksanchez.com/go/rest-websockets/repository"
	"hajduksanchez.com/go/rest-websockets/server"
)

const (
	HASH_COST int = 8
)

// Request for SignUp and Login
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = AuthRequest{}
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

		// Try to hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), HASH_COST)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // Error from server
			return
		}

		// Create and try to insert user
		var user = models.User{
			Email:    request.Email,
			Password: string(hashedPassword), // Convert password hashed into string
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

func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode login data into a struct
		var request = AuthRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		user, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // Error getting user
			return
		}
		if user == nil {
			http.Error(w, "Invalid credential", http.StatusUnauthorized) // User not found
			return
		}

		// Decode hash password and compare with user password credentials
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
			http.Error(w, "Invalid credential", http.StatusUnauthorized) // Password not valid
			return
		}

		// Generate JWT token
		claims := models.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				// Set token expires time
				ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Generate token
		tokenString, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Generate and send Login Response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Token: tokenString,
		})
	}
}

// Get user based on Auth token
func UserHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Token
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		// Validate Token
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Try to get data from Token validating if token is valid
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			// Get user by ID from the Token Payload
			user, err := repository.GetUserById(r.Context(), claims.UserId)
			// Error getting user
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Response user
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
		} else {
			// Error with Token
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
