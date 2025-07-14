package v1

import (
	"wentee/blog/app/di"

	"github.com/gin-gonic/gin"
)

func RegistryUserRouter(rg *gin.RouterGroup, container *di.Container) {
	rg.POST("", container.UserRouter.CreateUser)
	rg.GET("", container.UserRouter.ListUsers)

	id_rg := rg.Group("/:id")
	{
		id_rg.GET("", container.UserRouter.GetUser)
		id_rg.PATCH("", container.UserRouter.UpdateUser)
		id_rg.DELETE("", container.UserRouter.DeleteUser)
	}
}
