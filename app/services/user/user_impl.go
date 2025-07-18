package user

import (
	"context"
	"net/http"
	UserModel "wentee/blog/app/model/mongodb/user"
	UserRepo "wentee/blog/app/repo/user"
	"wentee/blog/app/schema/apperror"
	"wentee/blog/app/schema/apperror/errcode"
	"wentee/blog/app/schema/basemodel"
	UserSchema "wentee/blog/app/schema/user"
	"wentee/blog/app/utils"
	"wentee/blog/app/utils/mongo/imongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService struct {
	userRepo      UserRepo.IUserRepository
	passwordUtils utils.IPasswrodUtils
	objectCreator imongo.IObjectIdCreator
}

func NewUserService(userRepo UserRepo.IUserRepository, passwordUtils utils.IPasswrodUtils, objectCreator imongo.IObjectIdCreator) *UserService {
	return &UserService{
		userRepo:      userRepo,
		passwordUtils: passwordUtils,
		objectCreator: objectCreator,
	}
}

func (svc *UserService) CountUsers(ctx context.Context) (total int64, err error) {
	return svc.userRepo.CountUsers(ctx)
}

func (svc *UserService) RegistryUser(ctx context.Context, createUser *UserSchema.UserCreate) error {

	if user, _ := svc.userRepo.GetUserByEmail(ctx, createUser.Email); user != nil {
		return apperror.New(http.StatusBadRequest, errcode.USER_EXIST, nil)
	}

	salt, err := svc.passwordUtils.GenerateSalt(32)
	if err != nil {
		return err
	}

	hashedPasswroed, err := svc.passwordUtils.HashPassword(createUser.Password, salt)

	if err != nil {
		return err
	}

	if err := svc.userRepo.CreateUser(ctx, &UserModel.UserDocument{
		Email:    createUser.Email,
		Username: createUser.Username,
		Password: hashedPasswroed,
		Salt:     salt,
	}); err != nil {
		return err
	}

	return nil
}

func (svc *UserService) GetUserById(ctx context.Context, id string) (*UserSchema.UserInfo, error) {
	oId, err := svc.objectCreator.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	userDoc, err := svc.userRepo.GetUserById(ctx, oId, options.FindOne().SetProjection(bson.M{UserModel.FieldPassword: 0}))

	if err != nil {
		return nil, err
	}

	return &UserSchema.UserInfo{
		Id:       userDoc.Id,
		Email:    userDoc.Email,
		Username: userDoc.Username,
	}, nil
}

func (svc *UserService) ListUsers(ctx context.Context, baseQuery *basemodel.BaseQuery) (users []UserSchema.UserInfo, err error) {

	userDocs, err := svc.userRepo.QueryUsers(ctx, baseQuery)

	if err != nil {
		return
	}

	for _, userDoc := range userDocs {
		users = append(users, UserSchema.UserInfo{
			Id:       userDoc.Id,
			Email:    userDoc.Email,
			Username: userDoc.Username,
		})
	}

	return

}

func (svc *UserService) UpdateUserById(ctx context.Context, id string, userUpdate *UserSchema.UserUpdate) (err error) {
	oId, err := svc.objectCreator.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	err = svc.userRepo.UpdateUserById(ctx, oId, userUpdate)
	return
}

func (svc *UserService) DeleteUserById(ctx context.Context, id string) (err error) {
	oid, err := svc.objectCreator.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	err = svc.userRepo.DeleteUserById(ctx, oid)
	return
}
