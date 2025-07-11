package appinit

import (
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return router
}

func RegistryUserRouter(rg *gin.RouterGroup, appCtx *AppContext, middleware *gin.HandlerFunc) {
}
