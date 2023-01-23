package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
	"hajduksanchez.com/go/rest-websockets/models"
	"hajduksanchez.com/go/rest-websockets/repository"
	"hajduksanchez.com/go/rest-websockets/server"
	"hajduksanchez.com/go/rest-websockets/utils"
)

// Struct to insert or update post
type UpsertPostRequest struct {
	PostContent string `json:"post_content"`
}

type PostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"post_content"`
}

type PostUpdateResponse struct {
	Message string `json:"message"`
}

// Handler to insert a new post into DB
func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := utils.ValidateAuthorizationToken(s, w, r)
		// Try to get data from Token validating if token is valid
		if err == nil {
			var postRequest = UpsertPostRequest{}
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

// Handler to get a POST from DB by ID
func GetPostById(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r) // Get Path parameters to get ID of post like 'post/:ID'
		post, err := repository.GetPostById(r.Context(), params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post) // Response post data
	}
}

// Handler to update data from specified post
func UpdatePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r) // Get Path parameters to get ID of post like 'post/:ID'
		claims, err := utils.ValidateAuthorizationToken(s, w, r)
		// Try to get data from Token validating if token is valid
		if err == nil {
			var postRequest = UpsertPostRequest{}
			if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Create post Model
			post := models.Post{
				Id:      params["id"],
				Content: postRequest.PostContent,
				UserId:  claims.UserId,
			}

			// Update post
			err = repository.UpdatePost(r.Context(), &post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Send response
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(PostUpdateResponse{
				Message: "Post updated successfully",
			})
		} else {
			// Error with Token
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Handler to delete a specific post
func DeletePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r) // Get Path parameters to get ID of post like 'post/:ID'
		claims, err := utils.ValidateAuthorizationToken(s, w, r)
		// Try to get data from Token validating if token is valid
		if err == nil {

			// Delete post
			err = repository.DeletePost(r.Context(), params["id"], claims.UserId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Send response
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(PostUpdateResponse{
				Message: "Post deleted successfully",
			})
		} else {
			// Error with Token
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
