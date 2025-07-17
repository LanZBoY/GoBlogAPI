package user

import (
	"context"
	"errors"
	"net/http"
	UserModel "wentee/blog/app/model/mongodb/user"
	"wentee/blog/app/schema/apperror"
	"wentee/blog/app/schema/apperror/errcode"
	"wentee/blog/app/schema/basemodel"
	UserSchema "wentee/blog/app/schema/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo struct {
	UserCollection IUserCollection
}

func NewUserRepo(userCollection *mongo.Collection) *UserRepo {

	return &UserRepo{
		UserCollection: userCollection,
	}
}
func (repo *UserRepo) CountUsers(ctx context.Context) (int64, error) {
	return repo.UserCollection.CountDocuments(ctx, bson.M{})
}

func (repo *UserRepo) CreateUser(ctx context.Context, createUser *UserModel.UserDocument) error {
	_, err := repo.UserCollection.InsertOne(ctx, createUser)
	return err
}

func (repo *UserRepo) QueryUsers(ctx context.Context, query *basemodel.BaseQuery) ([]UserModel.UserDocument, error) {
	var userDocs []UserModel.UserDocument

	cur, err := repo.UserCollection.Find(ctx, bson.M{}, options.Find().SetSkip(query.Skip).SetLimit(query.Limit))

	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &userDocs); err != nil {
		return nil, err
	}

	return userDocs, nil
}

func (repo *UserRepo) GetUserById(ctx context.Context, id primitive.ObjectID, opts ...*options.FindOneOptions) (*UserModel.UserDocument, error) {
	var userDoc UserModel.UserDocument
	if err := repo.UserCollection.FindOne(ctx, bson.M{UserModel.FieldId: id}, opts...).Decode(&userDoc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.New(http.StatusNotFound, errcode.USER_NOT_FOUND, err)
		}
		return nil, err
	}
	return &userDoc, nil
}

func (repo *UserRepo) GetUserByEmail(ctx context.Context, email string, opts ...*options.FindOneOptions) (*UserModel.UserDocument, error) {
	var userDoc UserModel.UserDocument
	if err := repo.UserCollection.FindOne(ctx, bson.M{UserModel.FieldEmail: email}, opts...).Decode(&userDoc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.New(http.StatusNotFound, errcode.USER_NOT_FOUND, err)
		}
		return nil, err
	}
	return &userDoc, nil
}

func (repo *UserRepo) UpdateUserById(ctx context.Context, id primitive.ObjectID, updateData *UserSchema.UserUpdate, opts ...*options.UpdateOptions) (err error) {
	_, err = repo.UserCollection.UpdateOne(ctx, bson.M{UserModel.FieldId: id}, bson.M{"$set": updateData})
	return
}

func (repo *UserRepo) DeleteUserById(ctx context.Context, id primitive.ObjectID, opts ...*options.DeleteOptions) (err error) {
	_, err = repo.UserCollection.DeleteOne(ctx, bson.M{UserModel.FieldId: id}, opts...)
	return
}
