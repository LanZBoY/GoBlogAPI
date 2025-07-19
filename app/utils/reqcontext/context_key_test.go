package reqcontext

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	apperror "wentee/blog/app/schema/apperror"
	"wentee/blog/app/schema/apperror/errcode"
	AuthSchema "wentee/blog/app/schema/auth"
)

func TestGetUserInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	want := AuthSchema.JWTUserInfo{Id: "123"}
	c.Set(USER_INFO, want)

	userInfo, err := GetUserInfo(c)
	assert.NoError(t, err)
	assert.Equal(t, want, userInfo)
}

func TestGetUserInfo_Missing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	_, err := GetUserInfo(c)
	assert.Error(t, err)
	if appErr, ok := err.(apperror.AppError); ok {
		assert.Equal(t, errcode.USER_NOT_FOUND, appErr.Code)
		assert.Equal(t, http.StatusBadRequest, appErr.Status)
	}
}

func TestGetUserInfo_BadType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	c.Set(USER_INFO, 1)

	_, err := GetUserInfo(c)
	assert.Error(t, err)
	if appErr, ok := err.(apperror.AppError); ok {
		assert.Equal(t, errcode.TYPE_ASSERTION_ERROR, appErr.Code)
		assert.Equal(t, http.StatusInternalServerError, appErr.Status)
	}
}
