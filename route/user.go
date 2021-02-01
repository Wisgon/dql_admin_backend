package route

import (
	"dql_admin_backend/middleware"
	"dql_admin_backend/services"
)

func Users() {
	userRouter := Router.Group("/user").Use(middleware.JWTAuth())
	{
		userRouter.GET("/userInfo", services.GetUserInfo)
		userRouter.POST("/logout", services.Logout)
	}
	Router.POST("/user/regist", services.RegistUser)
	Router.POST("/user/login", services.Login)
}
