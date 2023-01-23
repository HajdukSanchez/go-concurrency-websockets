package utils

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"hajduksanchez.com/go/rest-websockets/models"
	"hajduksanchez.com/go/rest-websockets/server"
)

// Validate token info and return claims from token or error
func ValidateAuthorizationToken(s server.Server, w http.ResponseWriter, r *http.Request) (*models.AppClaims, error) {
	// Get Token
	tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
	// Validate Token
	token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Config().JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	// Try to get data from Token validating if token is valid
	if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
