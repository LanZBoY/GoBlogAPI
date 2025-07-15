package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTClaims struct {
	UserInfo JWTUserInfo
	jwt.RegisteredClaims
}

type JWTUserInfo struct {
	Id primitive.ObjectID
}

type LoginInfo struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
