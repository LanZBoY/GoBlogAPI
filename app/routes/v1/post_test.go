package v1_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"wentee/blog/app/di"
	PostRouter "wentee/blog/app/handler/post"
	v1 "wentee/blog/app/routes/v1"
)

func TestRegistryPostRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	c := &di.Container{PostRouter: &PostRouter.PostRouter{}}

	v1.RegistryPostRouter(r.Group("/posts"), c)

	checks := []gin.RouteInfo{
		{Method: "POST", Path: "/posts", Handler: handlerName(c.PostRouter.CreatePost)},
		{Method: "GET", Path: "/posts", Handler: handlerName(c.PostRouter.ListPosts)},
		{Method: "GET", Path: "/posts/:id", Handler: handlerName(c.PostRouter.GetPost)},
		{Method: "PATCH", Path: "/posts/:id", Handler: handlerName(c.PostRouter.UpdatePost)},
		{Method: "DELETE", Path: "/posts/:id", Handler: handlerName(c.PostRouter.DeletePost)},
	}

	for _, exp := range checks {
		found := false
		for _, rt := range r.Routes() {
			if rt.Method == exp.Method && rt.Path == exp.Path && rt.Handler == exp.Handler {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected route %+v registered", exp)
		}
	}
}
