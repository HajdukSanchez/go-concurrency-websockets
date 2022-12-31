package models

import "github.com/golang-jwt/jwt"

type AppClaims struct {
	UserId             string `json:"userId"`
	jwt.StandardClaims        // AppClaims contains all properties of the package
}
