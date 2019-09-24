package main

import (
	"errors"

	"github.com/seerx/goql"
	"github.com/seerx/goql/examples/util"
)

// StudentRequest 查询参数
type StudentRequest struct {
	ID int
}

// Student 返回的数据结构
type Student struct {
	ID    int
	Name  string
	Class string
}

// StudentByID 功能函数
func StudentByID(req *StudentRequest) (*Student, error) {
	for _, s := range students {
		if s.ID == req.ID {
			return s, nil
		}
	}
	return nil, errors.New("没有找到")
}

// AllStudents 获取学生列表
func AllStudents() ([]*Student, error) {
	return students, nil
}

func main() {
	g := goql.Get()
	g.RegisterQuery(StudentByID)
	g.RegisterQuery(AllStudents)
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
