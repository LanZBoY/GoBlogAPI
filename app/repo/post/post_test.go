package post_test

import (
	"context"
	"errors"
	"testing"
	PostModel "wentee/blog/app/model/mongodb/post"
	PostRepo "wentee/blog/app/repo/post"
	"wentee/blog/app/schema/basemodel"
	PostScheam "wentee/blog/app/schema/post"
	"wentee/blog/app/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewPostRepo(t *testing.T) {
	postCol := new(testutils.MockCollection)
	mockMongoUtils := new(testutils.MockMongoUtils)
	repo := PostRepo.NewPostRepo(postCol, mockMongoUtils)
	assert.NotNil(t, repo)
}

func TestCreatePost(t *testing.T) {
	postCol := new(testutils.MockCollection)
	mockMongoUtils := new(testutils.MockMongoUtils)

	ctx := context.TODO()
	postCreate := &PostScheam.PostCreate{}
	createdBy := primitive.NewObjectID()

	postCol.On("InsertOne", ctx, mock.Anything).Return(&mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil)

	repo := PostRepo.NewPostRepo(postCol, mockMongoUtils)

	err := repo.CreatePost(ctx, postCreate, &createdBy)
	assert.NoError(t, err)
}

func TestListPosts(t *testing.T) {

	tests := []struct {
		name      string
		ctx       context.Context
		query     *basemodel.BaseQuery
		mockSetup func(ctx context.Context, query *basemodel.BaseQuery, mc *testutils.MockCollection, mockUtils *testutils.MockMongoUtils, mockCursur *testutils.MockCursor)
		wantTotal int64
		wantErr   bool
	}{
		{
			name: "Normal",
			ctx:  context.TODO(),
			query: &basemodel.BaseQuery{
				Skip:  0,
				Limit: 10,
			},
			mockSetup: func(ctx context.Context, query *basemodel.BaseQuery, mc *testutils.MockCollection, mockUtils *testutils.MockMongoUtils, mockCursur *testutils.MockCursor) {
				mockUtils.On("CountDocumentWithPipeline", ctx, mc, mock.Anything).Return(int64(10), nil)
				mockCursur.On("All", ctx, mock.Anything).Run(func(args mock.Arguments) {
					out := args.Get(1).(*[]PostModel.PostWithCreatorDocument)
					*out = []PostModel.PostWithCreatorDocument{}
				}).Return(nil)

				mockCursur.On("Close", ctx).Return(nil)
				mc.On("Aggregate", ctx, mock.Anything).Return(mockCursur, nil)
			},
			wantTotal: 10,
			wantErr:   false,
		},
		{
			name: "Count Pipeline Error",
			ctx:  context.TODO(),
			query: &basemodel.BaseQuery{
				Skip:  0,
				Limit: 10,
			},
			mockSetup: func(ctx context.Context, query *basemodel.BaseQuery, mc *testutils.MockCollection, mockUtils *testutils.MockMongoUtils, mockCursur *testutils.MockCursor) {
				mockUtils.On("CountDocumentWithPipeline", ctx, mc, mock.Anything).Return(int64(0), errors.New("Something Happen"))
			},
			wantTotal: 0,
			wantErr:   true,
		},
		{
			name: "Aggregate Error",
			ctx:  context.TODO(),
			query: &basemodel.BaseQuery{
				Skip:  0,
				Limit: 10,
			},
			mockSetup: func(ctx context.Context, query *basemodel.BaseQuery, mc *testutils.MockCollection, mockUtils *testutils.MockMongoUtils, mockCursur *testutils.MockCursor) {
				mockUtils.On("CountDocumentWithPipeline", ctx, mc, mock.Anything).Return(int64(10), nil)
				mc.On("Aggregate", ctx, mock.Anything).Return((*mongo.Cursor)(nil), errors.New("Aggregate Error"))
			},
			wantTotal: 0,
			wantErr:   true,
		},
		{
			name: "Cursor All Error",
			ctx:  context.TODO(),
			query: &basemodel.BaseQuery{
				Skip:  0,
				Limit: 10,
			},
			mockSetup: func(ctx context.Context, query *basemodel.BaseQuery, mc *testutils.MockCollection, mockUtils *testutils.MockMongoUtils, mockCursur *testutils.MockCursor) {
				mockUtils.On("CountDocumentWithPipeline", ctx, mc, mock.Anything).Return(int64(10), nil)
				mockCursur.On("All", ctx, mock.Anything).Return(errors.New("CursorError"))
				mockCursur.On("Close", ctx).Return(nil)
				mc.On("Aggregate", ctx, mock.Anything).Return(mockCursur, nil)
			},
			wantTotal: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		mockCol := new(testutils.MockCollection)
		mockUtils := new(testutils.MockMongoUtils)
		mockCursor := new(testutils.MockCursor)

		tt.mockSetup(tt.ctx, tt.query, mockCol, mockUtils, mockCursor)

		repo := PostRepo.NewPostRepo(mockCol, mockUtils)

		total, _, err := repo.ListPosts(tt.ctx, tt.query)

		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.Equal(t, tt.wantTotal, total)
		}

	}
}

func TestGetPostById(t *testing.T) {

	tests := []struct {
		name      string
		ctx       context.Context
		id        primitive.ObjectID
		mockSetup func(ctx context.Context, id primitive.ObjectID, mockCol *testutils.MockCollection, mockCursor *testutils.MockCursor)
		wantErr   bool
	}{
		{
			name: "Normal",
			ctx:  context.TODO(),
			id:   primitive.NewObjectID(),
			mockSetup: func(ctx context.Context, id primitive.ObjectID, mockCol *testutils.MockCollection, mockCursor *testutils.MockCursor) {
				mockCol.On("Aggregate", ctx, mock.Anything).Return(mockCursor, nil)
				mockCursor.On("Next", ctx).Return(true)
				mockCursor.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
					out := args.Get(0).(*PostModel.PostWithCreatorDocument)
					*out = PostModel.PostWithCreatorDocument{}
				}).Return(nil)
				mockCursor.On("Close", ctx).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Aggregate Error",
			ctx:  context.TODO(),
			id:   primitive.NewObjectID(),
			mockSetup: func(ctx context.Context, id primitive.ObjectID, mockCol *testutils.MockCollection, mockCursor *testutils.MockCursor) {
				mockCol.On("Aggregate", ctx, mock.Anything).Return((*mongo.Cursor)(nil), errors.New("Some Error"))

			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		mockCol := new(testutils.MockCollection)
		mockCursor := new(testutils.MockCursor)
		mockUtils := new(testutils.MockMongoUtils)
		tt.mockSetup(tt.ctx, tt.id, mockCol, mockCursor)

		repo := PostRepo.NewPostRepo(mockCol, mockUtils)

		_, err := repo.GetPostById(tt.ctx, tt.id)

		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestUpdatePostById(t *testing.T) {
	mockCol := new(testutils.MockCollection)
	mockUtils := new(testutils.MockMongoUtils)
	ctx := context.TODO()
	id := primitive.NewObjectID()
	updateData := &PostScheam.PostUpdate{}
	mockCol.On("UpdateOne", ctx, bson.M{PostModel.FieldId: id}, bson.M{"$set": updateData}).Return(&mongo.UpdateResult{}, nil)

	repo := PostRepo.NewPostRepo(mockCol, mockUtils)
	err := repo.UpdatePostById(ctx, id, updateData)
	assert.NoError(t, err)
}

func TestDeletePostById(t *testing.T) {
	mockCol := new(testutils.MockCollection)
	mockUtils := new(testutils.MockMongoUtils)
	ctx := context.TODO()
	id := primitive.NewObjectID()
	mockCol.On("DeleteOne", ctx, bson.M{PostModel.FieldId: id}).Return(&mongo.UpdateResult{}, nil)
	repo := PostRepo.NewPostRepo(mockCol, mockUtils)

	err := repo.DeletePostById(ctx, id)
	assert.NoError(t, err)
}
