package post

import "go.mongodb.org/mongo-driver/bson/primitive"

type PostCreate struct {
	Title   string  `json:"title" bind:"rquired"`
	Content *string `json:"content" bind:"omitempty"`
}

type Creator struct {
	Id       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
}

type PostList struct {
	Id      primitive.ObjectID `json:"id"`
	Title   string             `json:"title"`
	Creator Creator            `json:"creator"`
}
