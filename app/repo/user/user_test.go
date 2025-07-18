package user_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	UserModel "wentee/blog/app/model/mongodb/user"
	UserRepo "wentee/blog/app/repo/user"
	"wentee/blog/app/schema/apperror"
	"wentee/blog/app/schema/apperror/errcode"
	"wentee/blog/app/schema/basemodel"
	UserSchema "wentee/blog/app/schema/user"
	"wentee/blog/app/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestNewUserRepo(t *testing.T) {
	userCol := new(testutils.MockCollection)

	repo := UserRepo.NewUserRepo(userCol)

	assert.NotNil(t, repo)
}

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
		name       string
		ctx        context.Context
		query      *basemodel.BaseQuery
		mockSetup  func(*testutils.MockCollection, *testutils.MockCursor, context.Context, *options.FindOptions)
		wantNormal bool
	}{
		{
			name: "Success",
			query: &basemodel.BaseQuery{
				Skip:  0,
				Limit: 10,
			},
			ctx: context.TODO(),
			mockSetup: func(userCol *testutils.MockCollection, mockCursor *testutils.MockCursor, ctx context.Context, findOpts *options.FindOptions) {
				mockCursor.On("All", ctx, mock.Anything).Return(nil)
				userCol.On("Find", ctx, mock.Anything, findOpts).
					Return(mockCursor, nil) // 假設用假的 cursor
			},
			wantNormal: true,
		},
		{
			name: "Limit Error",
			query: &basemodel.BaseQuery{
				Skip:  0,
				Limit: 0,
			},
			ctx: context.TODO(),
			mockSetup: func(userCol *testutils.MockCollection, mockCursor *testutils.MockCursor, ctx context.Context, findOpts *options.FindOptions) {
				userCol.On("Find", ctx, mock.Anything, findOpts).
					Return((*testutils.MockCursor)(nil), errors.New("Litmit must be positive"))
				mockCursor.On("All", ctx, mock.Anything).Return(nil)
			},
			wantNormal: false,
		},
		{
			name: "Cursor Error",
			query: &basemodel.BaseQuery{
				Skip:  0,
				Limit: 0,
			},
			ctx: context.TODO(),
			mockSetup: func(userCol *testutils.MockCollection, mockCursor *testutils.MockCursor, ctx context.Context, findOpts *options.FindOptions) {
				mockCursor.On("All", ctx, mock.Anything).Return(errors.New("MockErr"))
				userCol.On("Find", ctx, mock.Anything, findOpts).
					Return(mockCursor, nil)
			},
			wantNormal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userCol := new(testutils.MockCollection)
			mockCursor := new(testutils.MockCursor)
			findOpts := options.Find()
			findOpts.SetSkip(tt.query.Skip)
			findOpts.SetLimit(tt.query.Limit)

			tt.mockSetup(userCol, mockCursor, tt.ctx, findOpts)

			repo := &UserRepo.UserRepo{
				UserCollection: userCol,
			}
			_, err := repo.QueryUsers(tt.ctx, tt.query)

			if tt.wantNormal {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

		})
	}
}

func TestGetUserById(t *testing.T) {

	tests := []struct {
		name      string
		ctx       context.Context
		id        primitive.ObjectID
		mockSetup func(*testutils.MockCollection, *testutils.MockSingleResult, context.Context, primitive.ObjectID)
		wantErr   bool
	}{
		{
			name: "Normal",
			ctx:  context.TODO(),
			id:   primitive.NewObjectID(),
			mockSetup: func(mc *testutils.MockCollection, msr *testutils.MockSingleResult, ctx context.Context, id primitive.ObjectID) {

				msr.On("Decode", mock.Anything).Return(nil)
				mc.On("FindOne", ctx, bson.M{UserModel.FieldId: id}).Return(msr)
			},
			wantErr: false,
		},
		{
			name: "NotFound Error",
			ctx:  context.TODO(),
			id:   primitive.NewObjectID(),
			mockSetup: func(mc *testutils.MockCollection, msr *testutils.MockSingleResult, ctx context.Context, id primitive.ObjectID) {

				msr.On("Decode", mock.Anything).Return(mongo.ErrNoDocuments)
				mc.On("FindOne", ctx, bson.M{UserModel.FieldId: id}).Return(msr)
			},
			wantErr: true,
		},
		{
			name: "Unknown Error",
			ctx:  context.TODO(),
			id:   primitive.NewObjectID(),
			mockSetup: func(mc *testutils.MockCollection, msr *testutils.MockSingleResult, ctx context.Context, id primitive.ObjectID) {

				msr.On("Decode", mock.Anything).Return(errors.New("Mock"))
				mc.On("FindOne", ctx, bson.M{UserModel.FieldId: id}).Return(msr)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := new(testutils.MockCollection)
			msr := new(testutils.MockSingleResult)

			tt.mockSetup(mc, msr, tt.ctx, tt.id)

			repo := &UserRepo.UserRepo{UserCollection: mc}

			_, err := repo.GetUserById(tt.ctx, tt.id)

			if tt.wantErr {
				appError := &apperror.AppError{}
				if errors.As(err, appError) {
					assert.Equal(t, http.StatusNotFound, appError.Status)
					assert.Equal(t, errcode.USER_NOT_FOUND, appError.Code)
					assert.Equal(t, mongo.ErrNoDocuments, appError.Err)
				} else {
					assert.Error(t, err)
				}

			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestGetUserByEmail(t *testing.T) {

	tests := []struct {
		name      string
		ctx       context.Context
		email     string
		mockSetup func(*testutils.MockCollection, *testutils.MockSingleResult, context.Context, string)
		wantErr   bool
	}{
		{
			name:  "Normal",
			ctx:   context.TODO(),
			email: "superjimalex@ispi.com.tw",
			mockSetup: func(mc *testutils.MockCollection, msr *testutils.MockSingleResult, ctx context.Context, email string) {

				msr.On("Decode", mock.Anything).Return(nil)
				mc.On("FindOne", ctx, bson.M{UserModel.FieldEmail: email}).Return(msr)
			},
			wantErr: false,
		},
		{
			name:  "NotFound Error",
			ctx:   context.TODO(),
			email: "superjimalex@ispi.com.tw",
			mockSetup: func(mc *testutils.MockCollection, msr *testutils.MockSingleResult, ctx context.Context, email string) {

				msr.On("Decode", mock.Anything).Return(mongo.ErrNoDocuments)
				mc.On("FindOne", ctx, bson.M{UserModel.FieldEmail: email}).Return(msr)
			},
			wantErr: true,
		},
		{
			name:  "Unknown Error",
			ctx:   context.TODO(),
			email: "superjimalex@ispi.com.tw",
			mockSetup: func(mc *testutils.MockCollection, msr *testutils.MockSingleResult, ctx context.Context, email string) {

				msr.On("Decode", mock.Anything).Return(errors.New("Mock"))
				mc.On("FindOne", ctx, bson.M{UserModel.FieldEmail: email}).Return(msr)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := new(testutils.MockCollection)
			msr := new(testutils.MockSingleResult)

			tt.mockSetup(mc, msr, tt.ctx, tt.email)

			repo := &UserRepo.UserRepo{UserCollection: mc}

			_, err := repo.GetUserByEmail(tt.ctx, tt.email)

			if tt.wantErr {
				appError := &apperror.AppError{}
				if errors.As(err, appError) {
					assert.Equal(t, http.StatusNotFound, appError.Status)
					assert.Equal(t, errcode.USER_NOT_FOUND, appError.Code)
					assert.Equal(t, mongo.ErrNoDocuments, appError.Err)
				} else {
					assert.Error(t, err)
				}

			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateUserById(t *testing.T) {

	tests := []struct {
		name       string
		ctx        context.Context
		id         primitive.ObjectID
		updateData *UserSchema.UserUpdate
		mockSetup  func(*testutils.MockCollection, context.Context, primitive.ObjectID, *UserSchema.UserUpdate)
	}{
		{
			name:       "Normal",
			ctx:        context.TODO(),
			id:         primitive.NewObjectID(),
			updateData: &UserSchema.UserUpdate{},
			mockSetup: func(mc *testutils.MockCollection, ctx context.Context, id primitive.ObjectID, updateData *UserSchema.UserUpdate) {
				mc.On("UpdateOne", ctx, bson.M{UserModel.FieldId: id}, bson.M{"$set": updateData}).Return(&mongo.UpdateResult{}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := new(testutils.MockCollection)
			tt.mockSetup(mc, tt.ctx, tt.id, tt.updateData)

			repo := &UserRepo.UserRepo{
				UserCollection: mc,
			}

			err := repo.UpdateUserById(tt.ctx, tt.id, tt.updateData)

			assert.NoError(t, err)
		})
	}

}
func TestDeleteUserById(t *testing.T) {

	tests := []struct {
		name      string
		ctx       context.Context
		id        primitive.ObjectID
		mockSetup func(*testutils.MockCollection, context.Context, primitive.ObjectID)
	}{
		{
			name: "Normal",
			ctx:  context.TODO(),
			id:   primitive.NewObjectID(),
			mockSetup: func(mc *testutils.MockCollection, ctx context.Context, id primitive.ObjectID) {
				mc.On("DeleteOne", ctx, bson.M{UserModel.FieldId: id}).Return(&mongo.DeleteResult{DeletedCount: 1}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := new(testutils.MockCollection)
			tt.mockSetup(mc, tt.ctx, tt.id)

			repo := &UserRepo.UserRepo{
				UserCollection: mc,
			}

			err := repo.DeleteUserById(tt.ctx, tt.id)

			assert.NoError(t, err)
		})
	}
}
