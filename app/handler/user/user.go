package user

import (
	"net/http"
	"wentee/blog/app/schema/apperror"
	"wentee/blog/app/schema/apperror/errcode"
	"wentee/blog/app/schema/auth"
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
	var userCreate UserSchema.UserCreate

	if err := c.ShouldBindJSON(&userCreate); err != nil {
		c.Error(err)
		return
	}

	if err := api.userSvc.RegistryUser(&userCreate); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusCreated)
}

func (api *UserRouter) GetUser(c *gin.Context) {
	id := c.Param("id")
	userData, err := api.userSvc.GetUserById(id)

	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, basemodel.BaseResponse{Data: userData})

}

func (api *UserRouter) ListUsers(c *gin.Context) {

	baseQuery := basemodel.NewDefaultQuery()

	if err := c.ShouldBindQuery(&baseQuery); err != nil {
		c.Error(err)
		return
	}
	total, err := api.userSvc.CountUsers()

	if err != nil {
		c.Error(err)
		return
	}

	users, err := api.userSvc.ListUsers(&baseQuery)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, basemodel.BaseListResponse{Total: total, Data: users})
}

func (api *UserRouter) GetMe(c *gin.Context) {
	userInfoAny, ok := c.Get(reqcontext.USER_INFO)
	if !ok {
		c.Error(apperror.New(http.StatusBadRequest, errcode.USER_NOT_FOUND, nil))
		return
	}
	userInfo, ok := userInfoAny.(auth.JWTUserInfo)
	if !ok {
		c.Error(apperror.New(http.StatusInternalServerError, errcode.TPYE_ASSERTION_ERROR, nil))
		return
	}

	userData, err := api.userSvc.GetUserById(userInfo.Id)

	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, basemodel.BaseResponse{Data: userData})
}

func (api *UserRouter) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var userUpdate UserSchema.UserUpdate

	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.Error(err)
		return
	}
	if err := api.userSvc.UpdateUserById(id, userUpdate); err != nil {
		c.Error(err)
		return
	}

}

func (api *UserRouter) DeleteUser(c *gin.Context) {
	err := api.userSvc.DeleteUserById(c.Param("id"))

	if err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
