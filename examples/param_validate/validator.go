package main

import (
	"github.com/seerx/goql"
	"github.com/seerx/goql/examples/util"
	"github.com/seerx/goql/pkg/require"
)

// Param 参数
type Param struct {
	Name   string  `json:"name" gql:"limit=3<=$v<=20,regexp=\\w+,desc=姓名"`
	Age    int     `json:"age" gql:"limit=0<=$v<200,desc=年龄"`
	Weight float64 `json:"weight" gql:"limit=20.0<$v<300,desc=体重(公斤)"`
}

type Loader struct {
	ExplodeParams bool
}

// SubmitInfo resolver 函数
func (Loader) SubmitInfo(p *Param, r *require.Requirement) (string, error) {
	r.Requires("name")
	return "OK", nil
}

func init() {
	g := goql.Get()
	g.RegisterQuery(Loader{})
}

func main() {
	util.StartService(8080)
}
