package v1

import (
	"wentee/blog/app/di"
	"wentee/blog/app/middleware"

	"github.com/gin-gonic/gin"
)

func RegistryUserRouter(rg *gin.RouterGroup, container *di.Container) {
	rg.POST("", middleware.ErrorHandlerWrapper(container.UserRouter.CreateUser))
	rg.GET("", middleware.ErrorHandlerWrapper(container.UserRouter.ListUsers))

	id_rg := rg.Group("/:id")
	{
		id_rg.GET("", middleware.ErrorHandlerWrapper(container.UserRouter.GetUser))
		id_rg.PATCH("", middleware.ErrorHandlerWrapper(container.UserRouter.UpdateUser))
		id_rg.DELETE("", middleware.ErrorHandlerWrapper(container.UserRouter.DeleteUser))
	}
}
