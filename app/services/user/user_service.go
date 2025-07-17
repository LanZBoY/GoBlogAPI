package user

import (
	"context"
	"wentee/blog/app/schema/basemodel"
	UserSchema "wentee/blog/app/schema/user"
)

type IUserService interface {
	CountUsers(context.Context) (int64, error)
	RegistryUser(context.Context, *UserSchema.UserCreate) error
	GetUserById(context.Context, string) (*UserSchema.UserInfo, error)
	ListUsers(context.Context, *basemodel.BaseQuery) ([]UserSchema.UserInfo, error)
	UpdateUserById(context.Context, string, *UserSchema.UserUpdate) error
	DeleteUserById(context.Context, string) error
}
