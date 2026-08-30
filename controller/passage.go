package controller

import (
	"github/fetwtgrah/BirdNest/model"

	"github.com/gin-gonic/gin"
)

func AddPost(c *gin.Context) {
	var post model.Passage
	if err := c.ShouldBindJSON(&post); err != nil {

	}
}
