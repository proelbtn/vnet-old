package main

import (
	"net"
	"time"

	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/repositories"
)

func main() {
	networkManager := repositories.NewNetworkManger()
	containerManager, err := repositories.NewContainerManager()
	if err != nil {
		panic(err)
	}

	manager := repositories.NewLaboratoryManager(containerManager, networkManager)

	network, err := entities.NewNetwork("test")
	if err != nil {
		panic(err)
	}

	container, err := entities.NewContainer(
		"test",
		"docker.io/library/redis:alpine",
	)
	if err != nil {
		panic(err)
	}

	addr, net, err := net.ParseCIDR("192.168.0.1/24")
	if err != nil {
		panic(err)
	}

	port, err := entities.NewPort("eth0", network, []*entities.Address{
		{Addr: &addr, Net: net},
	})
	if err != nil {
		panic(err)
	}

	container.Ports = append(container.Ports, port)

	lab, err := entities.NewLaboratory("test", []*entities.Container{container}, []*entities.Network{network})
	if err != nil {
		panic(err)
	}

	manager.Start(lab)

	time.Sleep(5 * time.Second)

	manager.Stop(lab)
}
