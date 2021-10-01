package commands

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var StopCommand = &cli.Command{
	Name:  "stop",
	Usage: "Stop laboratory",
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
	Action: stop,
}

func stop(c *cli.Context) error {
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

	err = usecase.StopLaboratory(req)
	if err != nil {
		return err
	}

	return nil
}
