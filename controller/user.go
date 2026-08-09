package controller

import (
	"github/fetwtgrah/BirdNest/config"
	"github/fetwtgrah/BirdNest/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "msg": "数据绑定错误"})
	}
	config.Db.Create(&user)
}
