package main

import (
	"context"
	"time"

	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/repositories"
)

func main() {
	ctx := context.Background()

	manager, err := repositories.NewContainerManager()
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

	_, err = entities.NewLaboratory(
		"test",
		entities.WithContainer(container),
	)
	if err != nil {
		panic(err)
	}

	_, err = manager.Create(ctx, container)
	if err != nil {
		panic(err)
	}

	err = manager.Start(ctx, container)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	err = manager.Stop(ctx, container)
	if err != nil {
		panic(err)
	}

	err = manager.Delete(ctx, container)
	if err != nil {
		panic(err)
	}
}
