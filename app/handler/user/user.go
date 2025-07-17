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

// @summary 建立使用者
// @description 建立一般使用者
// @security BasicAuth
// @tags User
// @accept application/json
// @produce application/json
// @param UserCreateData body UserSchema.UserCreate true "使用者建立資訊"
// @Success	201
// @router /users [POST]
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

// @summary 取得使用者資訊
// @description 取得使用者
// @security BasicAuth
// @tags User
// @accept application/json
// @produce application/json
// @param id path string true "使用者識別ID"
// @Success	200 {object} basemodel.BaseResponse{data=UserSchema.UserInfo}
// @router /users/{id} [GET]
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

// @summary 使用者列表
// @description 取得使用者列表
// @security BasicAuth
// @tags User
// @accept application/json
// @produce application/json
// @Success	200 {object} basemodel.BaseListResponse{total=int, data=[]UserSchema.UserInfo}
// @router /users [GET]
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

// @summary 取得我的資訊
// @description 取得我的基本資訊
// @security BasicAuth
// @tags User
// @accept application/json
// @produce application/json
// @Success	200 {object} basemodel.BaseResponse{data=UserSchema.UserInfo}
// @router /users/me [GET]
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

// @summary 刪除使用者
// @description 刪除使用者
// @security BasicAuth
// @tags User
// @accept application/json
// @produce application/json
// @param id path string true "使用者識別ID"
// @Success	204
// @router /users/{id} [DELETE]
func (api *UserRouter) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	err := api.userSvc.DeleteUserById(ctx, c.Param("id"))

	if err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// @summary 更新使用者資訊
// @description 更新使用者
// @security BasicAuth
// @tags User
// @accept application/json
// @produce application/json
// @param id path string true "使用者識別ID"
// @Param UserUpdate body UserSchema.UserUpdate true "更新使用者資訊欄位"
// @Success	200
// @router /users/{id} [PATCH]
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
