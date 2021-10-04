package commands

import (
	"github.com/urfave/cli/v2"
)

var StartCommand = &cli.Command{
	Name:   "start",
	Usage:  "Start laboratory",
	Flags:  append([]cli.Flag{}, commonFlags...),
	Action: start,
}

func start(c *cli.Context) error {
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

	return usecase.StartLaboratory(req)
}
