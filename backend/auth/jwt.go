package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claim struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var JWTKEY = []byte("My_JWT_Key")

func GenerateJWTToken(username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claim := &Claim{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(JWTKEY)
}

func ValidateJWTToken(tokenString string) (*Claim, error) {
	claims := &Claim{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWTKEY, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
