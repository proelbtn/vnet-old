package main

import (
	"fmt"
	"io/ioutil"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var createCommand = &cli.Command{
	Name:  "create",
	Usage: "Create laboratory",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "manifest",
			Value: "./lab.yaml",
			Usage: "manifest for laboratory",
		},
	},
	Action: start,
}

func loadManifest(manifestPath string) (*Laboratory, error) {
	var lab Laboratory

	manifest, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(manifest, &lab)
	if err != nil {
		return nil, err
	}

	return &lab, nil
}

func start(c *cli.Context) error {
	manifestPath := c.String("manifest")

	lab, err := loadManifest(manifestPath)
	if err != nil {
		return err
	}

	req := lab.ToWritableLaboratory()

	fmt.Printf("%+v\n", req)

	return nil
}
