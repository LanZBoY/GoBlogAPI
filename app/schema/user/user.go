package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserCreate struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdate struct {
	Email    *string `json:"email" bson:"Email,omitempty"`
	Username *string `json:"username" bson:"Username,omitempty"`
}

type UserInfo struct {
	Id       primitive.ObjectID `json:"id"`
	Email    string             `json:"email" binding:"email"`
	Username string             `json:"username"`
}
