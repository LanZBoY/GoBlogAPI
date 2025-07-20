package v1_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"wentee/blog/app/di"
	UserRouter "wentee/blog/app/handler/user"
	v1 "wentee/blog/app/routes/v1"
)

func TestRegistryUserRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	c := &di.Container{UserRouter: &UserRouter.UserRouter{}}

	v1.RegistryUserRouter(r.Group("/users"), c)

	checks := []gin.RouteInfo{
		{Method: "POST", Path: "/users", Handler: handlerName(c.UserRouter.CreateUser)},
		{Method: "GET", Path: "/users", Handler: handlerName(c.UserRouter.ListUsers)},
		{Method: "GET", Path: "/users/me", Handler: handlerName(c.UserRouter.GetMe)},
		{Method: "GET", Path: "/users/:id", Handler: handlerName(c.UserRouter.GetUser)},
		{Method: "PATCH", Path: "/users/:id", Handler: handlerName(c.UserRouter.UpdateUser)},
		{Method: "DELETE", Path: "/users/:id", Handler: handlerName(c.UserRouter.DeleteUser)},
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
