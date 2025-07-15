package user

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	FieldId       = "_id"
	FieldEmail    = "Email"
	FieldUsername = "Username"
)

type UserDocument struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"Email"`
	Username string             `bson:"Username"`
	Password string             `bson:"Password"`
	Salt     string             `bson:"Salt"`
}
