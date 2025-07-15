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
	authMiddleware := middleware.AuthMiddleware{}

	v1.RegistryAuthRouter(r.Group("/auth"), container)
	v1.RegistryUserRouter(r.Group("/users", authMiddleware.RequiredAuth()), container)
	v1.RegistryPostRouter(r.Group("/posts", authMiddleware.RequiredAuth()), container)

	return r
}
