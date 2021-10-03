package main

import (
	"context"
	"net"
	"time"

	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/repositories"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	networkManager := repositories.NewNetworkManger()
	containerManager, err := repositories.NewContainerManager()
	if err != nil {
		panic(err)
	}

	manager := repositories.NewLaboratoryManager(containerManager, networkManager)

	network, err := entities.NewNetwork("test", 1500)
	if err != nil {
		panic(err)
	}

	addr, err := entities.NewIPAddress("192.168.0.1/24")
	if err != nil {
		panic(err)
	}

	port, err := entities.NewPort("eth0", network, []*net.IPNet{addr})
	if err != nil {
		panic(err)
	}

	container, err := entities.NewContainer(
		"test",
		"docker.io/nicolaka/netshoot:latest",
		[]*entities.Port{port},
		nil,
		nil,
	)
	if err != nil {
		panic(err)
	}

	lab, err := entities.NewLaboratory("test", []*entities.Container{container}, []*entities.Network{network})
	if err != nil {
		panic(err)
	}

	err = manager.Start(ctx, lab)
	if err != nil {
		panic(err)
	}

	time.Sleep(30 * time.Second)

	err = manager.Stop(ctx, lab)
	if err != nil {
		panic(err)
	}
}
