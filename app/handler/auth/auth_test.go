package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"wentee/blog/app/middleware"
	"wentee/blog/app/schema/apperror"
	"wentee/blog/app/schema/apperror/errcode"
	AuthSchema "wentee/blog/app/schema/auth"
	AuthSvc "wentee/blog/app/services/auth"
)

type mockAuthService struct{ mock.Mock }

func (m *mockAuthService) TryLogin(ctx context.Context, loginInfo *AuthSchema.LoginInfo) (string, error) {
	args := m.Called(ctx, loginInfo)
	return args.String(0), args.Error(1)
}

func setupRouter(svc AuthSvc.IAuthService) *gin.Engine {
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	ar := NewAuthRouter(svc)
	router.POST("/login", ar.Login)
	return router
}

func TestAuthRouter_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("ValidJSON", func(t *testing.T) {
		svc := new(mockAuthService)
		router := setupRouter(svc)
		svc.On("TryLogin", mock.Anything, mock.MatchedBy(func(li *AuthSchema.LoginInfo) bool {
			return li.Email == "a@b.c" && li.Password == "pw"
		})).Return("token", nil)

		body := `{"email":"a@b.c","password":"pw"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var resp struct {
			Data AuthSchema.TokenResponse `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "token", resp.Data.Token)
		svc.AssertExpectations(t)
	})

	t.Run("MissingFields", func(t *testing.T) {
		svc := new(mockAuthService)
		router := setupRouter(svc)

		body := `{"email":"a@b.c"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		svc.AssertNotCalled(t, "TryLogin", mock.Anything, mock.Anything)
	})

	t.Run("WrongPassword", func(t *testing.T) {
		svc := new(mockAuthService)
		router := setupRouter(svc)
		svc.On("TryLogin", mock.Anything, mock.MatchedBy(func(li *AuthSchema.LoginInfo) bool {
			return li.Email == "a@b.c" && li.Password == "bad"
		})).Return("", apperror.New(http.StatusNotFound, errcode.USER_NOT_FOUND, nil))

		body := `{"email":"a@b.c","password":"bad"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), errcode.USER_NOT_FOUND)
		svc.AssertExpectations(t)
	})
}
