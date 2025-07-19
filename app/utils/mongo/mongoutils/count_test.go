package mongoutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"wentee/blog/app/testutils"
	"wentee/blog/app/utils/mongo/imongo"
)

type mockAggregator struct{ mock.Mock }

func (m *mockAggregator) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (imongo.Cursor, error) {
	args := m.Called(append([]any{ctx, pipeline}, opts...)...)
	return args.Get(0).(imongo.Cursor), args.Error(1)
}

func TestCountDocumentWithPipeline(t *testing.T) {
	ctx := context.TODO()
	pipeline := bson.A{bson.D{{Key: "stage", Value: 1}}}
	cursor := new(testutils.MockCursor)
	cursor.On("All", ctx, mock.Anything).Run(func(args mock.Arguments) {
		out := args.Get(1).(*[]bson.M)
		*out = []bson.M{{"total": int32(3)}}
	}).Return(nil)
	cursor.On("Close", ctx).Return(nil)

	aggr := new(mockAggregator)
	aggr.On("Aggregate", ctx, pipeline).Return(cursor, nil)

	m := &MongoUtils{}
	total, err := m.CountDocumentWithPipeline(ctx, aggr, pipeline)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	cursor.AssertExpectations(t)
	aggr.AssertExpectations(t)
}

func TestCountDocumentWithPipeline_Error(t *testing.T) {
	ctx := context.TODO()
	pipeline := bson.A{}
	aggr := new(mockAggregator)
	aggr.On("Aggregate", ctx, pipeline).Return((*testutils.MockCursor)(nil), assert.AnError)

	m := &MongoUtils{}
	total, err := m.CountDocumentWithPipeline(ctx, aggr, pipeline)
	assert.Error(t, err)
	assert.Equal(t, int64(0), total)
}
