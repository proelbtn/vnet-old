package commands

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

var GetCommand = &cli.Command{
	Name:  "get",
	Usage: "Get information in laboratory",
	Flags: append([]cli.Flag{}, commonFlags...),
	Subcommands: []*cli.Command{
		GetPortNameCommand,
	},
}

var GetPortNameCommand = &cli.Command{
	Name:      "port",
	Usage:     "Get port name in container",
	ArgsUsage: "CONTAINER PORT",
	Action:    getPortName,
}

func getPortName(c *cli.Context) error {
	lab, err := initialize(c)
	if err != nil {
		return err
	}

	usecase, err := getUsecase()
	if err != nil {
		return err
	}

	req, err := lab.ToWritableLaboratory()
	if err != nil {
		return err
	}

	if c.Args().Len() != 2 {
		return errors.New("container name and port name are needed")
	}

	containerName := c.Args().Get(0)
	portName := c.Args().Get(1)

	name, err := usecase.GetPortName(req, containerName, portName)
	if err != nil {
		return err
	}

	fmt.Println(name)
	return nil
}
