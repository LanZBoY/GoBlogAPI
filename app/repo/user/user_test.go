package user_test

import (
	"context"
	"testing"
	"wentee/blog/app/repo/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockUserCollection struct {
	mock.Mock
}

func appendCallArgs[T any](fixed []any, variadic []T) []any {
	out := make([]any, 0, len(fixed)+len(variadic))

	out = append(out, fixed...)

	for _, v := range variadic {
		out = append(out, v)
	}
	return out
}

func (m *MockUserCollection) CountDocuments(ctx context.Context, filter any, opt ...*options.CountOptions) (int64, error) {
	params := appendCallArgs([]any{ctx, filter}, opt)
	args := m.Called(params...)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserCollection) Find(ctx context.Context, filter any, opt ...*options.FindOptions) (*mongo.Cursor, error) {
	params := appendCallArgs([]any{ctx, filter}, opt)
	args := m.Called(params...)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockUserCollection) FindOne(ctx context.Context, filter any, opt ...*options.FindOneOptions) *mongo.SingleResult {
	params := appendCallArgs([]any{ctx, filter}, opt)
	args := m.Called(params...)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockUserCollection) UpdateOne(ctx context.Context, filter any, updateDoc any, opt ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	params := appendCallArgs([]any{ctx, filter, updateDoc}, opt)
	args := m.Called(params...)

	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockUserCollection) DeleteOne(ctx context.Context, filter any, opt ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	params := appendCallArgs([]any{ctx, filter}, opt)
	args := m.Called(params...)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (m *MockUserCollection) InsertOne(ctx context.Context, data any, opt ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	params := appendCallArgs([]any{ctx, data}, opt)
	args := m.Called(params...)

	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func TestCountUser(t *testing.T) {
	userCol := new(MockUserCollection)
	userCol.On("CountDocuments", mock.Anything, bson.M{}).Return(int64(10), nil)

	userRepo := user.UserRepo{UserCollection: userCol}
	count, err := userRepo.CountUsers(context.TODO())

	assert.NoError(t, err)
	assert.Equal(t, int64(10), count)
	userCol.AssertCalled(t, "CountDocuments", mock.Anything, bson.M{})
}
