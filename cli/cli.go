package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"

	_ "github.com/blackdreamers/hotshot/cli/new"
	"github.com/blackdreamers/hotshot/cmd"
)

func Run() {
	cmd.Init(
		&cli.App{
			Name:        "hotshot",
			Version:     "0.0.1",
			Usage:       "hotshot is a code initialization cli tool",
			Description: "visit https://github.com/blackdreamers/hotshot for more information",
		},
	)

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}
