package post

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostCreate struct {
	Title   string  `json:"title" bind:"rquired"`
	Content *string `json:"content" bind:"omitempty"`
}

type PostUpdate struct {
	Title   *string `json:"title" bson:"Title,omitempty"`
	Content *string `json:"content" bson:"Content,omitempty"`
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

type Post struct {
	PostList
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
