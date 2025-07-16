package post

import (
	"net/http"
	"wentee/blog/app/schema/basemodel"
	PostSchema "wentee/blog/app/schema/post"
	PostSvc "wentee/blog/app/services/post"
	"wentee/blog/app/utils/reqcontext"

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

	userInfo, err := reqcontext.GetUserInfo(c)

	if err != nil {
		c.Error(err)
		return
	}

	var postCreate PostSchema.PostCreate
	if err := c.ShouldBindJSON(&postCreate); err != nil {
		c.Error(err)
		return
	}
	if err := api.postSvc.CreatePost(&postCreate, userInfo.Id); err != nil {
		c.Error(err)
		return
	}
}

func (api *PostRouter) ListPosts(c *gin.Context) {
	query := basemodel.NewDefaultQuery()
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(err)
		return
	}

	total, posts, err := api.postSvc.ListPosts(&query)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, basemodel.BaseListResponse{Total: total, Data: posts})
}

func (api *PostRouter) GetPost(c *gin.Context) {
	id := c.Param("id")

	post, err := api.postSvc.GetPostById(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, basemodel.BaseResponse{Data: post})
}

func (api *PostRouter) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var updateData PostSchema.PostUpdate

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.Error(err)
		return
	}
	if err := api.postSvc.UpdatePostById(id, &updateData); err != nil {
		c.Error(err)
		return
	}

}

func (api *PostRouter) DeletePost(c *gin.Context) {
	id := c.Param("id")

	if err := api.postSvc.DeletePostById(id); err != nil {
		c.Error(err)
		return
	}

}

// func (api *PostRouter) ListMyPosts(c *gin.Context) {

// }
