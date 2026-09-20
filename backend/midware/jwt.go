package midware

import (
	"context"
	"fmt"
	"github/fetwtgrah/BirdNest/backend/configs"
	"github/fetwtgrah/BirdNest/backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckToken(c *gin.Context) {
	token := c.GetHeader("token")
	claim, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		c.Abort()
		return
	}
	key := fmt.Sprintf("user:name:%s", claim.Username)
	RcToken := configs.Rc.Get(context.Background(), key).Val()
	if token != RcToken {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "权限错误",
		})
		c.Abort()
		return
	}
	c.Set("user_id", claim.UserID)
	c.Set("user_name", claim.Username)
	c.Next()
}
