package user

import (
	"context"
	"wentee/blog/app/model/mongodb"
	"wentee/blog/app/schema/basemodel"
	UserSchema "wentee/blog/app/schema/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo struct {
	userCollection *mongo.Collection
}

func NewUserRepo(userCollection *mongo.Collection) *UserRepo {

	return &UserRepo{
		userCollection: userCollection,
	}
}
func (repo *UserRepo) CountUsers() (int64, error) {
	return repo.userCollection.CountDocuments(context.TODO(), bson.M{})
}

func (repo *UserRepo) CreateUser(createUser *mongodb.UserDocument) error {
	_, err := repo.userCollection.InsertOne(context.TODO(), createUser)
	return err
}

func (repo *UserRepo) QueryUsers(query *basemodel.BaseQuery) ([]mongodb.UserDocument, error) {
	var userDocs []mongodb.UserDocument

	cur, err := repo.userCollection.Find(context.TODO(), bson.M{}, options.Find().SetSkip(query.Skip).SetLimit(query.Limit))

	if err != nil {
		return nil, err
	}

	if err := cur.All(context.TODO(), &userDocs); err != nil {
		return nil, err
	}

	return userDocs, nil
}

func (repo *UserRepo) GetUserById(id primitive.ObjectID, opts ...*options.FindOneOptions) (*mongodb.UserDocument, error) {
	var userDoc mongodb.UserDocument
	if err := repo.userCollection.FindOne(context.TODO(), bson.M{"_id": id}, opts...).Decode(&userDoc); err != nil {
		return nil, err
	}
	return &userDoc, nil
}

func (repo *UserRepo) GetUserByUserName(username string, opts ...*options.FindOneOptions) (*mongodb.UserDocument, error) {
	var userDoc mongodb.UserDocument
	if err := repo.userCollection.FindOne(context.TODO(), bson.M{"Username": username}, opts...).Decode(&userDoc); err != nil {
		return nil, err
	}
	return &userDoc, nil
}

func (repo *UserRepo) UpdateUserById(id primitive.ObjectID, updateData UserSchema.UserUpdate, opts ...*options.UpdateOptions) (err error) {
	_, err = repo.userCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": updateData})
	return
}

func (repo *UserRepo) DeleteUserById(id primitive.ObjectID, opts ...*options.DeleteOptions) (err error) {
	_, err = repo.userCollection.DeleteOne(context.TODO(), bson.M{"_id": id}, opts...)
	return
}
