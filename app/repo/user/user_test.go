package user_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	UserModel "wentee/blog/app/model/mongodb/user"
	UserRepo "wentee/blog/app/repo/user"
	"wentee/blog/app/schema/basemodel"
	"wentee/blog/app/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func TestQueryUsers(t *testing.T) {

	tests := []struct {
		ctx   context.Context
		query *basemodel.BaseQuery
	}{
		{
			ctx: context.TODO(),
			query: &basemodel.BaseQuery{
				Skip:  0,
				Limit: 10,
			},
		},
		{
			ctx: context.TODO(),
			query: &basemodel.BaseQuery{
				Skip:  0,
				Limit: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Skip = %d, Limit = %d", tt.query.Skip, tt.query.Limit), func(t *testing.T) {
			userCol := new(testutils.MockCollection)

			findOpts := options.Find()
			findOpts.SetSkip(tt.query.Skip)
			findOpts.SetLimit(tt.query.Limit)

			if tt.query.Limit > 0 {
				userCol.On("Find", tt.ctx, mock.Anything, findOpts).
					Return(&mongo.Cursor{}, nil) // 假設用假的 cursor
			} else {
				userCol.On("Find", tt.ctx, mock.Anything, findOpts).
					Return((*mongo.Cursor)(nil), errors.New("Litmit must be positive"))
			}

			repo := &UserRepo.UserRepo{
				UserCollection: userCol,
			}
			_, err := repo.UserCollection.Find(tt.ctx, bson.M{}, findOpts)

			if tt.query.Limit > 0 {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

		})
	}
}
