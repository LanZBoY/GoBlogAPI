package post

import (
	"context"
	"testing"
	PostModel "wentee/blog/app/model/mongodb/post"
	UserModel "wentee/blog/app/model/mongodb/user"
	"wentee/blog/app/schema/basemodel"
	PostSchema "wentee/blog/app/schema/post"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockPostRepo struct{ mock.Mock }

func (m *mockPostRepo) CreatePost(ctx context.Context, postCreate *PostSchema.PostCreate, createdBy *primitive.ObjectID) error {
	args := m.Called(ctx, postCreate, createdBy)
	return args.Error(0)
}
func (m *mockPostRepo) ListPosts(ctx context.Context, query *basemodel.BaseQuery) (int64, []PostModel.PostWithCreatorDocument, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(int64), args.Get(1).([]PostModel.PostWithCreatorDocument), args.Error(2)
}
func (m *mockPostRepo) GetPostById(ctx context.Context, id primitive.ObjectID) (PostModel.PostWithCreatorDocument, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(PostModel.PostWithCreatorDocument), args.Error(1)
}
func (m *mockPostRepo) UpdatePostById(ctx context.Context, id primitive.ObjectID, updateData *PostSchema.PostUpdate) error {
	args := m.Called(ctx, id, updateData)
	return args.Error(0)
}
func (m *mockPostRepo) DeletePostById(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestPostService_CreatePost(t *testing.T) {
	repo := new(mockPostRepo)
	svc := &PostService{postRepo: repo}
	create := &PostSchema.PostCreate{Title: "t"}
	oid := primitive.NewObjectID()
	repo.On("CreatePost", mock.Anything, create, &oid).Return(nil)

	err := svc.CreatePost(context.TODO(), create, oid.Hex())
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostService_CreatePost_BadId(t *testing.T) {
	repo := new(mockPostRepo)
	svc := &PostService{postRepo: repo}
	err := svc.CreatePost(context.TODO(), &PostSchema.PostCreate{}, "bad")
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestPostService_ListPosts(t *testing.T) {
	repo := new(mockPostRepo)
	svc := &PostService{postRepo: repo}
	query := &basemodel.BaseQuery{Skip: 0, Limit: 1}
	oid := primitive.NewObjectID()
	repo.On("ListPosts", mock.Anything, query).Return(int64(1), []PostModel.PostWithCreatorDocument{
		{
			PostDocument: PostModel.PostDocument{
				Id:    oid,
				Title: "t",
			},
			Creator: UserModel.UserDocument{
				Id:       oid,
				Username: "u",
			},
		},
	}, nil)

	total, posts, err := svc.ListPosts(context.TODO(), query)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, posts, 1)
	assert.Equal(t, oid, posts[0].Id)
	repo.AssertExpectations(t)
}

func TestPostService_GetPostById(t *testing.T) {
	repo := new(mockPostRepo)
	svc := &PostService{postRepo: repo}
	oid := primitive.NewObjectID()
	repo.On("GetPostById", mock.Anything, oid).Return(PostModel.PostWithCreatorDocument{PostDocument: PostModel.PostDocument{Id: oid, Title: "t", Content: ptr("c"), CreatedAt: time.Now()}, Creator: UserModel.UserDocument{Id: oid, Username: "u"}}, nil)

	post, err := svc.GetPostById(context.TODO(), oid.Hex())
	assert.NoError(t, err)
	assert.Equal(t, oid, post.Id)
	repo.AssertExpectations(t)
}

func TestPostService_GetPostById_BadId(t *testing.T) {
	svc := &PostService{postRepo: new(mockPostRepo)}
	_, err := svc.GetPostById(context.TODO(), "bad")
	assert.Error(t, err)
}

func TestPostService_UpdatePostById(t *testing.T) {
	repo := new(mockPostRepo)
	svc := &PostService{postRepo: repo}
	oid := primitive.NewObjectID()
	updateData := &PostSchema.PostUpdate{Title: ptr("new")}
	repo.On("UpdatePostById", mock.Anything, oid, updateData).Return(nil)

	err := svc.UpdatePostById(context.TODO(), oid.Hex(), updateData)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostService_UpdatePostById_BadId(t *testing.T) {
	repo := new(mockPostRepo)
	svc := &PostService{postRepo: repo}
	err := svc.UpdatePostById(context.TODO(), "bad", &PostSchema.PostUpdate{})
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestPostService_DeletePostById(t *testing.T) {
	repo := new(mockPostRepo)
	svc := &PostService{postRepo: repo}
	oid := primitive.NewObjectID()
	repo.On("DeletePostById", mock.Anything, oid).Return(nil)

	err := svc.DeletePostById(context.TODO(), oid.Hex())
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostService_DeletePostById_BadId(t *testing.T) {
	repo := new(mockPostRepo)
	svc := &PostService{postRepo: repo}
	err := svc.DeletePostById(context.TODO(), "bad")
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func ptr(s string) *string { return &s }
