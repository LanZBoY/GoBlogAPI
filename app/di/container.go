package di

import (
	"wentee/blog/app/appinit"
	"wentee/blog/app/config"
	AuthRouter "wentee/blog/app/handler/auth"
	UserRouter "wentee/blog/app/handler/user"
	"wentee/blog/app/model/mongodb"
	UserRepo "wentee/blog/app/repo/user"
	AuthSvc "wentee/blog/app/services/auth"
	UserService "wentee/blog/app/services/user"
)

type Container struct {
	UserRouter *UserRouter.UserRouter
	AuthRouter *AuthRouter.AuthRouter
}

func InitContainer(appCtx *appinit.AppContext) *Container {
	mainDB := appCtx.MongoClient.Database(config.MOGNO_DATABASE)

	userRepo := UserRepo.NewUserRepo(mainDB.Collection(mongodb.UserCollection))
	userSvc := UserService.NewUserService(userRepo)
	userRouter := UserRouter.NewUserRouter(userSvc)

	authSvc := AuthSvc.NewAuthService(userRepo)
	authRouter := AuthRouter.NewAuthRouter(authSvc)

	return &Container{
		UserRouter: userRouter,
		AuthRouter: authRouter,
	}
}
