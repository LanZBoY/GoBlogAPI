package post

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	FieldId        = "_id"
	FieldTitle     = "Title"
	FieldContent   = "Content"
	FieldCreatedAt = "CreatedAt"
	FieldCreatedBy = "CreatedBy"
)

type PostDocument struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"Title"`
	Content   string             `bson:"Content"`
	CreatedAt time.Time          `bson:"CreatedAt"`
	CreatedBy primitive.ObjectID `bson:"CreatedBy"`
}
