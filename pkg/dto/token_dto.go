package dto

import "github.com/golang-jwt/jwt"

type AuthHandler struct {
	secretKey        string
	refreshSecretKey string
	expiresIn        string
	refreshExpiresIn string
}

type JWTClaims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.StandardClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}