package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/handler"
	"github.com/seerx/goql"

	"github.com/julienschmidt/httprouter"
)

func Hello() (string, error) {
	return "Hello httprouter", nil
}

func init() {
	g := goql.Get()
	g.RegisterQuery(Hello)
}

func main() {
	router := httprouter.New()
	svr := &http.Server{Addr: fmt.Sprintf(":%d", 8080)}
	svr.Handler = router

	g := goql.Get()
	handle := g.CreateHandler(&handler.Config{
		Pretty:   true,
		GraphiQL: true,
	})

	router.GET("/graphql", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		handle.ServeHTTP(writer, request)
	})
	router.POST("/graphql", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		handle.ServeHTTP(writer, request)
	})

	log.Fatal(svr.ListenAndServe())
}
