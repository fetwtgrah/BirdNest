package main

import (
	"fmt"
	"github/fetwtgrah/BirdNest/backend/configs"
	"github/fetwtgrah/BirdNest/backend/controller"
	"github/fetwtgrah/BirdNest/backend/midware"
	"github/fetwtgrah/BirdNest/backend/model"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := configs.InitConfig(); err != nil {
		panic("系统配置出错" + err.Error())
	}
	configs.InitDb()
	configs.InitRedis()
	err := configs.Db.AutoMigrate(&model.User{}, &model.Passage{}, &model.Tag{})
	if err != nil {
		fmt.Println("数据库连接错误" + err.Error())
	}
	r := gin.Default()
	v1 := r.Group("/api/v1")

	v1.POST("/code", controller.SendCode)

	user := v1.Group("/user")
	{
		user.POST("", controller.RegisterByCode)
		user.POST("/login", controller.LoginByPassword)
	}
	passage := v1.Group("/passage")
	passage.Use(midware.CheckToken)
	{
		passage.POST("", controller.AddPassage)
		passage.GET("/:id", controller.GetPassage)
		passage.GET("/all", controller.GetAllPassage)
		passage.PUT("/:id", controller.UpdatePassage)
		passage.DELETE("/:id", controller.DeletePassage)

		passage.GET("/tag/:tag", controller.GetPassageByTag)
	}

	if err := r.Run(":8080"); err != nil {
		panic("路由连接出错" + err.Error())
	}
}
