package main

import (
	"encoding/json"
	"fmt"

	"github.com/seerx/goql"
	"github.com/seerx/goql/examples/util"
	"github.com/seerx/goql/pkg/require"
)

type Addr struct {
	Address  string `json:"address"`
	PostCode int    `json:"postCode" gql:"limit=10000<$v<=99999"`
}

// Param 参数
type Param struct {
	Name    string  `json:"name" gql:"limit=3<=$v<=20,regexp=\\w+,desc=姓名"`
	Age     int64   `json:"age" gql:"limit=0<=$v<200,desc=年龄"`
	Weight  float32 `json:"weight" gql:"limit=20.0<$v<300,desc=体重(公斤)"`
	Address Addr    `json:"address"`
}

type Loader struct {
	ExplodeParams bool
}

// SubmitInfo resolver 函数
func (Loader) SubmitInfo(p *Param, r *require.Requirement) (string, error) {
	r.Requires("name")

	data, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	return "OK", nil
}

func init() {
	g := goql.Get()
	g.RegisterQuery(Loader{})
}

func main() {
	util.StartService(8080)
}
