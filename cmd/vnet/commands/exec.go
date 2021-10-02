package commands

import (
	"errors"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var ExecCommand = &cli.Command{
	Name:  "exec",
	Usage: "Execute command in laboratory",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Value: false,
			Usage: "debug",
		},
		&cli.StringFlag{
			Name:  "manifest",
			Value: "./lab.yaml",
			Usage: "manifest for laboratory",
		},
	},
	ArgsUsage: "CONTAINER CMD [ARG...]",
	Action:    exec,
}

func exec(c *cli.Context) error {
	usecase, err := getUsecase()
	if err != nil {
		return err
	}

	if c.Bool("debug") {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return err
		}
		zap.ReplaceGlobals(logger)
	} else {
		logger, err := zap.NewProduction()
		if err != nil {
			return err
		}
		zap.ReplaceGlobals(logger)
	}

	manifestPath := c.String("manifest")

	lab, err := loadManifest(manifestPath)
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
