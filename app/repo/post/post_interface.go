package post

import (
	"context"
	PostModel "wentee/blog/app/model/mongodb/post"
	"wentee/blog/app/schema/basemodel"
	PostSchema "wentee/blog/app/schema/post"
	"wentee/blog/app/utils/mongo/imongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IPostCollection interface {
	imongo.IInsertOne
	imongo.IAggregate
	imongo.IUpdateOne
	imongo.IDeleteOne
}

type PostCollectionAdapter struct {
	Collection *mongo.Collection
}

// Aggregate implements IPostCollection.
func (p *PostCollectionAdapter) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (imongo.Cursor, error) {
	return p.Collection.Aggregate(ctx, pipeline, opts...)
}

// DeleteOne implements IPostCollection.
func (p *PostCollectionAdapter) DeleteOne(ctx context.Context, filter any, opts ...*options.DeleteOptions) (imongo.DeleteResult, error) {
	return p.Collection.DeleteOne(ctx, filter, opts...)
}

// InsertOne implements IPostCollection.
func (p *PostCollectionAdapter) InsertOne(ctx context.Context, document any, opts ...*options.InsertOneOptions) (imongo.InsertOneResult, error) {
	return p.Collection.InsertOne(ctx, document, opts...)
}

// UpdateOne implements IPostCollection.
func (p *PostCollectionAdapter) UpdateOne(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (imongo.UpdateResult, error) {
	return p.Collection.UpdateOne(ctx, filter, update, opts...)
}

type IPostRepository interface {
	CreatePost(ctx context.Context, postCreate *PostSchema.PostCreate, createdBy *primitive.ObjectID) (err error)
	ListPosts(ctx context.Context, query *basemodel.BaseQuery) (total int64, posts []PostModel.PostWithCreatorDocument, err error)
	GetPostById(ctx context.Context, id primitive.ObjectID) (post PostModel.PostWithCreatorDocument, err error)
	UpdatePostById(ctx context.Context, id primitive.ObjectID, updateData *PostSchema.PostUpdate) (err error)
	DeletePostById(ctx context.Context, id primitive.ObjectID) (err error)
}
