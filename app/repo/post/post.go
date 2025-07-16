package post

import (
	"context"
	"time"
	PostModel "wentee/blog/app/model/mongodb/post"
	"wentee/blog/app/schema/basemodel"
	PostSchema "wentee/blog/app/schema/post"
	"wentee/blog/app/utils/mongo/pipefactory"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepo struct {
	postColletion *mongo.Collection
}

func NewPostRepo(postColletion *mongo.Collection) *PostRepo {
	return &PostRepo{
		postColletion: postColletion,
	}
}

func (repo *PostRepo) CreatePost(postCreate *PostSchema.PostCreate, createdBy *primitive.ObjectID) (err error) {
	_, err = repo.postColletion.InsertOne(context.TODO(), PostModel.PostDocument{
		Title:     postCreate.Title,
		Content:   postCreate.Content,
		CreatedAt: time.Now().UTC(),
		CreatedBy: *createdBy,
	})
	return
}

func (repo *PostRepo) ListPosts(query *basemodel.BaseQuery) (total int64, posts []PostModel.PostWithCreatorDocument, err error) {
	totalPipe, queryPipe := getPostWithCreatorListPipeline(query.Skip, query.Limit)
	total, err = pipefactory.GetCount(context.TODO(), repo.postColletion, totalPipe)
	if err != nil {
		return
	}
	cursor, err := repo.postColletion.Aggregate(context.TODO(), queryPipe)
	if err != nil {
		return
	}
	err = cursor.All(context.TODO(), &posts)
	return
}

func (repo *PostRepo) GetPostById(id primitive.ObjectID) (post PostModel.PostWithCreatorDocument, err error) {
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "User"},
					{
						Key:   "let",
						Value: bson.D{{Key: "userId", Value: "$CreatedBy"}},
					},
					{Key: "pipeline",
						Value: bson.A{
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
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$Creator"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
				},
			},
		},
		bson.D{{Key: "$project", Value: bson.D{{Key: "CreatedBy", Value: 0}}}},
	}

	cursor, err := repo.postColletion.Aggregate(context.TODO(), pipeline)

	if err != nil {
		return
	}

	if cursor.Next(context.TODO()) {
		cursor.Decode(&post)
	}

	return
}

func (repo *PostRepo) UpdatePostById(id primitive.ObjectID, updateData *PostSchema.PostUpdate) (err error) {
	_, err = repo.postColletion.UpdateOne(context.TODO(), bson.M{PostModel.FieldId: id}, bson.M{"$set": updateData})
	return
}

func (repo *PostRepo) DeletePostById(id primitive.ObjectID) (err error) {
	_, err = repo.postColletion.DeleteOne(context.TODO(), bson.M{PostModel.FieldId: id})
	return
}
