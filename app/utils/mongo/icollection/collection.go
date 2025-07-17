package icollection

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ICountDocuments interface {
	CountDocuments(context.Context, any, ...*options.CountOptions) (int64, error)
}

type IFind interface {
	Find(context.Context, any, ...*options.FindOptions) (*mongo.Cursor, error)
}

type IFindOne interface {
	FindOne(context.Context, any, ...*options.FindOneOptions) *mongo.SingleResult
}

type IInsertOne interface {
	InsertOne(context.Context, any, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
}

type IUpdateOne interface {
	UpdateOne(context.Context, any, any, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type IDeleteOne interface {
	DeleteOne(context.Context, any, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

type IAggregate interface {
	Aggregate(context.Context, interface{}, ...*options.AggregateOptions) (*mongo.Cursor, error)
}
