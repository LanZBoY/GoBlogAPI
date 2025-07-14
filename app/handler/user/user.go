package user

import (
	"net/http"
	"wentee/blog/app/schema/apperror"
	"wentee/blog/app/schema/apperror/errcode"
	"wentee/blog/app/schema/basemodel"
	UserSchema "wentee/blog/app/schema/user"
	UserSvc "wentee/blog/app/service/user"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	userService *UserSvc.UserService
}

func NewUserRouter(userSvc *UserSvc.UserService) *UserRouter {
	return &UserRouter{
		userService: userSvc,
	}
}

func (api *UserRouter) CreateUser(c *gin.Context) {
	var userCreate UserSchema.UserCreate

	if err := c.ShouldBindBodyWithJSON(&userCreate); err != nil {
		c.Error(apperror.New(http.StatusBadRequest, errcode.USER_EXIST, errcode.Message(errcode.USER_EXIST), err))
		return
	}
	if err := api.userService.RegistryUser(&userCreate); err != nil {
		c.Error(err)
		return
	}
}

func (api *UserRouter) GetUser(c *gin.Context) {
	id := c.Param("id")
	userData, err := api.userService.GetUserById(id)

	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, basemodel.BaseResponse{Data: userData})

}

func (api *UserRouter) ListUsers(c *gin.Context) {

	var baseQuery basemodel.BaseQuery

	if err := c.ShouldBindQuery(&baseQuery); err != nil {
		c.Error(err)
		return
	}

	api.userService.ListUsers(&baseQuery)

}

func (api *UserRouter) UpdateUser(c *gin.Context) {
}

func (api *UserRouter) DeleteUser(c *gin.Context) {
}
