package web

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserID int `json:"userId"`
	jwt.RegisteredClaims
}
