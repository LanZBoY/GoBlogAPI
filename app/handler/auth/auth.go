package auth

import (
	"net/http"
	AuthSchema "wentee/blog/app/schema/auth"
	"wentee/blog/app/schema/basemodel"
	AuthSvc "wentee/blog/app/services/auth"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	authSvc *AuthSvc.AuthService
}

func NewAuthRouter(authSvc *AuthSvc.AuthService) *AuthRouter {

	return &AuthRouter{
		authSvc: authSvc,
	}
}

func (api *AuthRouter) Login(c *gin.Context) {
	var loginInfo AuthSchema.LoginInfo

	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.Error(err)
		return
	}
	tokenString, err := api.authSvc.TryLogin(&loginInfo)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, basemodel.BaseResponse{Data: AuthSchema.TokenResponse{Token: tokenString}})
}
