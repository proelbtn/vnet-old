package repositories

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/containerd/console"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/mattn/go-shellwords"
	"github.com/opencontainers/runtime-spec/specs-go"
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

func (v *ContainerManager) getLogger(con *entities.Container) *zap.Logger {
	return zap.L().With(
		zap.String("Name", con.Name),
		zap.String("Laboratory.Name", con.Laboratory.Name),
	)
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

func (v *ContainerManager) Exec(ctx context.Context, spec *entities.Container, args managers.ExecArgs) error {
	return v.exec(ctx, spec, args)
}

func getNamespaceName(spec *entities.Container) string {
	return DEFAULT_NAMESPACE
}

func getContainerName(spec *entities.Container) string {
	return fmt.Sprintf("%s-%s", spec.Laboratory.Name, spec.Name)
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

// refs: https://github.com/containerd/containerd/blob/63b7e5771e8914f3c36c707f3b5fc4846b11997b/cmd/ctr/commands/run/run_unix.go#L89
func (v *ContainerManager) createContainer(ctx context.Context, spec *entities.Container) (uint32, error) {
	ctx = namespaces.WithNamespace(ctx, getNamespaceName(spec))
	logger := v.getLogger(spec)

	logger.Debug("creating Container")
	name := getContainerName(spec)

	logger.Debug("getting image", zap.String("ImageName", spec.ImageName))
	image, err := v.getImage(ctx, spec.ImageName)
	if err != nil {
		return 0, err
	}

	logger.Debug("creating container")
	mounts := make([]specs.Mount, len(spec.Volumes))
	for i, volume := range spec.Volumes {
		mounts[i] = specs.Mount{
			Source:      volume.Source,
			Destination: volume.Destination,
			Type:        "none",
			Options:     []string{"bind"},
		}
	}

	container, err := v.client.NewContainer(
		ctx, name,
		containerd.WithNewSnapshot(name, image),
		containerd.WithNewSpec(
			oci.WithImageConfig(image),
			oci.WithPrivileged,
			oci.WithAllDevicesAllowed,
			oci.WithHostDevices,
			oci.WithMounts(mounts),
		),
	)
	if err != nil {
		return 0, err
	}

	logger.Debug("creating task")
	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return 0, err
	}

	logger.Debug("created Container", zap.Uint32("pid", task.Pid()))
	return task.Pid(), nil
}

func (v *ContainerManager) deleteContainer(ctx context.Context, spec *entities.Container) error {
	ctx = namespaces.WithNamespace(ctx, getNamespaceName(spec))
	logger := v.getLogger(spec)

	logger.Debug("deleting Container")
	name := getContainerName(spec)

	container, err := v.findContainer(ctx, name)
	if err != nil {
		// TODO: handling error
		return nil
	}

	task, _ := container.Task(ctx, nil)
	// TODO: handling error

	if task != nil {
		logger.Debug("deleting task")
		_, err = task.Delete(ctx)
		if err != nil {
			return err
		}
	}

	logger.Debug("deleting container")
	err = container.Delete(ctx, containerd.WithSnapshotCleanup)

	logger.Debug("deleted Container")
	return err
}

func (v *ContainerManager) startTask(ctx context.Context, con *entities.Container) error {
	ctx = namespaces.WithNamespace(ctx, getNamespaceName(con))
	logger := v.getLogger(con)

	logger.Debug("starting task")
	name := getContainerName(con)

	container, err := v.findContainer(ctx, name)
	if err != nil {
		return err
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return err
	}

	err = task.Start(ctx)
	if err != nil {
		return err
	}
	logger.Debug("started task")

	logger.Debug("executing commands")
	spec, err := container.Spec(ctx)
	if err != nil {
		return err
	}

	for _, cmd := range con.Commands {
		args, err := shellwords.Parse(cmd)
		if err != nil {
			return err
		}

		logger := logger.With(zap.Any("command", args))

		pspec := spec.Process
		pspec.Args = args

		logger.Debug("executing command")
		proc, err := task.Exec(ctx, "vnet-exec", pspec, cio.NullIO)
		if err != nil {
			return err
		}

		pchan, err := proc.Wait(ctx)
		if err != nil {
			return err
		}

		logger.Debug("starting command")
		if err = proc.Start(ctx); err != nil {
			return err
		}

		logger.Debug("waiting command")
		status := <-pchan

		logger.Debug("executed command", zap.Uint32("code", status.ExitCode()))

		if _, err := proc.Delete(ctx); err != nil {
			return err
		}
	}

	return err
}

func (v *ContainerManager) stopTask(ctx context.Context, spec *entities.Container) error {
	ctx = namespaces.WithNamespace(ctx, getNamespaceName(spec))
	logger := v.getLogger(spec)

	logger.Debug("killing task")
	name := getContainerName(spec)

	container, err := v.findContainer(ctx, name)
	if err != nil {
		return err
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		// TODO: handling errors which isn't `not found` error
		return nil
	}

	err = task.Kill(ctx, syscall.SIGKILL)

	logger.Debug("killing task")
	return err
}

func (v *ContainerManager) exec(ctx context.Context, con *entities.Container, args managers.ExecArgs) error {
	ctx = namespaces.WithNamespace(ctx, getNamespaceName(con))
	logger := v.getLogger(con).With(zap.Any("command", args))

	logger.Debug("starting task")
	container, err := v.findContainer(ctx, getContainerName(con))
	if err != nil {
		return err
	}

	logger.Debug("executing commands")
	spec, err := container.Spec(ctx)
	if err != nil {
		return err
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return err
	}

	// TODO: it's really okay?
	// ref: https://github.com/containerd/containerd/blob/main/cmd/ctr/commands/tasks/exec.go
	pspec := spec.Process
	pspec.Args = args.Args
	pspec.Terminal = true

	logger.Debug("executing command")
	consol := console.Current()
	if err := consol.SetRaw(); err != nil {
		return err
	}

	proc, err := task.Exec(ctx, "vnet-exec", pspec, cio.NewCreator(
		cio.WithStreams(consol, consol, nil),
		cio.WithTerminal,
	))
	if err != nil {
		return err
	}

	pchan, err := proc.Wait(ctx)
	if err != nil {
		return err
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	logger.Debug("starting command")
	if err = proc.Start(ctx); err != nil {
		return err
	}

	logger.Debug("waiting command")
	select {
	case <-sigs:
		logger.Debug("signal received, killing command")
		err := proc.Kill(ctx, syscall.SIGKILL)
		if err != nil {
			return err
		}
		logger.Debug("killed command")
	case <-pchan:
		status := <-pchan
		logger.Debug("executed command", zap.Uint32("code", status.ExitCode()))
	}

	if err := consol.Reset(); err != nil {
		return err
	}

	if _, err := proc.Delete(ctx); err != nil {
		return err
	}

	logger.Debug("executed command")
	return err
}
