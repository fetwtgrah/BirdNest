package controller

import (
	"github/fetwtgrah/BirdNest/configs"
	"github/fetwtgrah/BirdNest/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddPassageRequest struct {
	Content string `json:"content" binding:"required"`
}

func AddPassage(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	user_name, _ := c.Get("user_name")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "未登录",
		})
		return
	}
	var req AddPassageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	passage := model.Passage{
		UserID:  user_id.(int64),
		Content: req.Content,
	}
	if err := configs.Db.Create(&passage).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":     "上传失败",
			"err":     err.Error(),
			"当前请求的id": passage.UserID,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":    "上传成功",
		"author": user_name,
	})
}
