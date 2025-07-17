package testutils

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockCollection struct {
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

func (m *MockCollection) CountDocuments(ctx context.Context, filter any, opt ...*options.CountOptions) (int64, error) {
	params := appendCallArgs([]any{ctx, filter}, opt)
	args := m.Called(params...)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCollection) Find(ctx context.Context, filter any, opt ...*options.FindOptions) (*mongo.Cursor, error) {
	params := appendCallArgs([]any{ctx, filter}, opt)
	args := m.Called(params...)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockCollection) FindOne(ctx context.Context, filter any, opt ...*options.FindOneOptions) *mongo.SingleResult {
	params := appendCallArgs([]any{ctx, filter}, opt)
	args := m.Called(params...)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockCollection) UpdateOne(ctx context.Context, filter any, updateDoc any, opt ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	params := appendCallArgs([]any{ctx, filter, updateDoc}, opt)
	args := m.Called(params...)

	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockCollection) DeleteOne(ctx context.Context, filter any, opt ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	params := appendCallArgs([]any{ctx, filter}, opt)
	args := m.Called(params...)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (m *MockCollection) InsertOne(ctx context.Context, data any, opt ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	params := appendCallArgs([]any{ctx, data}, opt)
	args := m.Called(params...)

	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}
