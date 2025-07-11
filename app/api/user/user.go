package user

import (
	"wentee/blog/app/service/user"

	"github.com/gin-gonic/gin"
)

type UserAPI struct {
	UserService user.UserService
}

func (api *UserAPI) GetUser(c *gin.Context) {

}
