package post

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"wentee/blog/app/middleware"
	AuthSchema "wentee/blog/app/schema/auth"
	"wentee/blog/app/schema/basemodel"
	PostSchema "wentee/blog/app/schema/post"
	"wentee/blog/app/utils/reqcontext"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockPostService struct{ mock.Mock }

func (m *mockPostService) CreatePost(ctx context.Context, postCreate *PostSchema.PostCreate, createdByString string) error {
	args := m.Called(ctx, postCreate, createdByString)
	return args.Error(0)
}

func (m *mockPostService) ListPosts(ctx context.Context, query *basemodel.BaseQuery) (int64, []PostSchema.PostList, error) {
	args := m.Called(ctx, query)
	var posts []PostSchema.PostList
	if v := args.Get(1); v != nil {
		posts = v.([]PostSchema.PostList)
	}
	return args.Get(0).(int64), posts, args.Error(2)
}

func (m *mockPostService) GetPostById(ctx context.Context, id string) (PostSchema.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(PostSchema.Post), args.Error(1)
}

func (m *mockPostService) UpdatePostById(ctx context.Context, id string, updateData *PostSchema.PostUpdate) error {
	args := m.Called(ctx, id, updateData)
	return args.Error(0)
}

func (m *mockPostService) DeletePostById(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupRouter(svc *mockPostService, withUser bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())
	if withUser {
		r.Use(func(c *gin.Context) {
			c.Set(reqcontext.USER_INFO, AuthSchema.JWTUserInfo{Id: "uid"})
		})
	}
	pr := &PostRouter{postSvc: svc}
	r.POST("/posts", pr.CreatePost)
	r.GET("/posts", pr.ListPosts)
	r.GET("/posts/:id", pr.GetPost)
	r.PATCH("/posts/:id", pr.UpdatePost)
	r.DELETE("/posts/:id", pr.DeletePost)
	return r
}

func TestNewPostRouter(t *testing.T) {
	r := NewPostRouter(nil)
	assert.NotNil(t, r)
}

func TestCreatePost(t *testing.T) {
	svc := new(mockPostService)
	router := setupRouter(svc, true)
	body := PostSchema.PostCreate{Title: "t"}
	bs, _ := json.Marshal(body)
	svc.On("CreatePost", mock.Anything, &body, "uid").Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestCreatePost_BindError(t *testing.T) {
	svc := new(mockPostService)
	router := setupRouter(svc, true)
	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(`{"content":"c"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	svc.AssertExpectations(t)
}

func TestCreatePost_ServiceError(t *testing.T) {
	svc := new(mockPostService)
	router := setupRouter(svc, true)
	body := PostSchema.PostCreate{Title: "t"}
	bs, _ := json.Marshal(body)
	svc.On("CreatePost", mock.Anything, &body, "uid").Return(errors.New("svc"))

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestListPosts(t *testing.T) {
	svc := new(mockPostService)
	router := setupRouter(svc, false)
	query := basemodel.NewDefaultQuery()
	oid := primitive.NewObjectID()
	posts := []PostSchema.PostList{{Id: oid, Title: "t", Creator: PostSchema.Creator{Id: oid, Username: "u"}}}
	svc.On("ListPosts", mock.Anything, &query).Return(int64(1), posts, nil)

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Total int64                 `json:"total"`
		Data  []PostSchema.PostList `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, oid, resp.Data[0].Id)
	svc.AssertExpectations(t)
}

func TestListPosts_Error(t *testing.T) {
	svc := new(mockPostService)
	router := setupRouter(svc, false)
	query := basemodel.NewDefaultQuery()
	svc.On("ListPosts", mock.Anything, &query).Return(int64(0), nil, errors.New("svc"))

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestListPosts_BindError(t *testing.T) {
	svc := new(mockPostService)
	router := setupRouter(svc, false)
	req := httptest.NewRequest(http.MethodGet, "/posts?limit=bad", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	svc.AssertExpectations(t)
}

func TestGetPost(t *testing.T) {
	svc := new(mockPostService)
	router := setupRouter(svc, false)
	oid := primitive.NewObjectID()
	post := PostSchema.Post{PostList: PostSchema.PostList{Id: oid, Title: "t", Creator: PostSchema.Creator{Id: oid, Username: "u"}}, Content: "c"}
	svc.On("GetPostById", mock.Anything, "id").Return(post, nil)

	req := httptest.NewRequest(http.MethodGet, "/posts/id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp basemodel.BaseResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	got, _ := json.Marshal(resp.Data)
	var postResp PostSchema.Post
	json.Unmarshal(got, &postResp)
	assert.Equal(t, oid, postResp.Id)
	svc.AssertExpectations(t)
}

func TestGetPost_Error(t *testing.T) {
	svc := new(mockPostService)
	router := setupRouter(svc, false)
	svc.On("GetPostById", mock.Anything, "id").Return(PostSchema.Post{}, errors.New("e"))

	req := httptest.NewRequest(http.MethodGet, "/posts/id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestUpdatePost(t *testing.T) {
	svc := new(mockPostService)
	router := setupRouter(svc, false)
	title := "new"
	upd := PostSchema.PostUpdate{Title: &title}
	bs, _ := json.Marshal(upd)
	svc.On("UpdatePostById", mock.Anything, "id", &upd).Return(nil)

	req := httptest.NewRequest(http.MethodPatch, "/posts/id", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestUpdatePost_Error(t *testing.T) {
	svc := new(mockPostService)
	router := setupRouter(svc, false)
	title := "new"
	upd := PostSchema.PostUpdate{Title: &title}
	svc.On("UpdatePostById", mock.Anything, "id", &upd).Return(errors.New("e"))
	bs, _ := json.Marshal(upd)

	req := httptest.NewRequest(http.MethodPatch, "/posts/id", bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestUpdatePost_BindError(t *testing.T) {
	svc := new(mockPostService)
	router := setupRouter(svc, false)
	req := httptest.NewRequest(http.MethodPatch, "/posts/id", bytes.NewBufferString(`{"title":1}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestDeletePost(t *testing.T) {
	svc := new(mockPostService)
	router := setupRouter(svc, false)
	svc.On("DeletePostById", mock.Anything, "id").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/posts/id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestDeletePost_Error(t *testing.T) {
	svc := new(mockPostService)
	router := setupRouter(svc, false)
	svc.On("DeletePostById", mock.Anything, "id").Return(errors.New("e"))

	req := httptest.NewRequest(http.MethodDelete, "/posts/id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}
