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

func (api *UserRouter) CreateUser(c *gin.Context) error {
	var userCreate UserSchema.UserCreate

	if err := c.ShouldBindBodyWithJSON(&userCreate); err != nil {
		return apperror.New(http.StatusBadRequest, errcode.USER_EXIST, errcode.Message(errcode.USER_EXIST), err)
	}
	if err := api.userService.RegistryUser(&userCreate); err != nil {
		return err
	}
	return nil
}

func (api *UserRouter) GetUser(c *gin.Context) error {
	id := c.Param("id")
	userData, err := api.userService.GetUserById(id)

	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, basemodel.BaseResponse{Data: userData})
	return nil
}

func (api *UserRouter) ListUsers(c *gin.Context) error {

	var baseQuery basemodel.BaseQuery

	if err := c.ShouldBindQuery(&baseQuery); err != nil {
		return err
	}

	api.userService.ListUsers(&baseQuery)

	return nil
}

func (api *UserRouter) UpdateUser(c *gin.Context) error {

	return nil
}

func (api *UserRouter) DeleteUser(c *gin.Context) error {

	return nil
}
