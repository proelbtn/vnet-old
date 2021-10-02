package commands

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var StartCommand = &cli.Command{
	Name:  "start",
	Usage: "Start laboratory",
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
	Action: start,
}

func start(c *cli.Context) error {
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

	if err := tryStart(c); err != nil {
		zap.L().Error("could not start laboratory", zap.Error(err))
	}

	return nil
}

func tryStart(c *cli.Context) error {
	usecase, err := getUsecase()
	if err != nil {
		return err
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

	return usecase.StartLaboratory(req)
}
