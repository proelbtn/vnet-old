package commands

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

var GetCommand = &cli.Command{
	Name:  "get",
	Usage: "Get information in laboratory",
	Subcommands: []*cli.Command{
		GetPortNameCommand,
	},
}

var GetPortNameCommand = &cli.Command{
	Name:      "port",
	Usage:     "Get port name in container",
	Flags:     append([]cli.Flag{}, commonFlags...),
	ArgsUsage: "CONTAINER PORT",
	Action:    getPortName,
}

func getPortName(c *cli.Context) error {
	lab, err := initialize(c)
	if err != nil {
		return err
	}

	usecase, err := newUsecase(WithMockContainerManager)
	if err != nil {
		return err
	}

	if c.Args().Len() != 2 {
		return errors.New("container name and port name are needed")
	}

	containerName := c.Args().Get(0)
	portName := c.Args().Get(1)

	name := usecase.GetPortName(lab.Name, containerName, portName)

	fmt.Println(name)
	return nil
}
