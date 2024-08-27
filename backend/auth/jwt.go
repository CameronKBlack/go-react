package auth

import (
	"fmt"
	"time"

	"go-react/backend/config"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaim struct {
	User string `json:"user"`
	Role string `string:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string, role string) (string, error) {
	claims := CustomClaim{username,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["typ"] = "JWT"
	ss, err := token.SignedString(config.SigningKey)
	if err != nil {
		return "An error has occured when creating JWT", err
	}

	return ss, nil
}

func ValidateJWT(tokenString string) (bool, string) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.SigningKey, nil
	})
	if err != nil {
		return false, err.Error()
	}
	if claims, ok := token.Claims.(*CustomClaim); ok && token.Valid {
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return false, "token has expired"
		}
		return true, "Valid token"
	} else {
		return false, "invalid token claims"
	}
}
