package main

import (
	"github.com/seerx/goql"
	"github.com/seerx/goql/examples/util"
)

func main() {
	g := goql.Get()
	g.RegisterQuery(func() (string, error) {
		return "Hello goql!", nil
	})
	util.StartService(8080)
}
