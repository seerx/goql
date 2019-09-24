package main

import (
	"github.com/seerx/goql"
	"github.com/seerx/goql/examples/util"
)

func hello() (string, error) {
	return "Hello graphql!", nil
}

//go:generate echo hello
//go:generate go build ./hello.go

func main() {
	g := goql.Get()
	g.RegisterQuery(hello)
	util.StartService(8080)
}
