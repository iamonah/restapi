package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuth struct {
	secret string
	aud    string
	iss string
}

func NewJWTAuth(secret, aud, iss string) *JWTAuth{
	return &JWTAuth{
		secret: secret,
		aud: aud,
		iss: iss,
	}
}

func (j *JWTAuth) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWTAuth) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return false, fmt.Errorf("invalid signing signature, alg %v", t.Header["alg"])
		}
		return []byte(j.secret), nil
	}, jwt.WithAudience(j.aud))

	return token, err
}
