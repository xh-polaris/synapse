package main

import (
	"path/filepath"
	"runtime"

	basicuser "github.com/xh-polaris/synapse/biz/domain/basicuser/dal/model"
	thirdparty "github.com/xh-polaris/synapse/biz/domain/thirdparty/dal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var path2Model = map[string][]any{
	"/basicuser/dal/query":  {&basicuser.Auth{}, &basicuser.BasicUser{}, &basicuser.School{}},
	"/thirdparty/dal/query": {&thirdparty.ThirdPartyUser{}},
}

func getRootPath() string {
	_, filename, _, _ := runtime.Caller(0)
	// 获取当前文件目录
	currentDir := filepath.Dir(filename)
	// 返回到domain目录的相对路径
	return filepath.Join(currentDir, "../../domain")
}

func main() {
	// 初始化数据库连接

	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=UTC"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	root := getRootPath()
	
	for k, v := range path2Model {
		g := gen.NewGenerator(gen.Config{
			OutPath: filepath.Join(root, k), // 输出目录
			Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
		})
		g.UseDB(db) // 使用已存在的数据库连接
		g.ApplyBasic(v...)
		g.Execute()
	}
}