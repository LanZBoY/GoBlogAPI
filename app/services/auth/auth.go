package auth

import (
	"context"
	"net/http"
	"time"
	"wentee/blog/app/config"
	UserRepo "wentee/blog/app/repo/user"
	"wentee/blog/app/schema/apperror"
	"wentee/blog/app/schema/apperror/errcode"
	AuthSchema "wentee/blog/app/schema/auth"
	"wentee/blog/app/utils"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userRepo UserRepo.IGetUserByMail
}

func NewAuthService(userRepo *UserRepo.UserRepo) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (authSvc *AuthService) TryLogin(ctx context.Context, loginInfo *AuthSchema.LoginInfo) (tokenString string, err error) {

	userDoc, err := authSvc.userRepo.GetUserByEmail(ctx, loginInfo.Email)

	if err != nil {
		return
	}

	if !utils.VerifyPassword(userDoc.Password, loginInfo.Password, userDoc.Salt) {
		err = apperror.New(http.StatusNotFound, errcode.USER_NOT_FOUND, nil)
		return
	}

	claims := AuthSchema.JWTClaims{
		UserInfo: AuthSchema.JWTUserInfo{
			Id: userDoc.Id.Hex(),
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.SERVICE_NAME,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString([]byte(config.JWT_SECRET))

	return
}
