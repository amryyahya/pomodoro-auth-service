package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateTokens generates an access token and a refresh token for a given email
func GenerateTokens(email string, accessSecret string, refreshSecret string) (string, string, error) {
	// Generate Access Token
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["email"] = email
	accessClaims["exp"] = time.Now().Add(time.Minute * 15).Unix() // Access token expires in 15 minutes

	accessTokenString, err := accessToken.SignedString([]byte(accessSecret))
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["email"] = email
	refreshClaims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix() // Refresh token expires in a month

	refreshTokenString, err := refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// RefreshAccessToken refreshes an access token using a refresh token
func RefreshAccessToken(refreshTokenString string,accessSecret string, refreshSecret string) (string, error) {
	// Validate the refresh token
	claims, err := ValidateToken(refreshTokenString, []byte(refreshSecret))
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Extract email
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return "", fmt.Errorf("invalid email in token")
	}

	// Generate a new access token
	newAccessToken := jwt.New(jwt.SigningMethodHS256)
	newAccessClaims := newAccessToken.Claims.(jwt.MapClaims)
	newAccessClaims["email"] = email
	newAccessClaims["exp"] = time.Now().Add(time.Minute * 15).Unix() // New access token expires in 15 minutes

	return newAccessToken.SignedString([]byte(accessSecret))
}


func ValidateToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}