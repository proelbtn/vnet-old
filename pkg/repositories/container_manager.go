package repositories

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/usecases/managers"
	"go.uber.org/zap"
)

const (
	DEFAULT_NAMESPACE = "vnet"
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

func (v *ContainerManager) Create(ctx context.Context, spec *entities.Container) (uint32, error) {
	return v.createContainer(ctx, spec)
}

func (v *ContainerManager) Start(ctx context.Context, spec *entities.Container) error {
	return v.startTask(ctx, spec)
}

func (v *ContainerManager) Stop(ctx context.Context, spec *entities.Container) error {
	return v.stopTask(ctx, spec)
}

func (v *ContainerManager) Delete(ctx context.Context, spec *entities.Container) error {
	return v.deleteContainer(ctx, spec)
}

func getNamespaceName(spec *entities.Container) string {
	return DEFAULT_NAMESPACE
}

func getContainerID(spec *entities.Container) string {
	return spec.ID.String()
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

func (v *ContainerManager) getImage(ctx context.Context, name string) (containerd.Image, error) {
	zap.L().Debug("finding image", zap.String("name", name))
	images, err := v.client.ListImages(ctx, fmt.Sprintf("name==%s", name))
	if err != nil {
		return nil, err
	}
	if len(images) > 0 {
		zap.L().Debug("saved image found")
		return images[0], nil
	}

	zap.L().Debug("saved image not found, pulling")
	return v.client.Pull(ctx, name, containerd.WithPullUnpack)
}

func (v *ContainerManager) createContainer(ctx context.Context, spec *entities.Container) (uint32, error) {
	ctx = namespaces.WithNamespace(ctx, getNamespaceName(spec))
	logger := zap.L().With(
		zap.String("ID", spec.ID.String()),
		zap.String("Name", spec.Name),
		zap.String("Laboratory.ID", spec.Laboratory.ID.String()),
		zap.String("Laboratory.Name", spec.Laboratory.Name),
	)

	logger.Debug("creating Container")
	id := getContainerID(spec)

	logger.Debug("getting image", zap.String("ImageName", spec.ImageName))
	image, err := v.getImage(ctx, spec.ImageName)
	if err != nil {
		return 0, err
	}

	logger.Debug("creating container")
	container, err := v.client.NewContainer(
		ctx, id,
		containerd.WithNewSnapshot(id, image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return 0, err
	}

	logger.Debug("creating task")
	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return 0, err
	}

	logger.Debug("created Container")
	return task.Pid(), nil
}

func (v *ContainerManager) deleteContainer(ctx context.Context, spec *entities.Container) error {
	ctx = namespaces.WithNamespace(ctx, getNamespaceName(spec))
	logger := zap.L().With(
		zap.String("ID", spec.ID.String()),
		zap.String("Name", spec.Name),
		zap.String("Laboratory.ID", spec.Laboratory.ID.String()),
		zap.String("Laboratory.Name", spec.Laboratory.Name),
	)

	logger.Debug("deleting Container")
	id := getContainerID(spec)

	container, err := v.findContainer(ctx, id)
	if err != nil {
		return err
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return err
	}

	logger.Debug("deleting task")
	_, err = task.Delete(ctx)
	if err != nil {
		return err
	}

	logger.Debug("deleting container")
	err = container.Delete(ctx, containerd.WithSnapshotCleanup)

	logger.Debug("deleted Container")
	return err
}

func (v *ContainerManager) startTask(ctx context.Context, spec *entities.Container) error {
	ctx = namespaces.WithNamespace(ctx, getNamespaceName(spec))
	logger := zap.L().With(
		zap.String("ID", spec.ID.String()),
		zap.String("Name", spec.Name),
		zap.String("Laboratory.ID", spec.Laboratory.ID.String()),
		zap.String("Laboratory.Name", spec.Laboratory.Name),
	)

	logger.Debug("starting task")
	id := getContainerID(spec)

	container, err := v.findContainer(ctx, id)
	if err != nil {
		return err
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return err
	}

	err = task.Start(ctx)

	logger.Debug("started task")
	return err
}

func (v *ContainerManager) stopTask(ctx context.Context, spec *entities.Container) error {
	ctx = namespaces.WithNamespace(ctx, getNamespaceName(spec))
	logger := zap.L().With(
		zap.String("ID", spec.ID.String()),
		zap.String("Name", spec.Name),
		zap.String("Laboratory.ID", spec.Laboratory.ID.String()),
		zap.String("Laboratory.Name", spec.Laboratory.Name),
	)

	logger.Debug("killing task")
	id := getContainerID(spec)

	container, err := v.findContainer(ctx, id)
	if err != nil {
		return err
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return err
	}

	err = task.Kill(ctx, syscall.SIGKILL)

	logger.Debug("killing task")
	return err
}
