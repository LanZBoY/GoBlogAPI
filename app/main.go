package main

import (
	"wentee/blog/app/appinit"
	"wentee/blog/app/config"
	"wentee/blog/app/di"
	_ "wentee/blog/app/docs"
	"wentee/blog/app/routes"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title	Blog API
// @version 0.1.0
// @description 嗨~你好~
// @host localhost:8080
// @basePath /

// @securityDefinitions.apikey BaseAuth
// @in header
// @name Authorizaton
func main() {
	appCtx := setupCtx()
	container := di.InitContainer(appCtx)
	router := routes.SetupRouter(container)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":8080")

}

func setupCtx() *appinit.AppContext {
	newContext := &appinit.AppContext{
		MongoClient: appinit.GetMongoClient(options.Client().ApplyURI(config.MONGO_URI)),
	}

	return newContext
}
