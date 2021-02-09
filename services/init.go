package services

import (
	"dql_admin_backend/middleware"
	"dql_admin_backend/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RegistForm struct {
	model.User
	PhoneCode string `form:"phoneCode"` // 验证码
}

type Pagination struct {
	PageSize int `json:"pageSize"`
	PageNo   int `json:"pageNo"`
}

func JudgeAuthority(c *gin.Context, role string) bool {
	claims := c.MustGet("claims").(*middleware.CustomClaims) // 获取token解析出来的用户信息
	if !strings.Contains(claims.Roles, role) {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    1,
			"message": "只有" + role + "才有权限使用此接口",
		})
		return false
	}
	return true
}
