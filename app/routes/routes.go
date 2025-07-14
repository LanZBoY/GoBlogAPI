package routes

import (
	"wentee/blog/app/di"
	"wentee/blog/app/middleware"
	v1 "wentee/blog/app/routes/v1"

	"github.com/gin-gonic/gin"
)

func SetupRouter(container *di.Container) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandler())

	v1.RegistryUserRouter(r.Group("/users"), container)

	return r
}
