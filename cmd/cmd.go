package cmd

import (
	"os"
	"sort"

	"github.com/urfave/cli/v2"
)

var (
	defaultCmd *cli.App
	commands   []*cli.Command
)

func Init(app *cli.App) {
	defaultCmd = app
}

func Register(cmds ...*cli.Command) {
	commands = append(commands, cmds...)
}

func Run() error {
	defaultCmd.Commands = commands

	sort.Sort(cli.FlagsByName(defaultCmd.Flags))
	sort.Sort(cli.CommandsByName(defaultCmd.Commands))

	return defaultCmd.Run(os.Args)
}
