package main

import (
	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/repositories"
)

func main() {
	manager := repositories.NewNetworkManger()

	network, err := entities.NewNetwork("test")
	if err != nil {
		panic(err)
	}

	err = manager.Create(network)
	if err != nil {
		panic(err)
	}

	err = manager.Delete(network)
	if err != nil {
		panic(err)
	}
}
