package post

import PostRepo "wentee/blog/app/repo/post"

type PostService struct {
	postRepo *PostRepo.PostRepo
}

func NewPostService(postRepo *PostRepo.PostRepo) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}
