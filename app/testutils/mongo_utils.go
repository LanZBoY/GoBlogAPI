package testutils

import (
	"context"
	"wentee/blog/app/utils/mongo/imongo"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
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

func (m *MockCollection) CountDocuments(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error) {
	params := appendCallArgs([]any{ctx, filter}, opts)
	args := m.Called(params...)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCollection) Find(ctx context.Context, filter any, opts ...*options.FindOptions) (imongo.Cursor, error) {
	params := appendCallArgs([]any{ctx, filter}, opts)
	args := m.Called(params...)
	return args.Get(0).(imongo.Cursor), args.Error(1)
}

func (m *MockCollection) FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) imongo.SingleResult {
	params := appendCallArgs([]any{ctx, filter}, opts)
	args := m.Called(params...)
	return args.Get(0).(imongo.SingleResult)
}

func (m *MockCollection) UpdateOne(ctx context.Context, filter any, updateDoc any, opts ...*options.UpdateOptions) (imongo.UpdateResult, error) {
	params := appendCallArgs([]any{ctx, filter, updateDoc}, opts)
	args := m.Called(params...)

	return args.Get(0).(imongo.UpdateResult), args.Error(1)
}

func (m *MockCollection) DeleteOne(ctx context.Context, filter any, opts ...*options.DeleteOptions) (imongo.DeleteResult, error) {
	params := appendCallArgs([]any{ctx, filter}, opts)
	args := m.Called(params...)
	return args.Get(0).(imongo.DeleteResult), args.Error(1)
}

func (m *MockCollection) InsertOne(ctx context.Context, data any, opts ...*options.InsertOneOptions) (imongo.InsertOneResult, error) {
	params := appendCallArgs([]any{ctx, data}, opts)
	args := m.Called(params...)

	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockCollection) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (imongo.Cursor, error) {
	params := appendCallArgs([]any{ctx, pipeline}, opts)
	args := m.Called(params...)
	return args.Get(0).(imongo.Cursor), args.Error(1)
}

type MockCursor struct {
	mock.Mock
}

func (m *MockCursor) All(ctx context.Context, out interface{}) error {
	args := m.Called(ctx, out)
	return args.Error(0)
}

func (m *MockCursor) Next(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockCursor) Decode(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockCursor) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockSingleResult struct {
	mock.Mock
}

func (m *MockSingleResult) Decode(v any) error {
	args := m.Called(v)
	return args.Error(0)
}

type MockMongoUtils struct {
	mock.Mock
}

// CountDocumentWithPipeline implements mongoutils.DocumentCounter.
func (m *MockMongoUtils) CountDocumentWithPipeline(ctx context.Context, aggregator imongo.IAggregate, countPipeline bson.A) (total int64, err error) {
	args := m.Called(ctx, aggregator, countPipeline)
	return args.Get(0).(int64), args.Error(1)
}
