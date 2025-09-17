package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	jwt.RegisteredClaims

	UserID int64 `json:"userId"`
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims

	UserID int64 `json:"userId"`
}

func GenerateToken(userID int64) (string, error) {

	now := time.Now()
	tokenExpiryMinutes, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRY_MINUTES"))

	if err != nil {
		return "", err
	}

	claims := TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(tokenExpiryMinutes))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

// GenerateRefreshToken creates a new refresh token for the given user ID.
// The refresh token has a longer expiry time than the access token.
func GenerateRefreshToken(userID int64) (string, error) {
	now := time.Now()
	refreshTokenExpiryDays := 30 // Default to 30 days if env var not set

	if refreshDaysStr := os.Getenv("REFRESH_TOKEN_EXPIRY_DAYS"); refreshDaysStr != "" {
		if days, err := strconv.Atoi(refreshDaysStr); err == nil {
			refreshTokenExpiryDays = days
		}
	}

	claims := RefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * time.Duration(refreshTokenExpiryDays))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	claims := &RefreshTokenClaims{}

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
		return nil, errors.New("refresh token invalid")
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)

	if !ok {
		return nil, errors.New("token valid but couldn't parse claims")
	}

	return claims, nil
}

func FindToken(request *http.Request) (string, error) {
	// Get token from authorization header.
	bearer := request.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:], nil
	}
	return "", fmt.Errorf("invalid token: %s", bearer)
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
		return nil, errors.New("token invalid")
	}

	claims, ok := token.Claims.(*TokenClaims)

	if !ok {
		return nil, errors.New("token valid but couldn't parse claims")
	}

	return claims, nil
}
