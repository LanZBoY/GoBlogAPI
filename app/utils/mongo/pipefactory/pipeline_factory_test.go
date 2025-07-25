package pipefactory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestPipelineBuilder_Build(t *testing.T) {
	pb := NewPipelineBuilder()
	pb.AddOperations(
		bson.D{{Key: "stage1", Value: "val1"}},
		bson.D{{Key: "stage2", Value: "val2"}},
	).AddSort(bson.D{{Key: "sort", Value: 1}}).AddSkip(5).AddLimit(10)

	count := pb.BuildCountPipeline()
	query := pb.BuildQeuryPipeline()

	expectedStages := bson.A{
		bson.D{{Key: "stage1", Value: "val1"}},
		bson.D{{Key: "stage2", Value: "val2"}},
	}

	assert.Equal(t, append(expectedStages, bson.D{{Key: "$count", Value: "total"}}), count)
	assert.Equal(t, append(expectedStages,
		bson.D{{Key: "$sort", Value: bson.D{{Key: "sort", Value: 1}}}},
		bson.D{{Key: "$skip", Value: int64(5)}},
		bson.D{{Key: "$limit", Value: int64(10)}},
	), query)
}

func TestPipelineBuilder_NoOptions(t *testing.T) {
	pb := NewPipelineBuilder()
	pb.AddOperations(bson.D{{Key: "only", Value: 1}})
	count := pb.BuildCountPipeline()
	query := pb.BuildQeuryPipeline()

	expected := bson.A{bson.D{{Key: "only", Value: 1}}}

	assert.Equal(t, append(expected, bson.D{{Key: "$count", Value: "total"}}), count)
	assert.Equal(t, expected, query)
}

func TestPipelineBuilder_WithSkipLimitSort(t *testing.T) {
	pb := NewPipelineBuilder()
	pb.AddOperations(bson.D{{Key: "s1", Value: "val"}}).
		AddSkip(3).AddLimit(6).AddSort(bson.D{{Key: "sort", Value: -1}})

	query := pb.BuildQeuryPipeline()

	assert.Equal(t, bson.A{
		bson.D{{Key: "s1", Value: "val"}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "sort", Value: -1}}}},
		bson.D{{Key: "$skip", Value: int64(3)}},
		bson.D{{Key: "$limit", Value: int64(6)}},
	}, query)
}
