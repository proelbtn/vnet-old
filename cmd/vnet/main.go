package main

import (
	"fmt"
	"log"
	"os"

	"github.com/proelbtn/vnet/cmd/vnet/commands"
	"github.com/urfave/cli/v2"
)

var (
	version  string = "vX.X.X"
	revision string = "xxxxxxx"
)

func main() {
	app := &cli.App{
		Name:    "vnet",
		Usage:   "Virtual Network Laboratory",
		Version: fmt.Sprintf("%s (rev: %s)", version, revision),
		Commands: []*cli.Command{
			commands.StartCommand,
			commands.StopCommand,
			commands.RestartCommand,
			commands.ExecCommand,
			commands.ShowCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
