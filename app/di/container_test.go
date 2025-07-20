package di

import (
	"reflect"
	"testing"
	"unsafe"

	"wentee/blog/app/appinit"
	"wentee/blog/app/repo/post"
	"wentee/blog/app/repo/user"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

// getUnexportedField returns the value of an unexported struct field using reflection.
func getUnexportedField(obj interface{}, field string) interface{} {
	v := reflect.ValueOf(obj).Elem().FieldByName(field)
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Interface()
}

func TestInitContainer(t *testing.T) {
	client, err := mongo.NewClient()
	assert.NoError(t, err)
	appCtx := &appinit.AppContext{MongoClient: client}

	c := InitContainer(appCtx)
	if assert.NotNil(t, c) {
		assert.NotNil(t, c.UserRouter)
		assert.NotNil(t, c.AuthRouter)
		assert.NotNil(t, c.PostRouter)
	}

	userSvc := getUnexportedField(c.UserRouter, "userSvc")
	authSvc := getUnexportedField(c.AuthRouter, "authSvc")
	postSvc := getUnexportedField(c.PostRouter, "postSvc")

	assert.NotNil(t, userSvc)
	assert.NotNil(t, authSvc)
	assert.NotNil(t, postSvc)

	// user and auth services should share the same repository
	userRepo := getUnexportedField(userSvc, "userRepo").(*user.UserRepo)
	authRepo := getUnexportedField(authSvc, "userRepo").(*user.UserRepo)
	assert.Equal(t, userRepo, authRepo)

	// verify the repositories are built using the provided mongo client
	ucol := userRepo.UserCollection.(*user.UserCollectionAdapter).Collection
	pRepo := getUnexportedField(postSvc, "postRepo").(*post.PostRepo)
	pcol := getUnexportedField(pRepo, "postColletion").(post.IPostCollection)

	// underlying collection from PostRepo is an adapter as well
	if adapter, ok := pcol.(*post.PostCollectionAdapter); ok {
		assert.Equal(t, client, adapter.Collection.Database().Client())
	}
	assert.Equal(t, client, ucol.Database().Client())
}
