package routes_test

import (
	"github.com/gin-gonic/gin"
	"testing"
	"wentee/blog/app/di"
	AuthRouter "wentee/blog/app/handler/auth"
	PostRouter "wentee/blog/app/handler/post"
	UserRouter "wentee/blog/app/handler/user"
	"wentee/blog/app/routes"
)

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := &di.Container{
		UserRouter: &UserRouter.UserRouter{},
		AuthRouter: &AuthRouter.AuthRouter{},
		PostRouter: &PostRouter.PostRouter{},
	}

	r := routes.SetupRouter(c)

	if got := len(r.Handlers); got != 3 {
		t.Fatalf("expected 3 global middlewares, got %d", got)
	}

	expected := map[string]bool{
		"/auth/login": false,
		"/users":      false,
		"/users/me":   false,
		"/users/:id":  false,
		"/posts":      false,
		"/posts/:id":  false,
	}

	for _, info := range r.Routes() {
		if _, ok := expected[info.Path]; ok {
			expected[info.Path] = true
		}
	}

	for p, ok := range expected {
		if !ok {
			t.Errorf("route %s not registered", p)
		}
	}
}
