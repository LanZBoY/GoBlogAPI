package user

import (
	"context"
	UserModel "wentee/blog/app/model/mongodb/user"
	"wentee/blog/app/schema/basemodel"
	UserSchema "wentee/blog/app/schema/user"
	"wentee/blog/app/utils/mongo/imongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUserCollection interface {
	imongo.ICountDocuments
	imongo.IFind
	imongo.IFindOne
	imongo.IInsertOne
	imongo.IUpdateOne
	imongo.IDeleteOne
}

type UserCollectionAdapter struct {
	Collection *mongo.Collection
}

// CountDocuments implements IUserCollection.
func (u *UserCollectionAdapter) CountDocuments(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error) {
	return u.Collection.CountDocuments(ctx, filter, opts...)
}

// DeleteOne implements IUserCollection.
func (u *UserCollectionAdapter) DeleteOne(ctx context.Context, filter any, opts ...*options.DeleteOptions) (imongo.DeleteResult, error) {
	return u.Collection.DeleteOne(ctx, filter, opts...)
}

// Find implements IUserCollection.
func (u *UserCollectionAdapter) Find(ctx context.Context, filter any, opts ...*options.FindOptions) (imongo.Cursor, error) {
	return u.Collection.Find(ctx, filter, opts...)
}

// FindOne implements IUserCollection.
func (u *UserCollectionAdapter) FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) imongo.SingleResult {
	return u.Collection.FindOne(ctx, filter, opts...)
}

// InsertOne implements IUserCollection.
func (u *UserCollectionAdapter) InsertOne(ctx context.Context, doc any, opts ...*options.InsertOneOptions) (imongo.InsertOneResult, error) {
	return u.Collection.InsertOne(ctx, doc, opts...)
}

// UpdateOne implements IUserCollection.
func (u *UserCollectionAdapter) UpdateOne(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (imongo.UpdateResult, error) {
	return u.Collection.UpdateOne(ctx, filter, update, opts...)
}

type IUserRepository interface {
	CountUsers(context.Context) (int64, error)
	CreateUser(context.Context, *UserModel.UserDocument) error
	QueryUsers(context.Context, *basemodel.BaseQuery) ([]UserModel.UserDocument, error)
	GetUserById(context.Context, primitive.ObjectID, ...*options.FindOneOptions) (*UserModel.UserDocument, error)
	UpdateUserById(context.Context, primitive.ObjectID, *UserSchema.UserUpdate, ...*options.UpdateOptions) error
	DeleteUserById(context.Context, primitive.ObjectID, ...*options.DeleteOptions) error
	IGetUserByMail
}

type IGetUserByMail interface {
	GetUserByEmail(context.Context, string, ...*options.FindOneOptions) (*UserModel.UserDocument, error)
}
