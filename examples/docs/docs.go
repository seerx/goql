package main

import (
	"errors"

	"github.com/seerx/goql"
	"github.com/seerx/goql/examples/util"
)

// Student 学生信息
type Student struct {
	ID    int    `json:"id" gql:"desc=学生编号"`
	Name  string `json:"name" gql:"desc=姓名"`
	Class string `json:"class" gql:"desc=班级"`
}

// StudentRequest 查询参数
type StudentRequest struct {
	ID int `json:"id" gql:"desc=学生编号"`
}

// Loader resolver 函数承载结构
type Loader struct {
	ExplodeParams int // 只要定义了 ExplodeParams 字段(类型不限)，那么该结构的所有 Resolver 函数的跟参数 in 将会去掉
}

// QueryDesc Query 函数的说明信息
func (l *Loader) QueryDesc() string {
	return `查询学生信息`
}

// Query resolver 函数
func (l *Loader) Query(req *StudentRequest) (*Student, error) {
	for _, s := range students {
		if s.ID == req.ID {
			return s, nil
		}
	}
	return nil, errors.New("没有找到")
}

func main() {
	g := goql.Get()
	g.RegisterQuery(&Loader{})
	util.StartService(8080)
}

var students = []*Student{{
	ID:    1,
	Name:  "小明",
	Class: "1(1)",
}, {
	ID:    2,
	Name:  "小红",
	Class: "1(1)",
}, {
	ID:    3,
	Name:  "小玲",
	Class: "2(1)",
}, {
	ID:    1,
	Name:  "小飞",
	Class: "2(1)",
}}
