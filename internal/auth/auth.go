package auth

import "github.com/golang-jwt/jwt"

type Auth interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(token string) (bool, error)
}
