package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"hajduksanchez.com/go/rest-websockets/models"
	"hajduksanchez.com/go/rest-websockets/repository"
	"hajduksanchez.com/go/rest-websockets/server"
)

type InsertPostRequest struct {
	PostContent string `json:"post_content"`
}

type PostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"post_content"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
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
			var postRequest = InsertPostRequest{}
			if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Generate new ID
			id, err := ksuid.NewRandom()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Create post Model
			post := models.Post{
				Id:      id.String(),
				Content: postRequest.PostContent,
				UserId:  claims.UserId,
			}
			// Insert post
			err = repository.InsertPost(r.Context(), &post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Send response
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(PostResponse{
				Id:          post.Id,
				PostContent: post.Content,
			})
		} else {
			// Error with Token
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
