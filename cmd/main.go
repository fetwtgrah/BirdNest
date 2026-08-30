package main

import (
	"fmt"
	"github/fetwtgrah/BirdNest/configs"
	"github/fetwtgrah/BirdNest/controller"
	"github/fetwtgrah/BirdNest/midware"
	"github/fetwtgrah/BirdNest/model"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := configs.InitConfig(); err != nil {
		panic("系统配置出错" + err.Error())
	}
	configs.InitDb()
	configs.InitRedis()
	err := configs.Db.AutoMigrate(&model.User{}, &model.Passage{})
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
		passage.POST("", controller.AddPost)
	}
	
	if err := r.Run(":8080"); err != nil {
		panic("路由连接出错" + err.Error())
	}
}
