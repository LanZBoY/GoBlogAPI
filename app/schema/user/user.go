package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserUpdate struct {
	Username *string `json:"username" bson:"Username,omitempty"`
}

type UserInfo struct {
	Id       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
}
