package user

import (
	"context"
	UserModel "wentee/blog/app/model/mongodb/user"
	"wentee/blog/app/schema/basemodel"
	UserSchema "wentee/blog/app/schema/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
