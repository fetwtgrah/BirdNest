package midware

import (
	"context"
	"fmt"
	"github/fetwtgrah/BirdNest/configs"
	"github/fetwtgrah/BirdNest/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckToken(c *gin.Context) {
	token := c.GetHeader("token")
	claim, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "权限错误",
		})
		c.Abort()
		return
	}
	key := fmt.Sprintf("user:name:%s", claim["username"])
	RcToken := configs.Rc.Get(context.Background(), key).Val()
	if token != RcToken {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "权限错误",
		})
		c.Abort()
		return
	}
	userIDFloat, ok := claim["userid"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "token 内容异常"})
		c.Abort()
		return
	}
	c.Set("user_id", int64(userIDFloat))
	c.Next()
}
