package user

import (
	"net/http"
	"wentee/blog/app/schema/basemodel"
	UserSchema "wentee/blog/app/schema/user"
	UserSvc "wentee/blog/app/services/user"
	"wentee/blog/app/utils/reqcontext"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	userSvc *UserSvc.UserService
}

func NewUserRouter(userSvc *UserSvc.UserService) *UserRouter {
	return &UserRouter{
		userSvc: userSvc,
	}
}

func (api *UserRouter) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	var userCreate UserSchema.UserCreate

	if err := c.ShouldBindJSON(&userCreate); err != nil {
		c.Error(err)
		return
	}

	if err := api.userSvc.RegistryUser(ctx, &userCreate); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusCreated)
}

func (api *UserRouter) GetUser(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	userData, err := api.userSvc.GetUserById(ctx, id)

	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, basemodel.BaseResponse{Data: userData})

}

func (api *UserRouter) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()
	baseQuery := basemodel.NewDefaultQuery()

	if err := c.ShouldBindQuery(&baseQuery); err != nil {
		c.Error(err)
		return
	}
	total, err := api.userSvc.CountUsers(ctx)

	if err != nil {
		c.Error(err)
		return
	}

	users, err := api.userSvc.ListUsers(ctx, &baseQuery)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, basemodel.BaseListResponse{Total: total, Data: users})
}

func (api *UserRouter) GetMe(c *gin.Context) {
	ctx := c.Request.Context()
	userInfo, err := reqcontext.GetUserInfo(c)
	if err != nil {
		c.Error(err)
		return
	}

	userData, err := api.userSvc.GetUserById(ctx, userInfo.Id)

	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, basemodel.BaseResponse{Data: userData})
}

func (api *UserRouter) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	var userUpdate UserSchema.UserUpdate

	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.Error(err)
		return
	}
	if err := api.userSvc.UpdateUserById(ctx, id, userUpdate); err != nil {
		c.Error(err)
		return
	}

}

func (api *UserRouter) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	err := api.userSvc.DeleteUserById(ctx, c.Param("id"))

	if err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
