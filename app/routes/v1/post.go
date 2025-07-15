package v1

import (
	"wentee/blog/app/di"

	"github.com/gin-gonic/gin"
)

func RegistryPostRouter(rg *gin.RouterGroup, container *di.Container) {
	rg.POST("", container.PostRouter.CreatePost)
	rg.GET("", container.PostRouter.ListPosts)

	id_group := rg.Group("/:id")
	{
		id_group.GET("", container.PostRouter.GetPost)
		id_group.PATCH("", container.PostRouter.UpdatePost)
		id_group.DELETE("", container.PostRouter.DeletePost)
	}
}
