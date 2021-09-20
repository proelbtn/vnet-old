package main

import (
	"time"

	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/repositories"
)

func main() {
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

	_, err = entities.NewLaboratory("test", []*entities.Container{container}, nil)
	if err != nil {
		panic(err)
	}

	_, err = manager.Create(container)
	if err != nil {
		panic(err)
	}

	err = manager.Start(container)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	err = manager.Stop(container)
	if err != nil {
		panic(err)
	}

	err = manager.Delete(container)
	if err != nil {
		panic(err)
	}
}
