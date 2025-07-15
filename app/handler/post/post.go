package post

import (
	PostSvc "wentee/blog/app/services/post"

	"github.com/gin-gonic/gin"
)

type PostRouter struct {
	postSvc *PostSvc.PostService
}

func NewPostRouter(postSvc *PostSvc.PostService) *PostRouter {
	return &PostRouter{
		postSvc: postSvc,
	}
}

func (api *PostRouter) CreatePost(c *gin.Context) {

}

func (api *PostRouter) ListPosts(c *gin.Context) {

}

func (api *PostRouter) GetPost(c *gin.Context) {

}

func (api *PostRouter) UpdatePost(c *gin.Context) {

}

func (api *PostRouter) DeletePost(c *gin.Context) {

}

// func (api *PostRouter) ListMyPosts(c *gin.Context) {

// }
