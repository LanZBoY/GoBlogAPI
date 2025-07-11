package main

import (
	"fmt"
	"wentee/blog/app/appinit"
	"wentee/blog/app/config"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	router := appinit.NewRouter()
	appCtx := setupCtx()
	fmt.Printf("%v\n", appCtx)
	router.Run(":8080")

}

func setupCtx() *appinit.AppContext {
	newContext := &appinit.AppContext{
		MongoClient: appinit.GetMongoClient(options.Client().ApplyURI(config.MONGO_URI)),
	}

	return newContext
}
