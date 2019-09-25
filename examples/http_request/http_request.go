package main

import (
	"context"
	"net/http"

	"github.com/seerx/goql/examples/util"

	"github.com/seerx/goql"

	"github.com/graphql-go/graphql"
)

func init() {
	g := goql.Get()
	g.RegisterInject(injectRequest)
	g.RegisterQuery(showMeMyAddress)
}

// injectRequest  http.Request 注入函数
func injectRequest(ctx context.Context, r *http.Request, gp *graphql.ResolveParams) *http.Request {
	return r
}

// showMeMyAddress resolver 函数
func showMeMyAddress(r *http.Request) (string, error) {
	return r.RemoteAddr, nil
}

func main() {
	util.StartService(8080)
}
