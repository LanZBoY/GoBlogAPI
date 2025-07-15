package post

import "go.mongodb.org/mongo-driver/mongo"

type PostRepo struct {
	postColletion *mongo.Collection
}

func NewPostRepo(postColletion *mongo.Collection) *PostRepo {
	return &PostRepo{
		postColletion: postColletion,
	}
}
