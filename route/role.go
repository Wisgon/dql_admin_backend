package route

import (
	"dql_admin_backend/middleware"
	"dql_admin_backend/services"
)

func Roles() {
	roleRouter := Router.Group("/roleManagement").Use(middleware.JWTAuth())
	{
		roleRouter.POST("/getList", services.GetRoles)
	}
}
