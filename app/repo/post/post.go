package post

import (
	"context"
	"time"
	"wentee/blog/app/model/mongodb"
	PostModel "wentee/blog/app/model/mongodb/post"
	PostSchema "wentee/blog/app/schema/post"

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
		CreatedAt: time.Now(),
		CreatedBy: *createdBy,
	})
	return
}

func (repo *PostRepo) ListPost() (posts []PostModel.PostWithCreatorDocument, err error) {
	cursor, err := repo.postColletion.Aggregate(context.TODO(), mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{{Key: "from", Value: mongodb.UserCollection}, {Key: "localField", Value: "CreatedBy"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "Creator"}}}},
		{{Key: "$unwind", Value: "$Creator"}},
	})
	if err != nil {
		return
	}
	err = cursor.All(context.TODO(), &posts)
	return
}
