package repositories

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/containerd/console"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/mattn/go-shellwords"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/proelbtn/vnet/pkg/entities"
	"github.com/proelbtn/vnet/pkg/errors"
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

	return "", errors.ErrNotFound
}

func (v *ContainerManager) getLogger(con *entities.Container) *zap.Logger {
	return zap.L().With(
		zap.String("Name", con.Name),
		zap.String("Laboratory.Name", con.Laboratory.Name),
	)
}

func (v *ContainerManager) Create(ctx context.Context, spec *entities.Container) (uint32, error) {
	return v.create(ctx, spec)
}

func (v *ContainerManager) Start(ctx context.Context, spec *entities.Container) error {
	return v.startTask(ctx, spec)
}

func (v *ContainerManager) Stop(ctx context.Context, spec *entities.Container) error {
	return v.stopTask(ctx, spec)
}

func (v *ContainerManager) Delete(ctx context.Context, spec *entities.Container) error {
	return v.delete(ctx, spec)
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

// findImage tries to find image specified with name
func (v *ContainerManager) findImage(ctx context.Context, name string) (containerd.Image, error) {
	images, err := v.client.ListImages(ctx, fmt.Sprintf("name==%s", name))
	if err != nil {
		return nil, err
	}
	if len(images) > 0 {
		return images[0], nil
	}

	return nil, errors.ErrNotFound
}

// findContainer tries to find container associated with id
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

	return nil, errors.ErrNotFound
}

// findTask tries to find task bound to container
func (v *ContainerManager) findTask(ctx context.Context, container containerd.Container) (containerd.Task, error) {
	task, err := container.Task(ctx, nil)
	if err != nil {
		if !errors.Is(err, errdefs.ErrNotFound) {
			return nil, err
		}
	}

	if task != nil {
		return task, nil
	}

	return nil, errors.ErrNotFound
}

// ensureImageExists tries to get image or fetch it if image is not found
func (v *ContainerManager) ensureImageExists(ctx context.Context, name string) (containerd.Image, error) {
	logger := zap.L().With(zap.String("name", name))

	logger.Debug("ensuring image exists")

	logger.Debug("finding image")
	image, err := v.findImage(ctx, name)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	if image != nil {
		logger.Debug("image found")
		return image, nil
	}

	logger.Debug("image not found, pulling")
	return v.client.Pull(ctx, name, containerd.WithPullUnpack)
}

// ensureContainerExists tries to get container or create new one if container is not found
func (v *ContainerManager) ensureContainerExists(ctx context.Context, name string, opts ...containerd.NewContainerOpts) (containerd.Container, error) {
	logger := zap.L().With(zap.String("name", name))

	logger.Debug("ensuring container exists")

	logger.Debug("finding container")
	container, err := v.findContainer(ctx, name)
	if err != nil {
		if !errors.IsNotFound(err) {
			return nil, err
		}
	}

	if container != nil {
		logger.Debug("container found")
		return container, nil
	}

	logger.Debug("container not found, creating")
	return v.client.NewContainer(ctx, name, opts...)
}

// ensureTaskExists tries to get task or create new one if task is not found
func (v *ContainerManager) ensureTaskExists(ctx context.Context, container containerd.Container) (containerd.Task, error) {
	logger := zap.L().With(zap.String("name", container.ID()))

	logger.Debug("ensuring task exists")

	logger.Debug("finding task")
	task, err := v.findTask(ctx, container)
	if err != nil {
		if !errors.IsNotFound(err) {
			return nil, err
		}
	}

	if task != nil {
		logger.Debug("task found")
		return task, nil
	}

	logger.Debug("task not found, creating")

	devnull, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	return container.NewTask(ctx, cio.NewCreator(cio.WithStreams(devnull, devnull, devnull)))
}

func (v *ContainerManager) ensureContainerNotExist(ctx context.Context, name string) error {
	logger := zap.L().With(zap.String("name", name))

	logger.Debug("ensuring container not exist")

	container, err := v.findContainer(ctx, name)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if container == nil {
		logger.Warn("container not found")
		return nil
	}

	err = v.ensureTaskNotExist(ctx, container)
	if err != nil {
		return err
	}

	return container.Delete(ctx, containerd.WithSnapshotCleanup)
}

func (v *ContainerManager) ensureTaskNotExist(ctx context.Context, container containerd.Container) error {
	logger := zap.L().With(zap.String("name", container.ID()))

	logger.Debug("ensuring task not exist")

	task, err := v.findTask(ctx, container)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if task == nil {
		logger.Warn("task not found")
		return nil
	}

	status, err := task.Status(ctx)
	if err != nil {
		return err
	}

	if status.Status != containerd.Stopped {
		pchan, err := task.Wait(ctx)
		if err != nil {
			return err
		}

		if err := task.Kill(ctx, syscall.SIGKILL, containerd.WithKillAll); err != nil {
			return err
		}

		<-pchan
	}

	_, err = task.Delete(ctx)
	return err
}

func (v *ContainerManager) execute(ctx context.Context, task containerd.Task, id string, pspec *specs.Process, creator cio.Creator) error {
	logger := zap.L().With(zap.Any("args", pspec.Args))

	logger.Debug("executing command")
	proc, err := task.Exec(ctx, id, pspec, creator)
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

	return nil
}

// create container based on spec
func (v *ContainerManager) create(ctx context.Context, spec *entities.Container) (uint32, error) {
	ctx = namespaces.WithNamespace(ctx, getNamespaceName(spec))
	logger := v.getLogger(spec)

	logger.Debug("creating container")
	name := getContainerName(spec)

	image, err := v.ensureImageExists(ctx, spec.ImageName)
	if err != nil {
		return 0, err
	}

	mounts := make([]specs.Mount, len(spec.Volumes))
	for i, volume := range spec.Volumes {
		mounts[i] = specs.Mount{
			Source:      volume.Source,
			Destination: volume.Destination,
			Type:        "none",
			Options:     []string{"ro", "rbind"},
		}
	}

	opts := []containerd.NewContainerOpts{
		containerd.WithNewSnapshot(name, image),
		containerd.WithNewSpec(
			oci.WithImageConfig(image),
			oci.WithHostname(spec.Name),
			oci.WithPrivileged,
			oci.WithAllDevicesAllowed,
			oci.WithHostDevices,
			oci.WithMounts(mounts),
			oci.WithoutRunMount,
		),
	}

	container, err := v.ensureContainerExists(ctx, name, opts...)
	if err != nil {
		return 0, err
	}

	task, err := v.ensureTaskExists(ctx, container)
	if err != nil {
		return 0, err
	}

	logger.Debug("created container", zap.Uint32("pid", task.Pid()))
	return task.Pid(), nil
}

// delete container based on spec
func (v *ContainerManager) delete(ctx context.Context, spec *entities.Container) error {
	ctx = namespaces.WithNamespace(ctx, getNamespaceName(spec))
	logger := v.getLogger(spec)

	logger.Debug("deleting container")
	name := getContainerName(spec)

	err := v.ensureContainerNotExist(ctx, name)
	if err != nil {
		return err
	}

	logger.Debug("deleted container")
	return nil
}

// TODO: refactor
func (v *ContainerManager) startTask(ctx context.Context, con *entities.Container) error {
	ctx = namespaces.WithNamespace(ctx, getNamespaceName(con))
	logger := v.getLogger(con)

	logger.Debug("starting task")
	name := getContainerName(con)

	container, err := v.findContainer(ctx, name)
	if err != nil {
		return err
	}

	task, err := v.findTask(ctx, container)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	status, err := task.Status(ctx)
	if err != nil {
		return err
	}

	switch status.Status {
	case containerd.Running:
		return nil
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

		pspec := spec.Process
		pspec.Args = args

		if err := v.execute(ctx, task, "vnet-exec", pspec, cio.NullIO); err != nil {
			return err
		}
	}

	return err
}

// TODO: refactor
func (v *ContainerManager) stopTask(ctx context.Context, spec *entities.Container) error {
	ctx = namespaces.WithNamespace(ctx, getNamespaceName(spec))
	logger := v.getLogger(spec)

	logger.Debug("killing task")
	name := getContainerName(spec)

	container, err := v.findContainer(ctx, name)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if container == nil {
		logger.Warn("container not found")
		return nil
	}

	task, err := v.findTask(ctx, container)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if task == nil {
		logger.Warn("task not found")
		return nil
	}

	status, err := task.Status(ctx)
	if err != nil {
		return err
	}

	if status.Status != containerd.Running {
		return nil
	}

	return task.Kill(ctx, syscall.SIGKILL)
}

// TODO: refactor
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
