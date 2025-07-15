package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserInfo JWTUserInfo
	jwt.RegisteredClaims
}

type JWTUserInfo struct {
	Id string
}

type LoginInfo struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
