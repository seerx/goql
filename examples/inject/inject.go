package main

import (
	"fmt"

	"github.com/seerx/goql"
	"github.com/seerx/goql/examples/util"
)

// Loader resolver 函数承载结构
type Loader struct {
	Class *ClassInfo
}

// InjectToLoader 注入 ClassInfo 到 Loader.Class
func (l Loader) InjectToLoader() (*ClassInfo, error) {
	return l.Class, nil
}

// Inject 注入 ClassInfo 到参数 class
func Inject(class *ClassInfo) (*ClassInfo, error) {
	return class, nil
}

// ReadFromDB 测试自动关闭功能
func ReadFromDB(db *DBConnection) (string, error) {
	fmt.Println("使用数据库连接")
	return db.DB, nil
}

func main() {
	g := goql.Get()
	g.RegisterQuery(&Loader{})
	g.RegisterQuery(Inject)

	g.RegisterQuery(ReadFromDB)

	util.StartService(8080)
}
