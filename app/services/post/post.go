package post

import (
	PostRepo "wentee/blog/app/repo/post"
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

func (svc *PostService) ListPosts() (postSlice []PostSchena.PostList, err error) {
	posts, err := svc.postRepo.ListPost()

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
