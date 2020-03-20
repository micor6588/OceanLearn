package common

import (
	"fmt"
	"micor/ginessential/model"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

// InitDB 开启连接池
func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "371871"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",

		username,
		password,
		host,
		port,
		database,
		charset,
	)
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database err= " + err.Error())
	}
	//自动创建数据表
	db.AutoMigrate(&model.User{})
	return db
}

// GetDB 定义一个方法获取DB实例
func GetDB() *gorm.DB {
	return DB
}
