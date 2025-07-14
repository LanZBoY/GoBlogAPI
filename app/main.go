package main

import (
	"wentee/blog/app/appinit"
	"wentee/blog/app/config"
	"wentee/blog/app/di"
	"wentee/blog/app/routes"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	appCtx := setupCtx()
	container := di.InitContainer(appCtx)
	router := routes.SetupRouter(container)

	router.Run(":8080")

}

func setupCtx() *appinit.AppContext {
	newContext := &appinit.AppContext{
		MongoClient: appinit.GetMongoClient(options.Client().ApplyURI(config.MONGO_URI)),
	}

	return newContext
}
