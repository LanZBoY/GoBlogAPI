package post

import (
	"context"
	PostModel "wentee/blog/app/model/mongodb/post"
	"wentee/blog/app/schema/basemodel"
	PostSchema "wentee/blog/app/schema/post"
	"wentee/blog/app/utils/mongo/icollection"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPostCollection interface {
	icollection.IInsertOne
	icollection.IAggregate
	icollection.IUpdateOne
	icollection.IDeleteOne
}

type IPostRepository interface {
	CreatePost(ctx context.Context, postCreate *PostSchema.PostCreate, createdBy *primitive.ObjectID) (err error)
	ListPosts(ctx context.Context, query *basemodel.BaseQuery) (total int64, posts []PostModel.PostWithCreatorDocument, err error)
	GetPostById(ctx context.Context, id primitive.ObjectID) (post PostModel.PostWithCreatorDocument, err error)
	UpdatePostById(ctx context.Context, id primitive.ObjectID, updateData *PostSchema.PostUpdate) (err error)
	DeletePostById(ctx context.Context, id primitive.ObjectID) (err error)
}
