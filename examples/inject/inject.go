package main

import (
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

func main() {
	g := goql.Get()
	g.RegisterQuery(&Loader{})
	g.RegisterQuery(Inject)

	util.StartService(8080)
}
