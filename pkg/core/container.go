package core

import (
	"context"
	"errors"
	"strconv"
)

type ContainerAdapter interface {
	ContainerCreateService
	ContainerFindService
	ContainerExecService
}

type ContainerCreateService interface {
	Create(ctx context.Context, options ContainerCreateOptions) (Container, error)
	Start(ctx context.Context, ID string) error
	Stop(ctx context.Context, ID string) error
}

type ContainerExecService interface {
	Exec(ctx context.Context, options ContainerExecOptions) error
	Inspect(ctx context.Context, ID string) (ContainerInspectResult, error)
}

type ContainerFindService interface {
	Find(ctx context.Context, options ContainerFindOptions) (Container, error)
	FindAll(ctx context.Context, options ContainerFindOptions) ([]Container, error)
}

type Container struct {
	ID string
}

type ContainerCreateOptions struct {
	Image      string
	Name       string
	Ports      []string
	Envs       []string
	Volumes    []string
	WorkingDir string
}
type ContainerCreateResult struct {
	ID string
}

type ContainerStartOptions struct{}
type ContainerStartResult struct{}
type ContainerInspectResult struct {
	Running  bool
	ExitCode int
}

type ContainerFindOptions struct {
	Name    string
	Image   string
	Compose string
	Volume  string
	All     bool
}

type ContainerExecOptions struct {
	ID     string
	Cmd    []string
	Config *CommandConfig
	Output bool
}

type ContainerExecErr struct {
	Code int
}

func (e *ContainerExecErr) Error() string {
	return strconv.Itoa(e.Code)
}

var ErrContainerNotFound = errors.New("Can't find the container")
