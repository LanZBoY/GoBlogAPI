package post

import (
	PostRepo "wentee/blog/app/repo/post"
	"wentee/blog/app/schema/basemodel"
	PostSchema "wentee/blog/app/schema/post"
	PostSchena "wentee/blog/app/schema/post"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService struct {
	postRepo *PostRepo.PostRepo
}

func NewPostService(postRepo *PostRepo.PostRepo) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

func (svc *PostService) CreatePost(postCreate *PostSchema.PostCreate, createdByString string) (err error) {
	createdBy, err := primitive.ObjectIDFromHex(createdByString)
	if err != nil {
		return
	}
	err = svc.postRepo.CreatePost(postCreate, &createdBy)
	return
}

func (svc *PostService) ListPosts(query *basemodel.BaseQuery) (total int64, postSlice []PostSchena.PostList, err error) {
	total, posts, err := svc.postRepo.ListPosts(query)

	postSlice = make([]PostSchena.PostList, len(posts))

	for index, post := range posts {
		postSlice[index] = PostSchena.PostList{
			Id:      post.Id,
			Title:   post.Title,
			Creator: PostSchena.Creator{Id: post.Creator.Id, Username: post.Creator.Username},
		}
	}

	return
}

func (svc *PostService) GetPostById(id string) (post PostSchema.Post, err error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return
	}
	postDoc, err := svc.postRepo.GetPostById(oid)
	if err != nil {
		return
	}
	post = PostSchena.Post{
		PostList: PostSchena.PostList{
			Id:      postDoc.Id,
			Title:   postDoc.Title,
			Creator: PostSchena.Creator{Id: postDoc.Creator.Id, Username: postDoc.Creator.Username},
		},
		Content:   *postDoc.Content,
		CreatedAt: postDoc.CreatedAt,
	}
	return
}

func (svc *PostService) UpdatePostById(id string, updateData *PostSchema.PostUpdate) (err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	err = svc.postRepo.UpdatePostById(oid, updateData)
	return
}

func (svc *PostService) DeletePostById(id string) (err error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return
	}

	err = svc.postRepo.DeletePostById(oid)

	return
}
