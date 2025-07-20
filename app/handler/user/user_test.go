package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"wentee/blog/app/middleware"
	"wentee/blog/app/schema/apperror"
	"wentee/blog/app/schema/apperror/errcode"
	AuthSchema "wentee/blog/app/schema/auth"
	"wentee/blog/app/schema/basemodel"
	UserSchema "wentee/blog/app/schema/user"
	"wentee/blog/app/utils/reqcontext"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mockUserService satisfies user.IUserService using testify/mock.
type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) CountUsers(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUserService) RegistryUser(ctx context.Context, createUser *UserSchema.UserCreate) error {
	args := m.Called(ctx, createUser)
	return args.Error(0)
}

func (m *mockUserService) GetUserById(ctx context.Context, id string) (*UserSchema.UserInfo, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*UserSchema.UserInfo), args.Error(1)
}

func (m *mockUserService) ListUsers(ctx context.Context, baseQuery *basemodel.BaseQuery) ([]UserSchema.UserInfo, error) {
	args := m.Called(ctx, baseQuery)
	return args.Get(0).([]UserSchema.UserInfo), args.Error(1)
}

func (m *mockUserService) UpdateUserById(ctx context.Context, id string, userUpdate *UserSchema.UserUpdate) error {
	args := m.Called(ctx, id, userUpdate)
	return args.Error(0)
}

func (m *mockUserService) DeleteUserById(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupRouter(us *mockUserService, mws ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())
	group := r.Group("/users", mws...)
	ur := &UserRouter{userSvc: us}
	group.POST("", ur.CreateUser)
	group.GET("", ur.ListUsers)
	group.GET("/me", ur.GetMe)
	idGroup := group.Group("/:id")
	{
		idGroup.GET("", ur.GetUser)
		idGroup.PATCH("", ur.UpdateUser)
		idGroup.DELETE("", ur.DeleteUser)
	}
	return r
}

func TestCreateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mockUserService)
		router := setupRouter(mockSvc)
		reqBody := UserSchema.UserCreate{Email: "a@b.com", Username: "user", Password: "pwd"}
		body, _ := json.Marshal(reqBody)
		mockSvc.On("RegistryUser", mock.Anything, &reqBody).Return(nil)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "", strings.TrimSpace(w.Body.String()))
		mockSvc.AssertExpectations(t)
	})

	t.Run("ValidationError", func(t *testing.T) {
		mockSvc := new(mockUserService)
		router := setupRouter(mockSvc)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{}`))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(mockUserService)
		router := setupRouter(mockSvc)
		reqBody := UserSchema.UserCreate{Email: "a@b.com", Username: "user", Password: "pwd"}
		body, _ := json.Marshal(reqBody)
		mockSvc.On("RegistryUser", mock.Anything, &reqBody).Return(errors.New("svc"))
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mockUserService)
		router := setupRouter(mockSvc)
		uid := primitive.NewObjectID()
		out := &UserSchema.UserInfo{Id: uid, Email: "a@b.com", Username: "user"}
		mockSvc.On("GetUserById", mock.Anything, "1").Return(out, nil)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		router.ServeHTTP(w, req)
		var resp struct {
			Data UserSchema.UserInfo `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, out.Email, resp.Data.Email)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(mockUserService)
		router := setupRouter(mockSvc)
		mockSvc.On("GetUserById", mock.Anything, "1").Return((*UserSchema.UserInfo)(nil), apperror.New(http.StatusBadRequest, errcode.BAD_REQUEST, nil))
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestListUsers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mockUserService)
		router := setupRouter(mockSvc)
		q := basemodel.NewDefaultQuery()
		users := []UserSchema.UserInfo{{Id: primitive.NewObjectID(), Email: "a@b.com", Username: "u"}}
		mockSvc.On("CountUsers", mock.Anything).Return(int64(1), nil)
		mockSvc.On("ListUsers", mock.Anything, &q).Return(users, nil)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		router.ServeHTTP(w, req)
		var resp struct {
			Total int64                 `json:"total"`
			Data  []UserSchema.UserInfo `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, int64(1), resp.Total)
		mockSvc.AssertExpectations(t)
	})

	t.Run("BindError", func(t *testing.T) {
		mockSvc := new(mockUserService)
		router := setupRouter(mockSvc)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users?limit=bad", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(mockUserService)
		router := setupRouter(mockSvc)
		mockSvc.On("CountUsers", mock.Anything).Return(int64(0), errors.New("fail"))
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestGetMe(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mockUserService)
		mw := func(c *gin.Context) {
			c.Set(reqcontext.USER_INFO, AuthSchema.JWTUserInfo{Id: "1"})
			c.Next()
		}
		router := setupRouter(mockSvc, mw)
		out := &UserSchema.UserInfo{Id: primitive.NewObjectID(), Email: "a@b.com", Username: "user"}
		mockSvc.On("GetUserById", mock.Anything, "1").Return(out, nil)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
		router.ServeHTTP(w, req)
		var resp struct {
			Data UserSchema.UserInfo `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, out.Username, resp.Data.Username)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(mockUserService)
		mw := func(c *gin.Context) {
			c.Set(reqcontext.USER_INFO, AuthSchema.JWTUserInfo{Id: "1"})
			c.Next()
		}
		router := setupRouter(mockSvc, mw)
		mockSvc.On("GetUserById", mock.Anything, "1").Return((*UserSchema.UserInfo)(nil), errors.New("fail"))
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mockUserService)
		router := setupRouter(mockSvc)
		email := "a@b.com"
		up := UserSchema.UserUpdate{Email: &email}
		body, _ := json.Marshal(up)
		mockSvc.On("UpdateUserById", mock.Anything, "1", &up).Return(nil)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/users/1", bytes.NewReader(body))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(mockUserService)
		router := setupRouter(mockSvc)
		email := "a@b.com"
		up := UserSchema.UserUpdate{Email: &email}
		body, _ := json.Marshal(up)
		mockSvc.On("UpdateUserById", mock.Anything, "1", &up).Return(errors.New("fail"))
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/users/1", bytes.NewReader(body))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mockUserService)
		router := setupRouter(mockSvc)
		mockSvc.On("DeleteUserById", mock.Anything, "1").Return(nil)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(mockUserService)
		router := setupRouter(mockSvc)
		mockSvc.On("DeleteUserById", mock.Anything, "1").Return(errors.New("fail"))
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
