package reqcontext

import (
	"net/http"
	"wentee/blog/app/schema/apperror"
	"wentee/blog/app/schema/apperror/errcode"
	"wentee/blog/app/schema/auth"
	AuthSchema "wentee/blog/app/schema/auth"

	"github.com/gin-gonic/gin"
)

const USER_INFO = "UserInfo"

func GetUserInfo(c *gin.Context) (userInfo AuthSchema.JWTUserInfo, err error) {
	userInfoAny, ok := c.Get(USER_INFO)
	if !ok {
		err = apperror.New(http.StatusBadRequest, errcode.USER_NOT_FOUND, nil)
		return
	}
	userInfo, ok = userInfoAny.(auth.JWTUserInfo)
	if !ok {
		err = apperror.New(http.StatusInternalServerError, errcode.TYPE_ASSERTION_ERROR, nil)
		return
	}
	return
}
