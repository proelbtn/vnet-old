package repositories

import (
	"context"
	"errors"
	"os"
	"syscall"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/usecases/managers"
)

type ContainerManager struct {
	client *containerd.Client
}

var _ managers.ContainerManager = (*ContainerManager)(nil)

var (
	ErrNotFound = errors.New("not found")
)

var socketPathCandidates = []string{
	"/run/containerd/containerd.sock",
	"/var/run/containerd/containerd.sock",
	"/var/run/docker/containerd/containerd.sock",
}

func NewContainerManager() (*ContainerManager, error) {
	socketPath, err := FindContainerdSocketPath()
	if err != nil {
		return nil, err
	}

	client, err := containerd.New(socketPath)
	if err != nil {
		return nil, err
	}

	return &ContainerManager{
		client: client,
	}, nil
}

func FindContainerdSocketPath() (string, error) {
	for _, path := range socketPathCandidates {
		info, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", err
		}

		if info.Mode()&os.ModeSocket != 0 {
			return path, nil
		}
	}

	return "", ErrNotFound
}

func GetNamespaceName(spec *entities.Container) string {
	return spec.Laboratory.ID.String()
}

func GetContainerID(spec *entities.Container) string {
	return spec.ID.String()
}

func (v *ContainerManager) Create(spec *entities.Container) (uint32, error) {
	return v.createContainer(spec)
}

func (v *ContainerManager) Start(spec *entities.Container) error {
	return v.startTask(spec)
}

func (v *ContainerManager) Stop(spec *entities.Container) error {
	return v.stopTask(spec)
}

func (v *ContainerManager) Delete(spec *entities.Container) error {
	return v.deleteContainer(spec)
}

func (v *ContainerManager) createContainer(spec *entities.Container) (uint32, error) {
	ctx := namespaces.WithNamespace(context.Background(), GetNamespaceName(spec))
	id := GetContainerID(spec)

	image, err := v.client.Pull(ctx, spec.ImageName, containerd.WithPullUnpack)
	if err != nil {
		return 0, err
	}

	container, err := v.client.NewContainer(
		ctx, id,
		containerd.WithNewSnapshot(id, image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return 0, err
	}

	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return 0, err
	}

	return task.Pid(), nil
}

func (v *ContainerManager) findContainer(ctx context.Context, id string) (containerd.Container, error) {
	containers, err := v.client.Containers(ctx, "")
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		if container.ID() == id {
			return container, nil
		}
	}

	return nil, ErrNotFound
}

func (v *ContainerManager) deleteContainer(spec *entities.Container) error {
	ctx := namespaces.WithNamespace(context.Background(), GetNamespaceName(spec))
	id := GetContainerID(spec)

	container, err := v.findContainer(ctx, id)
	if err != nil {
		return err
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return err
	}

	_, err = task.Delete(ctx)
	if err != nil {
		return err
	}

	return container.Delete(ctx, containerd.WithSnapshotCleanup)
}

func (v *ContainerManager) startTask(spec *entities.Container) error {
	ctx := namespaces.WithNamespace(context.Background(), GetNamespaceName(spec))
	id := GetContainerID(spec)

	container, err := v.findContainer(ctx, id)
	if err != nil {
		return err
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return err
	}

	return task.Start(ctx)
}

func (v *ContainerManager) stopTask(spec *entities.Container) error {
	ctx := namespaces.WithNamespace(context.Background(), GetNamespaceName(spec))
	id := GetContainerID(spec)

	container, err := v.findContainer(ctx, id)
	if err != nil {
		return err
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return err
	}

	return task.Kill(ctx, syscall.SIGKILL)
}
