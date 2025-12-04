package main

import (
	basicuser "github.com/xh-polaris/synapse/biz/domain/basicuser/dal/model"
	thirdparty "github.com/xh-polaris/synapse/biz/domain/thirdparty/dal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 从go的结构生成数据库表

func main() {
	// MySQL 连接配置
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=UTC"

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&basicuser.Auth{}, &basicuser.BasicUser{}, &basicuser.School{}, &thirdparty.ThirdPartyUser{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}
}