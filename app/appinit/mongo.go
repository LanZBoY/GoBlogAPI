package appinit

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConnect wraps mongo.Connect so tests can replace it.
var MongoConnect = mongo.Connect

func GetMongoClient(opt *options.ClientOptions) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := MongoConnect(ctx, opt)

	if err != nil {
		panic("Mongo Connection Error!")
	}

	return client
}
