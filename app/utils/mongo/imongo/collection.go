package imongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Cursor interface {
	All(ctx context.Context, out interface{}) error
	Close(ctx context.Context) error
	Decode(v interface{}) error
	Next(ctx context.Context) bool
}

type SingleResult interface {
	Decode(v any) error
}

type UpdateResult interface {
}

type DeleteResult interface {
}

type InsertOneResult interface {
}

type ICountDocuments interface {
	CountDocuments(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error)
}

type IFind interface {
	Find(ctx context.Context, filter any, opts ...*options.FindOptions) (Cursor, error)
}

type IFindOne interface {
	FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) SingleResult
}

type IInsertOne interface {
	InsertOne(ctx context.Context, document any, opts ...*options.InsertOneOptions) (InsertOneResult, error)
}

type IUpdateOne interface {
	UpdateOne(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (UpdateResult, error)
}

type IDeleteOne interface {
	DeleteOne(ctx context.Context, filter any, opts ...*options.DeleteOptions) (DeleteResult, error)
}

type IAggregate interface {
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (Cursor, error)
}
