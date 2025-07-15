package user

import (
	"net/http"
	"wentee/blog/app/model/mongodb"
	UserRepo "wentee/blog/app/repo/user"
	"wentee/blog/app/schema/apperror"
	"wentee/blog/app/schema/apperror/errcode"
	"wentee/blog/app/schema/basemodel"
	UserSchema "wentee/blog/app/schema/user"
	"wentee/blog/app/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService struct {
	userRepo *UserRepo.UserRepo
}

func NewUserService(userRepo *UserRepo.UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (svc *UserService) CountUsers() (total int64, err error) {
	return svc.userRepo.CountUsers()
}

func (svc *UserService) RegistryUser(createUser *UserSchema.UserCreate) error {

	if user, _ := svc.userRepo.GetUserByUserName(createUser.Username); user != nil {
		return apperror.New(http.StatusBadRequest, errcode.USER_EXIST, nil)
	}

	salt, err := utils.GenerateSalt(32)
	if err != nil {
		return err
	}

	hashedPasswroed, err := utils.HashPassword(createUser.Password, salt)

	if err != nil {
		return err
	}

	svc.userRepo.CreateUser(&mongodb.UserDocument{
		Username: createUser.Username,
		Password: hashedPasswroed,
		Salt:     salt,
	})

	return nil
}

func (svc *UserService) GetUserById(id string) (*UserSchema.UserInfo, error) {
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	userDoc, err := svc.userRepo.GetUserById(oId, options.FindOne().SetProjection(bson.M{"Password": 0}))

	if err != nil {
		return nil, err
	}

	return &UserSchema.UserInfo{
		Username: userDoc.Username,
	}, nil
}

func (svc *UserService) ListUsers(baseQuery *basemodel.BaseQuery) (users []UserSchema.UserInfo, err error) {

	user_docs, err := svc.userRepo.QueryUsers(baseQuery)

	if err != nil {
		return
	}

	for _, user_doc := range user_docs {
		users = append(users, UserSchema.UserInfo{
			Username: user_doc.Username,
		})
	}

	return

}
