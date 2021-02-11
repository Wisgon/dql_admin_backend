package services

import (
	"dql_admin_backend/config"
	"dql_admin_backend/middleware"
	"dql_admin_backend/model"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateRole(c *gin.Context) {
	isAdmin := JudgeAuthority(c, "admin")
	if !isAdmin {
		return
	}

	role := model.Role{}
	if err := c.ShouldBind(&role); err != nil {
		log.Println("do edit bind fail!!!" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"message": "输入的数据字段不对",
			"code":    config.STATUS["InvalidParam"],
		})
		return
	} else {
		err := model.CreateRole(role)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "创建失败，请查看后台日志",
				"code":    config.STATUS["InternalError"],
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    config.STATUS["OK"],
		"message": "创建成功",
	})
}

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
		isAdmin := JudgeAuthority(c, "admin")
		if !isAdmin {
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

func GetAccessablePages(c *gin.Context) {
	claims := c.MustGet("claims").(*middleware.CustomClaims)
	var roles []string
	// 如果是admin，则还可以传一个role的参数来查询
	if strings.Contains(claims.Roles, "admin") {
		roleId := c.Query("role_id")
		if roleId != "" {
			roles = []string{roleId}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    config.STATUS["OK"],
				"message": "you are admin.",
				"data":    make(map[string]map[string]bool),
			})
			return
		}
	} else {
		rolesStrings := strings.Split(claims.Roles, "#")
		roles = rolesStrings[:len(rolesStrings)-1] // 去掉最后一个，因为最后一个为空

	}

	accessablePages, err := model.GetAccessablePages(roles)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    config.STATUS["InternalError"],
			"message": "get accessablePages error, see logs",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    config.STATUS["OK"],
		"message": "get accessablePages success.",
		"data":    accessablePages,
	})
}

func EditRole(c *gin.Context) {
	isAdmin := JudgeAuthority(c, "admin")
	if !isAdmin {
		return
	}

	role := model.Role{}
	if err := c.ShouldBind(&role); err != nil {
		log.Println("do edit bind fail!!!" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"message": "输入的数据字段不对",
			"code":    config.STATUS["InvalidParam"],
		})
		return
	} else {
		err := model.EditRole(role)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "保存失败，请查看后台日志",
				"code":    config.STATUS["InternalError"],
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    config.STATUS["OK"],
		"message": "保存成功",
	})
}

func DeleteRole(c *gin.Context) {
	isAdmin := JudgeAuthority(c, "admin")
	if !isAdmin {
		return
	}

	role := model.Role{}
	if err := c.ShouldBind(&role); err != nil {
		log.Println("do delete bind fail!!!" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"message": "输入的数据字段不对",
			"code":    config.STATUS["InvalidParam"],
		})
		return
	} else {
		err := model.DeleteRole(role)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "删除失败，请查看后台日志",
				"code":    config.STATUS["InternalError"],
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    config.STATUS["OK"],
		"message": "删除成功",
	})
}
