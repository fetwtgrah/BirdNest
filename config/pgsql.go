package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() {
	dsn := "host=localhost user=postgres port=5432 password=123456 dbname=postgres sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败" + err.Error())
	}
	Db = db
}
