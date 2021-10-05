package commands

import (
	"github.com/urfave/cli/v2"
)

var RestartCommand = &cli.Command{
	Name:   "restart",
	Usage:  "Restart laboratory",
	Flags:  append([]cli.Flag{}, commonFlags...),
	Action: restart,
}

func restart(c *cli.Context) error {
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

	if err := usecase.StopLaboratory(req); err != nil {
		return err
	}

	if err := usecase.StartLaboratory(req); err != nil {
		return err
	}

	return nil
}
