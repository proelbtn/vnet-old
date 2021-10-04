package commands

import (
	"github.com/urfave/cli/v2"
)

var StopCommand = &cli.Command{
	Name:   "stop",
	Usage:  "Stop laboratory",
	Flags:  append([]cli.Flag{}, commonFlags...),
	Action: stop,
}

func stop(c *cli.Context) error {
	lab, err := initialize(c)
	if err != nil {
		return err
	}

	usecase, err := newUsecase()
	if err != nil {
		return err
	}

	req, err := lab.ToWritableLaboratory()
	if err != nil {
		return err
	}

	return usecase.StopLaboratory(req)
}
