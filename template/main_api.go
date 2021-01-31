package template

var (
	MainAPI = `package main

import (
	"{{.Dir}}/handler"
	"github.com/blackdreamers/core/server"
)

func main() {
	// Init server
	server.Init(
		server.Name({{lower .Alias}}),
		server.Type(server.API),
		server.EnableDB(false),
	)

	// Register handles
	server.Handles(
		new(handler.{{title .Alias}}),
	)

	// Run server
	server.Run()
}
`
)
