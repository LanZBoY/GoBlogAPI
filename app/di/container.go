package di

import (
	"wentee/blog/app/appinit"
	"wentee/blog/app/config"
	UserRouter "wentee/blog/app/handler/user"
	"wentee/blog/app/model/mongodb"
	UserRepo "wentee/blog/app/repo/user"
	UserService "wentee/blog/app/service/user"
)

type Container struct {
	UserRouter *UserRouter.UserRouter
}

func InitContainer(appCtx *appinit.AppContext) *Container {
	mainDB := appCtx.MongoClient.Database(config.MOGNO_DATABASE)

	userRepo := UserRepo.NewUserRepo(mainDB.Collection(mongodb.UserCollection))
	userSvc := UserService.NewUserService(userRepo)
	userRouter := UserRouter.NewUserRouter(userSvc)

	return &Container{
		UserRouter: userRouter,
	}
}
