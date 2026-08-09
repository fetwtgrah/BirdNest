package main

import (
	"fmt"
	"github/fetwtgrah/BirdNest/config"
	"github/fetwtgrah/BirdNest/controller"
	"github/fetwtgrah/BirdNest/model"

	"github.com/gin-gonic/gin"
)

// 用户系统：
// register,login
func main() {
	config.InitDb()
	config.InitRedis()
	err := config.Db.AutoMigrate(&model.User{})
	if err != nil {
		fmt.Println("数据库连接错误" + err.Error())
	}
	r := gin.Default()
	r.POST("/user", controller.Register)
	if err := r.Run(":8080"); err != nil {
		panic("路由连接出错" + err.Error())
	}
}
