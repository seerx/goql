package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"

	"github.com/seerx/goql"
)

func init() {
	g := goql.Get()
	g.RegisterInject(injectCloseAble)
}

// DBConnection 运行后可以自动清理
type DBConnection struct {
	DB string
}

// Close 自动执行
func (conn *DBConnection) Close() error {
	conn.DB = "已关闭数据库连接"
	fmt.Println("清理工作", "关闭数据库连接")
	return nil
}

func injectCloseAble(ctx context.Context, r *http.Request, gp *graphql.ResolveParams) *DBConnection {
	// 建立数据库连接
	fmt.Println("准备工作", "建立数据库连接")
	return &DBConnection{
		DB: "已连接数据",
	}
}
