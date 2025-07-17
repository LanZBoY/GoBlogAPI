package user_test

import (
	"context"
	"testing"
	UserModel "wentee/blog/app/model/mongodb/user"
	UserRepo "wentee/blog/app/repo/user"
	"wentee/blog/app/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCountUser(t *testing.T) {
	userCol := new(testutils.MockCollection)
	userCol.On("CountDocuments", mock.Anything, bson.M{}).Return(int64(10), nil)

	userRepo := UserRepo.UserRepo{UserCollection: userCol}
	count, err := userRepo.CountUsers(context.TODO())

	assert.NoError(t, err)
	assert.Equal(t, int64(10), count)
	userCol.AssertCalled(t, "CountDocuments", mock.Anything, bson.M{})
}

func TestCreateUser(t *testing.T) {
	userCol := new(testutils.MockCollection)
	insertData := &UserModel.UserDocument{
		Id:       primitive.NewObjectID(),
		Email:    "WenTee@ispi.com.tw",
		Username: "Cool",
		Password: "string",
		Salt:     "Fake",
	}
	userCol.On("InsertOne", mock.Anything, insertData).Return(&mongo.InsertOneResult{InsertedID: insertData.Id}, nil)

	userRepo := UserRepo.UserRepo{UserCollection: userCol}
	err := userRepo.CreateUser(context.TODO(), insertData)
	assert.NoError(t, err)
}

// func TestQueryUsers(t *testing.T) {
// 	userCol := new(testutils.MockCollection)

// 	tests := []struct {
// 		ctx   context.Context
// 		query *basemodel.BaseQuery
// 	}{
// 		{
// 			ctx: context.TODO(),
// 			query: &basemodel.BaseQuery{
// 				Skip:  0,
// 				Limit: 10,
// 			},
// 		},
// 		{
// 			ctx: context.TODO(),
// 			query: &basemodel.BaseQuery{
// 				Skip:  0,
// 				Limit: 0,
// 			},
// 		},
// 	}
// 	userCol.On("Find", )
// }
