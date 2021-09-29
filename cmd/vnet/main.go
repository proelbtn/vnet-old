package main

import (
	"log"
	"os"

	"github.com/proelbtn/vnet/pkg/repositories"
	"github.com/proelbtn/vnet/pkg/usecases"
	"github.com/urfave/cli/v2"
)

var usecase *usecases.InstantLaboratoryUsecase = nil

func main() {
	app := &cli.App{
		Name:  "vnet",
		Usage: "Virtual Network Laboratory",
		Commands: []*cli.Command{
			startCommand,
			stopCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getUsecase() (*usecases.InstantLaboratoryUsecase, error) {
	if usecase != nil {
		return usecase, nil
	}

	networkManager := repositories.NewNetworkManger()
	containerManager, err := repositories.NewContainerManager()
	if err != nil {
		return nil, err
	}

	laboratoryManager := repositories.NewLaboratoryManager(containerManager, networkManager)
	usecase := usecases.NewInstantLaboratoryUsecase(laboratoryManager)

	return usecase, nil
}
