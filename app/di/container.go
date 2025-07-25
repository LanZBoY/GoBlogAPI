package di

import (
	"wentee/blog/app/appinit"
	"wentee/blog/app/config"
	AuthRouter "wentee/blog/app/handler/auth"
	PostRouter "wentee/blog/app/handler/post"
	UserRouter "wentee/blog/app/handler/user"
	"wentee/blog/app/repo"
	PostRepo "wentee/blog/app/repo/post"
	UserRepo "wentee/blog/app/repo/user"
	AuthSvc "wentee/blog/app/services/auth"
	PostSvc "wentee/blog/app/services/post"
	UserService "wentee/blog/app/services/user"
	"wentee/blog/app/utils"
	"wentee/blog/app/utils/mongo/mongoutils"
)

type Container struct {
	UserRouter *UserRouter.UserRouter
	AuthRouter *AuthRouter.AuthRouter
	PostRouter *PostRouter.PostRouter
}

func InitContainer(appCtx *appinit.AppContext) *Container {
	mainDB := appCtx.MongoClient.Database(config.MOGNO_DATABASE)
	docCounter := &mongoutils.MongoUtils{}
	passwordUtils := &utils.PasswordUtils{}
	objectCreator := &mongoutils.ObjectIdCreator{}

	userRepo := UserRepo.NewUserRepo(&UserRepo.UserCollectionAdapter{Collection: mainDB.Collection(repo.UserCollection)})
	userSvc := UserService.NewUserService(userRepo, passwordUtils, objectCreator)
	userRouter := UserRouter.NewUserRouter(userSvc)

	authSvc := AuthSvc.NewAuthService(userRepo, passwordUtils)
	authRouter := AuthRouter.NewAuthRouter(authSvc)

	poseRepo := PostRepo.NewPostRepo(&PostRepo.PostCollectionAdapter{Collection: mainDB.Collection(repo.PostCollection)}, docCounter)
	postSvc := PostSvc.NewPostService(poseRepo)
	postRouter := PostRouter.NewPostRouter(postSvc)

	return &Container{
		UserRouter: userRouter,
		AuthRouter: authRouter,
		PostRouter: postRouter,
	}
}
