package commands

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var GetCommand = &cli.Command{
	Name:  "get",
	Usage: "Get information in laboratory",
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
	Subcommands: []*cli.Command{
		GetPortNameCommand,
	},
}

var GetPortNameCommand = &cli.Command{
	Name:      "port",
	Usage:     "Get port name in container",
	ArgsUsage: "CONTAINER PORT",
	Action:    getPortName,
}

func getPortName(c *cli.Context) error {
	if c.Args().Len() != 2 {
		return errors.New("container name and port name are needed")
	}

	containerName := c.Args().Get(0)
	portName := c.Args().Get(1)

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

	name, err := usecase.GetPortName(req, containerName, portName)
	if err != nil {
		return err
	}

	fmt.Println(name)
	return nil
}
