package v1

import (
	"wentee/blog/app/di"

	"github.com/gin-gonic/gin"
)

func RegistryAuthRouter(rg *gin.RouterGroup, container *di.Container) {
	rg.POST("/login", container.AuthRouter.Login)
}
