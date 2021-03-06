package services

import (
	"dql_admin_backend/config"
	"dql_admin_backend/middleware"
	"dql_admin_backend/model"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func RegistUser(c *gin.Context) {
	var formData RegistForm
	if err := c.ShouldBind(&formData); err != nil {
		log.Println("regist bind fail!!!" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"message": "注册失败",
			"code":    config.STATUS["InvalidParam"],
		})
		return
	} else {
		newUser := formData.User
		err := newUser.CreateUser()
		if err != nil {
			switch err.Error() {
			case "手机已被注册":
				c.JSON(http.StatusOK, gin.H{
					"message": "手机已被注册",
					"code":    config.STATUS["InvalidParam"],
				})
				return
			case "用户名已被注册":
				c.JSON(http.StatusOK, gin.H{
					"message": "用户名已被注册",
					"code":    config.STATUS["InvalidParam"],
				})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "InternalError, see logs",
					"code":    config.STATUS["InternalError"],
				})
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"code":    config.STATUS["OK"],
	})
}

func Login(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		log.Println("login bind fail!!!" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"message": "数据错误",
			"code":    config.STATUS["InvalidParam"],
		})
		return
	} else {
		result, err := user.VerifyPwd()
		if err != nil {
			if err.Error() == "user not found!" {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    config.STATUS["NotFound"],
					"message": "用户名或密码错误",
				})
			} else {
				log.Println("verifypwd error: " + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    config.STATUS["InternalError"],
					"message": "verify password error, see logs",
				})
			}
			return
		}
		if !result {
			c.JSON(http.StatusOK, gin.H{
				"message": "用户名或密码错误",
				"code":    config.STATUS["InvalidParam"],
			})
			return
		}
		tokenNext(c, user)
	}

}

// 登录以后签发jwt
func tokenNext(c *gin.Context, user model.User) {
	j := &middleware.JWT{
		SigningKey: []byte(config.JwtSignString), // 唯一签名
	}
	roleString := ""
	for _, role := range user.Roles {
		roleString += role.RoleID + "#" // 每个role用#号分隔
	}
	sc, err := model.GetSystemConfig()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    config.STATUS["InternalError"],
			"message": "获取system config失败，see logs",
		})
		return
	}
	claims := middleware.CustomClaims{
		ID:                user.UID,
		Roles:             roleString,
		PermissionVersion: sc.SystemConfigs[0].PermissionVersion,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,       // 签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*7, // 过期时间 一周
			Issuer:    "sherlock",                     // 签名的发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    config.STATUS["NotFound"],
			"message": "获取token失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    config.STATUS["OK"],
		"message": "login success!",
		"data": gin.H{
			"accessToken": token,
		},
	})
}

// func GetRedisJWT(userName string) (err error, redisJWT string) {
// 	redisJWT, err = model.RedisClient.Get(userName).Result()
// 	return err, redisJWT
// }

func GetUserInfo(c *gin.Context) {
	claims := c.MustGet("claims").(*middleware.CustomClaims) // 获取token解析出来的用户信息
	user := model.User{
		UID: claims.ID,
	}
	err := user.GetUserInfo("id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    config.STATUS["InternalError"],
			"message": "get user info error, see logs",
		})
		return
	}
	permissions := []string{}
	for _, role := range user.Roles {
		permissions = append(permissions, role.RoleID)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    config.STATUS["OK"],
		"message": "get user info ok!",
		"data": gin.H{
			"avatar":      user.Avatar,
			"username":    user.UserName,
			"permissions": permissions,
		},
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

func GetUserList(c *gin.Context) {
	var searchQuery SearchQuery
	if err := c.ShouldBind(&searchQuery); err != nil {
		log.Println("get User list bind fail!!!" + err.Error())
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
		userList, err := model.GetUserList(searchQuery.PageSize, searchQuery.PageNo, searchQuery.Username, searchQuery.Fuzz)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    config.STATUS["InternalError"],
				"message": "get user list error, see logs",
			})
			return
		}
		// fmt.Println("userList:", userList.Users)
		c.JSON(http.StatusOK, gin.H{
			"code":    config.STATUS["OK"],
			"message": "get user list success.",
			"data":    userList.Users,
		})
	}
}

func UpdateUser(c *gin.Context) {
	isAdmin := JudgeAuthority(c, "admin")
	if !isAdmin {
		return
	}

	// res, err := c.GetRawData()
	// if err != nil {
	// 	log.Println("update user get raw error:" + err.Error())
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": "InternalError, see logs",
	// 		"code":    config.STATUS["InternalError"],
	// 	})
	// 	return
	// }
	//fmt.Println("res:", string(res))

	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		log.Println("login bind fail!!!" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"message": "数据错误",
			"code":    config.STATUS["InvalidParam"],
		})
		return
	} else {
		err = model.UpdateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "InternalError, see logs",
				"code":    config.STATUS["InternalError"],
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    config.STATUS["OK"],
			"message": "成功更新用户数据",
		})
	}

}

func DeleteUser(c *gin.Context) {
	isAdmin := JudgeAuthority(c, "admin")
	if !isAdmin {
		return
	}

	user := model.User{}
	if err := c.ShouldBind(&user); err != nil {
		log.Println("do delete bind fail!!!" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"message": "输入的数据字段不对",
			"code":    config.STATUS["InvalidParam"],
		})
		return
	} else {
		err := model.DeleteUser(user)
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
