package route

import (
	"dql_admin_backend/middleware"
	"dql_admin_backend/service"
)

func Users() {
	userRouter := Router.Group("/user").Use(middleware.JWTAuth())
	{
		userRouter.GET("/userInfo", service.GetUserInfo)
		userRouter.POST("/logout", service.Logout)
	}
	Router.POST("/user/regist", service.RegistUser)
	Router.POST("/user/login", service.Login)
}
