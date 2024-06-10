package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type TokenClaims struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateToken(userId int64) (string, error) {

	now := time.Now()
	tokenExpiryMinutes, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRY_MINUTES"))

	if err != nil {
		return "", err
	}

	claims := TokenClaims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(tokenExpiryMinutes))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func FindToken(request *http.Request) (string, error) {
	// Get token from authorization header.
	bearer := request.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:], nil
	}
	return "", fmt.Errorf("invalid token: " + bearer)
}

func ValidateToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token invalid")
	}

	claims, ok := token.Claims.(*TokenClaims)

	if !ok {
		return nil, fmt.Errorf("token valid but couldn't parse claims")
	}

	return claims, nil
}
