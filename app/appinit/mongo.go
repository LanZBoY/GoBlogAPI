package appinit

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoClient(opt *options.ClientOptions) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opt)

	if err != nil {
		panic("Mongo Connection Error!")
	}

	return client
}
