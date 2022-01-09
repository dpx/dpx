package mock

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// type DockerContainer struct{}

// type DockerImage struct{}

// type DockerClient struct {
// }

// func (d *DockerClient) NewClientWithOpts()

type Docker struct {
	ContainerCreateFn      func(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error)
	ContainerStartFn       func(ctx context.Context, cID string, options types.ContainerStartOptions) error
	ContainerStopFn        func(ctx context.Context, containerID string, timeout *time.Duration) error
	ContainerListFn        func(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ContainerExecCreateFn  func(ctx context.Context, container string, config types.ExecConfig) (types.IDResponse, error)
	ContainerExecAttachFn  func(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error)
	ContainerExecInspectFn func(ctx context.Context, execID string) (types.ContainerExecInspect, error)
	ContainerWaitFn        func(ctx context.Context, contID string, condition container.WaitCondition) (<-chan container.ContainerWaitOKBody, <-chan error)
	ContainerExecResizeFn  func(ctx context.Context, execID string, options types.ResizeOptions) error
}

func (d *Docker) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
	return d.ContainerCreateFn(ctx, config, hostConfig, networkingConfig, platform, containerName)
}

func (d *Docker) ContainerStart(ctx context.Context, cID string, options types.ContainerStartOptions) error {
	return d.ContainerStartFn(ctx, cID, options)
}

func (d *Docker) ContainerStop(ctx context.Context, cID string, t *time.Duration) error {
	return d.ContainerStopFn(ctx, cID, t)
}

func (d *Docker) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	return d.ContainerListFn(ctx, options)
}

func (d *Docker) ContainerExecCreate(ctx context.Context, id string, config types.ExecConfig) (types.IDResponse, error) {
	return d.ContainerExecCreateFn(ctx, id, config)
}

func (d *Docker) ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error) {
	return d.ContainerExecAttachFn(ctx, execID, config)
}

func (d *Docker) ContainerExecInspect(ctx context.Context, execID string) (types.ContainerExecInspect, error) {
	return d.ContainerExecInspectFn(ctx, execID)
}

func (d *Docker) ContainerWait(ctx context.Context, contID string, condition container.WaitCondition) (<-chan container.ContainerWaitOKBody, <-chan error) {
	return d.ContainerWaitFn(ctx, contID, condition)
}

func (d *Docker) ContainerExecResize(ctx context.Context, execID string, options types.ResizeOptions) error {
	return d.ContainerExecResizeFn(ctx, execID, options)
}
