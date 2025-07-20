package auth

import (
	"context"
	AuthSchema "wentee/blog/app/schema/auth"
)

type IAuthService interface {
	TryLogin(ctx context.Context, loginInfo *AuthSchema.LoginInfo) (tokenString string, err error)
}
