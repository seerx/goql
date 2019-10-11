package main

import (
	"github.com/seerx/goql"
	"github.com/seerx/goql/examples/util"
)

// StudentLoader 学生信息操作承载类
type StudentLoader struct {
	// 定义 api 前缀为 student
	Prefix string `gql:"prefix=student"`
}

// Hello hello
func (StudentLoader) Hello() (string, error) {
	return "Hello Student", nil
}

// SchoolLoader 学校信息操作承载类
type SchoolLoader struct {
	// 定义 api 前缀为 school
	Prefix string `gql:"prefix=school"`
}

// Hello hello
func (SchoolLoader) Hello() (string, error) {
	return "Hello School", nil
}

func init() {
	g := goql.Get()
	g.RegisterQuery(StudentLoader{})
	g.RegisterQuery(SchoolLoader{})
}

func main() {
	util.StartService(8080)
}
