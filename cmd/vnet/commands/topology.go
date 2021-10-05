package commands

import (
	"github.com/urfave/cli/v2"
)

var TopologyCommand = &cli.Command{
	Name:   "topology",
	Usage:  "show topology",
	Flags:  append([]cli.Flag{}, commonFlags...),
	Action: topology,
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
