package controller

import (
	"fmt"
	"github/fetwtgrah/BirdNest/backend/configs"
	"github/fetwtgrah/BirdNest/backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AddPassageRequest struct {
	Tags    []model.Tag `json:"tags"`
	Content string      `json:"content" binding:"required"`
}

func AddPassage(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}
	var user model.User
	if err := configs.Db.First(&user, user_id).Error; err != nil {
		c.JSON(400, gin.H{
			"msg": "该用户不存在",
		})
		return
	}
	user_name, _ := c.Get("user_name")

	var req AddPassageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	//err := configs.Db.Model(&user).Association("Passages").Append(model.Passage{Tags: req.Tags, Content: req.Content})
	passage := model.Passage{
		UserID:  user_id.(uint),
		Content: req.Content,
		Tags:    req.Tags,
	}
	err := configs.Db.Create(&passage).Error
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "上传失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":    "上传成功",
		"author": user_name,
		"tags":   req.Tags,
	})
}

func GetPassage(c *gin.Context) {
	id := c.Param("id")
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "用户信息获取失败"})
		return
	}
	var passage model.Passage
	if err := configs.Db.Preload("Tags").First(&passage, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "不存在该文章或文章id错误",
		})
		return
	}
	if passage.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"msg": "无权限查看别人的文章"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "文章获取成功",
		"data": passage,
	})
}

func GetAllPassage(c *gin.Context) {
	var passages []model.Passage
	user_id, ok := c.Get("user_id")
	user_name, _ := c.Get("user_name")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "用户id获取出错",
		})
		return
	}
	userID := user_id.(uint)
	err := configs.Db.Preload("Tags").Where("user_id = ?", userID).Find(&passages).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	success := fmt.Sprintf("成功获取用户：%v的全部文章,一共%v篇", user_name, len(passages))
	c.JSON(http.StatusOK, gin.H{
		"msg":  success,
		"data": passages,
	})
}

func UpdatePassage(c *gin.Context) {
	id := c.Param("id")
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "用户的id获取失败",
		})
		return
	}
	var passage model.Passage
	if err := configs.Db.First(&passage, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "文章不存在",
		})
		return
	}
	if passage.UserID != user_id.(uint) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无权限修改别人的文章",
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
	passage.Content = req.Content
	passage.Tags = req.Tags
	if err := configs.Db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&passage).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "更新失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新成功",
		"data": passage,
	})
}

func DeletePassage(c *gin.Context) {
	id := c.Param("id")
	user_id, _ := c.Get("user_id")
	var passage model.Passage
	if err := configs.Db.First(&passage, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "文章不存在",
		})
		return
	}
	if passage.UserID != user_id.(uint) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无权限修改别人的文章",
		})
		return
	}
	if err := configs.Db.Delete(&passage).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "删除失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func GetPassageByTag(c *gin.Context) {
	tagName := c.Param("tag")
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "用户信息获取失败"})
		return
	}

	var passages []model.Passage
	if err := configs.Db.
		Preload("Tags").
		Joins("JOIN passage_tags ON passage_tags.passage_id = passages.id").
		Joins("JOIN tags ON tags.id = passage_tags.tag_id").
		Where("passages.user_id = ? AND tags.tag_name = ?", userID.(uint), tagName).
		Find(&passages).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	res := fmt.Sprintf("%v的文章一共有%v篇", tagName, len(passages))
	c.JSON(http.StatusOK, gin.H{
		"msg":  res,
		"data": passages,
	})
}
