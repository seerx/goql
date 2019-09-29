package main

import (
	"encoding/json"
	"fmt"

	"github.com/seerx/goql/pkg/param"

	"github.com/seerx/goql"
	"github.com/seerx/goql/examples/util"
)

type Addr struct {
	Address  string `json:"address" gql:"desc=详细地址"`
	PostCode int    `json:"postCode" gql:"limit=10000<=$v<=99999,desc=邮政编码"`
}

// Param 参数
type Param struct {
	Name    string  `json:"name" gql:"limit=3<=$v<=20,regexp=\\w+,error=姓名必须使用英文字母或数字,desc=姓名"`
	Age     int64   `json:"age" gql:"limit=18<=$v,error=年龄必须在 18 周岁以上,desc=年龄"`
	Weight  float32 `json:"weight" gql:"limit=20.0<$v<300,desc=体重(公斤)"`
	Address Addr    `json:"address" gql:"desc=地址"`
}

type Loader struct {
	ExplodeParams bool
}

// SubmitInfo resolver 函数
func (Loader) SubmitInfo(p *Param, r *param.Requirement) (string, error) {
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
