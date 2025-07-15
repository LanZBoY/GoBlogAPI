package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserDocument struct {
	Id       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"Username"`
	Password string             `bson:"Password"`
	Salt     string             `bson:"Salt"`
}
