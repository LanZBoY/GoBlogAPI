package user

import (
	"context"
	"wentee/blog/app/schema/basemodel"
	UserSchema "wentee/blog/app/schema/user"
)

type IUserService interface {
	CountUsers(ctx context.Context) (total int64, err error)
	RegistryUser(ctx context.Context, createUser *UserSchema.UserCreate) error
	GetUserById(ctx context.Context, id string) (*UserSchema.UserInfo, error)
	ListUsers(ctx context.Context, baseQuery *basemodel.BaseQuery) (users []UserSchema.UserInfo, err error)
	UpdateUserById(ctx context.Context, id string, userUpdate *UserSchema.UserUpdate) (err error)
	DeleteUserById(ctx context.Context, id string) (err error)
}
