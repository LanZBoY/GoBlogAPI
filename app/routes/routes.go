package routes

import (
	"wentee/blog/app/di"
	v1 "wentee/blog/app/routes/v1"

	"github.com/gin-gonic/gin"
)

func SetupRouter(container *di.Container) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1.RegistryUserRouter(r.Group("/users"), container)

	return r
}
