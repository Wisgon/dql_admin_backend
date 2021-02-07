package services

import (
	"dql_admin_backend/config"
	"dql_admin_backend/model"
	"dql_admin_backend/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRoles(c *gin.Context) {
	var pagination Pagination
	if err := c.ShouldBind(&pagination); err != nil {
		log.Println("get role list bind fail!!!" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"message": "数据错误",
			"code":    config.STATUS["InvalidParam"],
		})
		return
	} else {
		isAdmin := utils.JudgeAuthority(c, "admin")
		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    config.STATUS["AuthForbidden"],
				"message": "only admin can use it.",
			})
			return
		}
		// get list
		roleList, err := model.GetRolesList(pagination.PageSize, pagination.PageNo)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    config.STATUS["InternalError"],
				"message": "get role list error, see logs",
			})
			return
		}
		// fmt.Println("roleList:", roleList)

		c.JSON(http.StatusOK, gin.H{
			"code":    config.STATUS["OK"],
			"message": "get role list success.",
			"data":    roleList.Roles,
		})
	}
}
