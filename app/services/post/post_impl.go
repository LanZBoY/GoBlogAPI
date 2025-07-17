package post

import (
	"context"
	PostRepo "wentee/blog/app/repo/post"
	"wentee/blog/app/schema/basemodel"
	PostSchema "wentee/blog/app/schema/post"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService struct {
	postRepo PostRepo.IPostRepository
}

func NewPostService(postRepo *PostRepo.PostRepo) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

func (svc *PostService) CreatePost(ctx context.Context, postCreate *PostSchema.PostCreate, createdByString string) (err error) {
	createdBy, err := primitive.ObjectIDFromHex(createdByString)
	if err != nil {
		return
	}
	err = svc.postRepo.CreatePost(ctx, postCreate, &createdBy)
	return
}

func (svc *PostService) ListPosts(ctx context.Context, query *basemodel.BaseQuery) (total int64, postSlice []PostSchema.PostList, err error) {
	total, posts, err := svc.postRepo.ListPosts(ctx, query)

	postSlice = make([]PostSchema.PostList, len(posts))

	for index, post := range posts {
		postSlice[index] = PostSchema.PostList{
			Id:      post.Id,
			Title:   post.Title,
			Creator: PostSchema.Creator{Id: post.Creator.Id, Username: post.Creator.Username},
		}
	}

	return
}

func (svc *PostService) GetPostById(ctx context.Context, id string) (post PostSchema.Post, err error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return
	}
	postDoc, err := svc.postRepo.GetPostById(ctx, oid)
	if err != nil {
		return
	}
	post = PostSchema.Post{
		PostList: PostSchema.PostList{
			Id:      postDoc.Id,
			Title:   postDoc.Title,
			Creator: PostSchema.Creator{Id: postDoc.Creator.Id, Username: postDoc.Creator.Username},
		},
		Content:   *postDoc.Content,
		CreatedAt: postDoc.CreatedAt,
	}
	return
}

func (svc *PostService) UpdatePostById(ctx context.Context, id string, updateData *PostSchema.PostUpdate) (err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	err = svc.postRepo.UpdatePostById(ctx, oid, updateData)
	return
}

func (svc *PostService) DeletePostById(ctx context.Context, id string) (err error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return
	}

	err = svc.postRepo.DeletePostById(ctx, oid)

	return
}
