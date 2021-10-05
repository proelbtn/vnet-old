package commands

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

var ShowCommand = &cli.Command{
	Name:  "show",
	Usage: "Show information about laboratory",
	Subcommands: []*cli.Command{
		ShowPortNameCommand,
		ShowTopologyCommand,
	},
}

var ShowPortNameCommand = &cli.Command{
	Name:      "port",
	Usage:     "Get port name in container",
	Flags:     append([]cli.Flag{}, commonFlags...),
	ArgsUsage: "CONTAINER PORT",
	Action:    getPortName,
}

var ShowTopologyCommand = &cli.Command{
	Name:   "topology",
	Usage:  "Show topology of laboratory",
	Flags:  append([]cli.Flag{}, commonFlags...),
	Action: topology,
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
		return errors.New("container name and port name are required")
	}

	containerName := c.Args().Get(0)
	portName := c.Args().Get(1)

	name := usecase.GetPortName(lab.Name, containerName, portName)

	fmt.Println(name)
	return nil
}

func topology(c *cli.Context) error {
	lab, err := initialize(c)
	if err != nil {
		return err
	}

	usecase, err := newUsecase(WithMockContainerManager, WithMockNetworkManager)
	if err != nil {
		return err
	}

	req, err := lab.ToWritableLaboratory()
	if err != nil {
		return err
	}

	topology, err := usecase.GetTopology(req)
	if err != nil {
		return err
	}

	println(topology)
	return nil
}
