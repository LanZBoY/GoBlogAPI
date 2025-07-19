package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"wentee/blog/app/schema/apperror"
	"wentee/blog/app/schema/apperror/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type testPayload struct {
	Name string `validate:"required"`
}

func TestErrorHandler_AppError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/", func(c *gin.Context) {
		c.Error(apperror.New(http.StatusBadRequest, errcode.BAD_REQUEST, nil))
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, errcode.BAD_REQUEST, resp["Code"])
}

func TestErrorHandler_ValidationErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandler())
	validate := validator.New()

	router.GET("/", func(c *gin.Context) {
		var p testPayload
		err := validate.Struct(p)
		c.Error(err)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}
