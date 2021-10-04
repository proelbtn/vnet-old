package commands

import (
	"io/ioutil"

	"github.com/proelbtn/vnet/pkg/repositories"
	"github.com/proelbtn/vnet/pkg/usecases"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const (
	FlagDebug        = "debug"
	FlagManifest     = "manifest"
	FlagOverrideName = "override-name"
)

var usecase *usecases.LaboratoryUsecase = nil
var commonFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:  FlagDebug,
		Value: false,
		Usage: "debug",
	},
	&cli.StringFlag{
		Name:  FlagManifest,
		Value: "./lab.yaml",
		Usage: "manifest for laboratory",
	},
	&cli.StringFlag{
		Name:  FlagOverrideName,
		Value: "",
		Usage: "override laboratory name",
	},
}

func initializeLogger(c *cli.Context) error {
	var logger *zap.Logger
	var err error

	if c.Bool(FlagDebug) {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)

	return nil
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

func initialize(c *cli.Context) (*Laboratory, error) {
	err := initializeLogger(c)
	if err != nil {
		return nil, err
	}

	manifestPath := c.String(FlagManifest)

	lab, err := loadManifest(manifestPath)
	if err != nil {
		return nil, err
	}

	overrideName := c.String(FlagOverrideName)
	if overrideName != "" {
		lab.Name = overrideName
	}

	return lab, err
}

func getUsecase() (*usecases.LaboratoryUsecase, error) {
	if usecase != nil {
		return usecase, nil
	}

	networkManager := repositories.NewNetworkManger()
	containerManager, err := repositories.NewContainerManager()
	if err != nil {
		return nil, err
	}

	laboratoryManager := repositories.NewLaboratoryManager(containerManager, networkManager)
	usecase := usecases.NewLaboratoryUsecase(laboratoryManager, containerManager, networkManager)

	return usecase, nil
}
