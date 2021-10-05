package main

import (
	"log"
	"os"

	"github.com/proelbtn/vnet/cmd/vnet/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "vnet",
		Usage: "Virtual Network Laboratory",
		Commands: []*cli.Command{
			commands.ExecCommand,
			commands.GetCommand,
			commands.StartCommand,
			commands.StopCommand,
			commands.TopologyCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
