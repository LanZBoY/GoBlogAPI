package appinit

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoConnect is used to establish the connection to MongoDB. It is defined as
// a variable so it can be replaced in tests.
var mongoConnect = mongo.Connect

func GetMongoClient(opt *options.ClientOptions) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongoConnect(ctx, opt)

	if err != nil {
		panic("Mongo Connection Error!")
	}

	return client
}
