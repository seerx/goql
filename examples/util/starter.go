package util

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/graphql-go/handler"
	"github.com/seerx/goql"
)

func init() {
	goql.Configure(&goql.SilentLogger{})
}

func StartService(port int) {
	g := goql.Get()
	handle := g.CreateHandler(&handler.Config{
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/", handle)
	fmt.Println("Service address: ")
	fmt.Printf("http://localhost:%d/", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

var commands = map[string]string{
	"windows": "cmd /c start",
	"darwin":  "open",
	"linux":   "xdg-open",
}

// OpenURL calls the OS default program for uri
func OpenURL(uri string) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	cmd := exec.Command(run, uri)
	return cmd.Start()
}
