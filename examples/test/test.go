package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/seerx/goql"

	"github.com/seerx/goql/pkg/require"

	"github.com/graphql-go/graphql"

	"github.com/graphql-go/handler"
)

type TT struct {
	Prefix        string `gql:"prefix=ttgo"`
	ExplodeParams bool
	Require       require.Requirement
	Aa            *Abd
}

func aa() (string, error) {
	return "", nil
}

func (t TT) Atta(r require.Requirement) (int, error) {
	t.Require.Requires("a")
	r.Requires("a")
	return 1, nil
}

func (t TT) Ttta() (int, error) {
	return 1, nil
}

func (t TT) Ttta1() error {
	return nil
}

func Ttta(p *Abc, r *require.Requirement) (int, error) {
	r.Requires("newa", "IO.A")
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(data))
	return 1, nil
}

func (t TT) Nb(p *Abc) (int, error) {
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(data))
	return 1000, nil
}

type AryItem struct {
	A string
	C int `gql:"limit=0<$v<10"`
}

type Abc struct {
	UI  string     `json:"newa" gql:"desc=字符串"`
	Ary []*AryItem `json:"ary"`
	A   *AryItem
	IO  AryItem
	B   string `json:"-"`
}

type Abd struct {
	A string `json:"newa" gql:"desc=1235678"`
	B string `json:"b"`
}

func (t TT) TttaDesc() string {
	return "123"
}

//var mp = map[reflect.Type]string{}
//
//func printType(typ reflect.Type) {
//	td := reflects.ParseType(typ)
//
//	if _, ok := mp[typ]; ok {
//		fmt.Println("已经存在" + typ.Name())
//	} else {
//		mp[typ] = ""
//	}
//	fmt.Println(td.Summary())
//	fmt.Println("-------------------------------")
//}
const port = 8080

func main() {
	g := goql.Get()
	//g.RegisterQuery(aa)
	g.RegisterQuery(Ttta)
	g.RegisterQuery(TT{})

	//var inj gqlx.InjectTemplate
	inj := func(ctx context.Context, r *http.Request, gp *graphql.ResolveParams) *Abd {
		return &Abd{
			A: "123456789",
		}
	}
	g.RegisterInject(inj)

	http.Handle("/", g.CreateHandler(&handler.Config{
		Pretty:   true,
		GraphiQL: true,
	}))

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	//var a int = 10
	//printType(reflect.TypeOf(a))
	//
	//ary := []int{1, 2}
	//printType(reflect.TypeOf(ary))
	//
	//printType(reflect.TypeOf([]*TT{}))
	//printType(reflect.TypeOf([]TT{}))
	//printType(reflect.TypeOf([]TT{}))
	//printType(reflect.TypeOf(TT{}))
	//printType(reflect.TypeOf(&TT{}))

	//gql := gqlx.NewGQL(nil)
	//gql.RegisterQuery(aa)
	//gql.RegisterQuery(Ttta)
	//gql.RegisterQuery(&TT{})

}
