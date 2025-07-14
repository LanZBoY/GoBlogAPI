package mongodb

type UserDocument struct {
	Username string `bson:"Username"`
	Password string `bson:"Password"`
	Salt     string `bson:"Salt"`
}
