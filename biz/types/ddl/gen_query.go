package main

import (
	basicuser "github.com/xh-polaris/synapse/biz/domain/basicuser/dal/model"
	thirdparty "github.com/xh-polaris/synapse/biz/domain/thirdparty/dal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var path2Model = map[string][]any{
	"/basicuser/dal/query":  {&basicuser.Auth{}, &basicuser.BasicUser{}},
	"/thirdparty/dal/query": {&thirdparty.ThirdPartyUser{}},
}

var root = "../../domain" // domain

func main() {
	// 初始化数据库连接

	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=UTC"
	db, _ := gorm.Open(mysql.Open(dsn))

	for k, v := range path2Model {
		g := gen.NewGenerator(gen.Config{
			OutPath: root + k, // 输出目录
			Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
		})
		g.UseDB(db) // 使用已存在的数据库连接
		g.ApplyBasic(v...)
		g.Execute()
	}
}
