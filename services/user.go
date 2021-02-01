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

type RegistForm struct {
	model.User
	PhoneCode string `form:"phoneCode"`
}

func RegistUser(c *gin.Context) {
	var formData RegistForm
	if c.ShouldBind(&formData) == nil {
		// fmt.Println(formData.PhoneCode)
		// fmt.Printf("struct: %+v \n\n", formData)
		newUser := formData.User
		err := newUser.CreateUser()
		if err != nil {
			switch err.Error() {
			case "手机已被注册":
				c.JSON(http.StatusOK, gin.H{
					"message": "手机已被注册",
					"code":    2,
				})
				return
			case "用户名已被注册":
				c.JSON(http.StatusOK, gin.H{
					"message": "用户名已被注册",
					"code":    3,
				})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "注册失败",
					"code":    -1,
				})
				return
			}
		}
	} else {
		log.Println("bind fail!!!")
		c.JSON(http.StatusOK, gin.H{
			"message": "注册失败",
			"code":    1,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"code":    0,
	})
}

func Login(c *gin.Context) {
	var user model.User
	_ = c.ShouldBind(&user)
	result, err := user.VerifyPwd()
	if err != nil {
		log.Println("verifypwd error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"code":    -1,
		})
	}
	if !result {
		c.JSON(http.StatusOK, gin.H{
			"message": "用户名或密码错误",
			"code":    1,
		})
		return
	}
	tokenNext(c, user)

}

// 登录以后签发jwt
func tokenNext(c *gin.Context, user model.User) {
	j := &middleware.JWT{
		SigningKey: []byte(config.JwtSignString), // 唯一签名
	}
	claims := middleware.CustomClaims{
		ID:       user.ID,
		UserName: user.UserName,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,       // 签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*7, // 过期时间 一周
			Issuer:    "sherlock",                     // 签名的发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "获取token失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "login success!",
		"token":   token,
	})
}

// func GetRedisJWT(userName string) (err error, redisJWT string) {
// 	redisJWT, err = model.RedisClient.Get(userName).Result()
// 	return err, redisJWT
// }

func GetUserInfo(c *gin.Context) {
	claims := c.MustGet("claims").(*middleware.CustomClaims) // 获取token解析出来的用户信息
	user := model.User{
		ID: claims.ID,
	}
	err := user.GetUserInfo("id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "get user info error",
		})
		return
	}
	// todo: permissions暂时写死
	permissions := []string{"admin", "editor"}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
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
