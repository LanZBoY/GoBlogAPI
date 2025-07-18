package imongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type INewObjectIDFromHex interface {
	ObjectIDFromHex(s string) (primitive.ObjectID, error)
}

type IObjectIdCreator interface {
	INewObjectIDFromHex
}
