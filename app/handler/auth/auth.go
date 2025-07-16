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

// @Summary 登入API
// @Tags	Auth
// @Accept	application/json
// @Produce	application/json
// @Param	loginInfo	body	AuthSchema.LoginInfo	true	"登入資訊"
// @Success	200	{object}	basemodel.BaseResponse{data=AuthSchema.TokenResponse}
// @Router	/auth/login	[post]
func (api *AuthRouter) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var loginInfo AuthSchema.LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.Error(err)
		return
	}
	tokenString, err := api.authSvc.TryLogin(ctx, &loginInfo)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, basemodel.BaseResponse{Data: AuthSchema.TokenResponse{Token: tokenString}})
}
