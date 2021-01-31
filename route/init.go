package route

import (
	"github.com/gin-gonic/gin"
)

var Router = gin.Default()

func init() {
	gin.SetMode(gin.DebugMode)
}
