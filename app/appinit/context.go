package appinit

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type AppContext struct {
	MongoClient *mongo.Client
}
