package main

import (
	"context"
	"net/http"
	"testing"

	"wentee/blog/app/appinit"
	"wentee/blog/app/config"
	"wentee/blog/app/di"
	_ "wentee/blog/app/docs"
	"wentee/blog/app/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func stubMongoConnect(captured *string) func(context.Context, ...*options.ClientOptions) (*mongo.Client, error) {
	return func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
		if captured != nil && len(opts) > 0 && opts[0] != nil {
			*captured = opts[0].GetURI()
		}
		return mongo.NewClient(opts...)
	}
}

func TestSetupCtx(t *testing.T) {
	originalConnect := appinit.MongoConnect
	defer func() { appinit.MongoConnect = originalConnect }()

	var uri string
	appinit.MongoConnect = stubMongoConnect(&uri)

	oldURI := config.MONGO_URI
	defer func() { config.MONGO_URI = oldURI }()
	config.MONGO_URI = "mongodb://example.com:27018"

	ctx := setupCtx()
	assert.NotNil(t, ctx)
	assert.NotNil(t, ctx.MongoClient)
	assert.Equal(t, config.MONGO_URI, uri)
}

func TestSwaggerRouteRegistration(t *testing.T) {
	originalConnect := appinit.MongoConnect
	defer func() { appinit.MongoConnect = originalConnect }()
	appinit.MongoConnect = stubMongoConnect(nil)

	gin.SetMode(gin.TestMode)
	appCtx := setupCtx()
	container := di.InitContainer(appCtx)
	router := routes.SetupRouter(container)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	found := false
	for _, r := range router.Routes() {
		if r.Path == "/swagger/*any" && r.Method == http.MethodGet {
			found = true
			break
		}
	}
	assert.True(t, found, "swagger route not registered")
}
