package post

import (
	"wentee/blog/app/utils/mongo/pipefactory"

	"go.mongodb.org/mongo-driver/bson"
)

func getPostWithCreatorListPipeline(skip, limit int64) (countPipeline, queryPipeline bson.A) {
	pb := pipefactory.NewPipelineBuilder()
	pb.AddOperations(
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "User"},
					{Key: "let", Value: bson.D{{Key: "userId", Value: "$CreatedBy"}}},
					{Key: "pipeline",
						Value: bson.A{
							bson.D{
								{Key: "$match",
									Value: bson.D{
										{Key: "$expr",
											Value: bson.D{
												{Key: "$eq",
													Value: bson.A{
														"$_id",
														"$$userId",
													},
												},
											},
										},
									},
								},
							},
							bson.D{
								{Key: "$project",
									Value: bson.D{
										{Key: "Password", Value: 0},
										{Key: "Salt", Value: 0},
									},
								},
							},
						},
					},
					{Key: "as", Value: "Creator"},
				},
			},
		},
		bson.D{{Key: "$project", Value: bson.D{{Key: "CreatedBy", Value: 0}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$Creator"}}}},
	).AddSkip(skip).AddLimit(limit).AddSort(bson.D{{Key: "CreatedAt", Value: -1}})
	countPipeline = pb.BuildCountPipeline()
	queryPipeline = pb.BuildQeuryPipeline()
	return
}
