package mongoutils

import (
	"context"
	"wentee/blog/app/utils/mongo/icollection"

	"go.mongodb.org/mongo-driver/bson"
)

func CountDocumentWithPipeline(ctx context.Context, aggregator icollection.IAggregate, countPipeline bson.A) (total int64, err error) {
	cursor, err := aggregator.Aggregate(ctx, countPipeline)

	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	var countResult []bson.M
	if err = cursor.All(ctx, &countResult); err != nil {
		return
	}

	if len(countResult) > 0 {
		total = int64(countResult[0]["total"].(int32))
	}
	return
}
