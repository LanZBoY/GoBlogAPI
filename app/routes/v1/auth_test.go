package v1_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"wentee/blog/app/di"
	AuthRouter "wentee/blog/app/handler/auth"
	v1 "wentee/blog/app/routes/v1"
)

func TestRegistryAuthRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	c := &di.Container{AuthRouter: &AuthRouter.AuthRouter{}}

	v1.RegistryAuthRouter(r.Group("/auth"), c)

	expected := gin.RouteInfo{
		Method:  "POST",
		Path:    "/auth/login",
		Handler: handlerName(c.AuthRouter.Login),
	}

	found := false
	for _, rt := range r.Routes() {
		if rt.Method == expected.Method && rt.Path == expected.Path && rt.Handler == expected.Handler {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected route %+v registered", expected)
	}
}
