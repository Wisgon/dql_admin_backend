package utils

import (
	"dql_admin_backend/middleware"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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
