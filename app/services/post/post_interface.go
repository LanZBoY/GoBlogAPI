package post

import (
	"context"
	"wentee/blog/app/schema/basemodel"
	PostSchema "wentee/blog/app/schema/post"
)

type IPostService interface {
	CreatePost(ctx context.Context, postCreate *PostSchema.PostCreate, createdByString string) (err error)
	ListPosts(ctx context.Context, query *basemodel.BaseQuery) (total int64, postSlice []PostSchema.PostList, err error)
	GetPostById(ctx context.Context, id string) (post PostSchema.Post, err error)
	UpdatePostById(ctx context.Context, id string, updateData *PostSchema.PostUpdate) (err error)
	DeletePostById(ctx context.Context, id string) (err error)
}
