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
	postSvc PostSvc.IPostService
}

func NewPostRouter(postSvc *PostSvc.PostService) *PostRouter {
	return &PostRouter{
		postSvc: postSvc,
	}
}

// @summary 建立貼文
// @description 建立貼文
// @security BasicAuth
// @tags Post
// @accept application/json
// @produce application/json
// @param PostCreateData body PostSchema.PostCreate true "貼文資訊"
// @Success	201
// @router /posts [POST]
func (api *PostRouter) CreatePost(c *gin.Context) {
	ctx := c.Request.Context()
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
	if err := api.postSvc.CreatePost(ctx, &postCreate, userInfo.Id); err != nil {
		c.Error(err)
		return
	}
}

// @summary 貼文列表
// @description 取得貼文列表
// @security BasicAuth
// @tags Post
// @accept application/json
// @produce application/json
// @Success	200 {object} basemodel.BaseListResponse{total=int, data=[]PostSchema.PostList}
// @router /posts [GET]
func (api *PostRouter) ListPosts(c *gin.Context) {
	ctx := c.Request.Context()
	query := basemodel.NewDefaultQuery()
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(err)
		return
	}

	total, posts, err := api.postSvc.ListPosts(ctx, &query)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, basemodel.BaseListResponse{Total: total, Data: posts})
}

// @summary 取得貼文內容
// @description 透過ID取得貼文內容
// @security BasicAuth
// @tags Post
// @accept application/json
// @produce application/json
// @param id path string true "貼文ID"
// @Success	200 {object} basemodel.BaseResponse{data=PostSchema.Post}
// @router /posts/{id} [GET]
func (api *PostRouter) GetPost(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	post, err := api.postSvc.GetPostById(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, basemodel.BaseResponse{Data: post})
}

// @summary 更新貼文內容
// @description 透過ID更新貼文內容
// @security BasicAuth
// @tags Post
// @accept application/json
// @produce application/json
// @param id path string true "貼文ID"
// @Param PostUpdate body PostSchema.PostUpdate true "更新貼文欄位"
// @Success	200
// @router /posts/{id} [PATCH]
func (api *PostRouter) UpdatePost(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	var updateData PostSchema.PostUpdate

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.Error(err)
		return
	}
	if err := api.postSvc.UpdatePostById(ctx, id, &updateData); err != nil {
		c.Error(err)
		return
	}

}

// @summary 刪除貼文
// @description 透過ID刪除貼文
// @security BasicAuth
// @tags Post
// @accept application/json
// @produce application/json
// @param id path string true "貼文ID"
// @Success	204
// @router /posts/{id} [DELETE]
func (api *PostRouter) DeletePost(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if err := api.postSvc.DeletePostById(ctx, id); err != nil {
		c.Error(err)
		return
	}

}
