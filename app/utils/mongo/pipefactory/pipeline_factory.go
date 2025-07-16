package pipefactory

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PipelineBuilder struct {
	stages    bson.A
	sortStage bson.D
	skip      *int64
	limit     *int64
}

func NewPipelineBuilder() *PipelineBuilder {
	return &PipelineBuilder{}
}

func (pb *PipelineBuilder) AddOperations(ops ...bson.D) *PipelineBuilder {
	for _, op := range ops {
		pb.stages = append(pb.stages, op)
	}
	return pb
}

func (pb *PipelineBuilder) AddSkip(skip int64) *PipelineBuilder {
	pb.skip = &skip
	return pb
}

func (pb *PipelineBuilder) AddLimit(limit int64) *PipelineBuilder {
	pb.limit = &limit
	return pb
}

func (pb *PipelineBuilder) AddSort(sort bson.D) *PipelineBuilder {
	pb.sortStage = sort
	return pb
}

func (pb *PipelineBuilder) BuildCountPipeline() bson.A {
	pipeline := bson.A{}
	pipeline = append(pipeline, pb.stages...)
	pipeline = append(pipeline, bson.D{{Key: "$count", Value: "total"}})
	return pipeline

}

func (pb *PipelineBuilder) BuildQeuryPipeline() bson.A {
	pipeline := bson.A{}
	pipeline = append(pipeline, pb.stages...)

	if len(pb.sortStage) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: pb.sortStage}})
	}
	if pb.skip != nil {
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: *pb.skip}})
	}
	if pb.limit != nil {
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: *pb.limit}})
	}

	return pipeline
}

func GetCount(ctx context.Context, col *mongo.Collection, countPipeline bson.A) (total int64, err error) {

	cursor, err := col.Aggregate(ctx, countPipeline)

	if err != nil {
		return
	}

	var countResult []bson.M
	if err = cursor.All(ctx, &countResult); err != nil {
		return
	}

	if len(countResult) > 0 {
		total = int64(countResult[0]["total"].(int32))
	}
	return
}
