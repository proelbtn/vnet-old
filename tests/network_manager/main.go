package main

import (
	"context"

	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/repositories"
)

func main() {
	ctx := context.Background()
	manager := repositories.NewNetworkManger()

	network, err := entities.NewNetwork("test", 1500)
	if err != nil {
		panic(err)
	}

	err = manager.Create(ctx, network)
	if err != nil {
		panic(err)
	}

	err = manager.Delete(ctx, network)
	if err != nil {
		panic(err)
	}
}
