package template

var (
	MainSRV = `package main

import (
	coredb "github.com/blackdreamers/core/db"
	"github.com/blackdreamers/core/server"
	"{{.Dir}}/handler"
	"{{.Dir}}/subscriber"
)

func main() {
	// DB repositories
	coredb.Repositories()

	// Register handles
	server.Handles(
		new(handler.{{title .Alias}}),
	)

	// Register subscribers
	server.Subscribers(
		new(subscriber.{{title .Alias}}),
	)

	// Run server
	server.Run()
}
`
)
