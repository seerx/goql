package main

import (
	"context"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/seerx/goql"
)

func init() {
	g := goql.Get()
	g.RegisterInject(InjectClass) // 注册注入函数
}

// 要注入的对象
type ClassInfo struct {
	Grade string `json:"grade"` // 年级
	Class string `json:"class"` // 班级
}

// 注入函数
func InjectClass(ctx context.Context, r *http.Request, gp *graphql.ResolveParams) *ClassInfo {
	return &ClassInfo{
		Grade: "一年级",
		Class: "1 班",
	}
}
