package commands

import (
	"errors"

	"github.com/urfave/cli/v2"
)

var ExecCommand = &cli.Command{
	Name:      "exec",
	Usage:     "Execute command in laboratory",
	Flags:     append([]cli.Flag{}, commonFlags...),
	ArgsUsage: "CONTAINER CMD [ARG...]",
	Action:    exec,
}

func exec(c *cli.Context) error {
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

	if c.Args().Len() < 2 {
		return errors.New("container name and cmd are needed")
	}

	name := c.Args().First()
	args := c.Args().Slice()[1:]

	return usecase.Exec(req, name, args)
}
