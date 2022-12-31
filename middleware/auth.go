package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"hajduksanchez.com/go/rest-websockets/models"
	"hajduksanchez.com/go/rest-websockets/server"
	"hajduksanchez.com/go/rest-websockets/utils"
)

var (
	// List of routes that not needs authentication
	NO_AUTH_NEEDED = []string{
		utils.Login,
		utils.Register,
	}
)

// Function to know if it is important to check token or not based on route
func shouldCheckToken(route string) bool {
	for _, route := range NO_AUTH_NEEDED {
		if strings.Contains(route, route) {
			return false
		}
	}
	return true
}

// Middleware returns next function handler specified
// That is because middleware works as a previous handler function that surround a handler function
// If everything ok, handler function passed will be executed
func CheckAuthMiddleware(s server.Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.URL.Path)
			// Validate if route needs to be authenticated
			if !shouldCheckToken(r.URL.Path) {
				next.ServeHTTP(w, r) // Continue with handler function of the specific path
				return
			}

			// Get Token and validate if user has permission based on this specific token
			tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
			_, err := jwt.ParseWithClaims(tokenString, models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r) // Continue with handler function
		})
	}
}
