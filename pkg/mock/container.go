package mock

import (
	"context"

	"github.com/dpx/dpx/pkg/core"
)

type ContainerAdapter struct {
	CreateFn  func(ctx context.Context, options core.ContainerCreateOptions) (core.Container, error)
	StartFn   func(ctx context.Context, ID string) error
	StopFn    func(ctx context.Context, ID string) error
	ExecFn    func(ctx context.Context, options core.ContainerExecOptions) error
	InspectFn func(ctx context.Context, ID string) (core.ContainerInspectResult, error)
	FindFn    func(ctx context.Context, options core.ContainerFindOptions) (core.Container, error)
	FindAllFn func(ctx context.Context, options core.ContainerFindOptions) ([]core.Container, error)
}

func (c *ContainerAdapter) Create(ctx context.Context, options core.ContainerCreateOptions) (core.Container, error) {
	return c.CreateFn(ctx, options)
}

func (c *ContainerAdapter) Start(ctx context.Context, ID string) error {
	return c.StartFn(ctx, ID)
}

func (c *ContainerAdapter) Stop(ctx context.Context, ID string) error {
	return c.StopFn(ctx, ID)
}

func (c *ContainerAdapter) Exec(ctx context.Context, options core.ContainerExecOptions) error {
	return c.ExecFn(ctx, options)
}

func (c *ContainerAdapter) Inspect(ctx context.Context, ID string) (core.ContainerInspectResult, error) {
	return c.InspectFn(ctx, ID)
}

func (c *ContainerAdapter) Find(ctx context.Context, options core.ContainerFindOptions) (core.Container, error) {
	return c.FindFn(ctx, options)
}

func (c *ContainerAdapter) FindAll(ctx context.Context, options core.ContainerFindOptions) ([]core.Container, error) {
	return c.FindAllFn(ctx, options)
}

type ContainerFn func(*ContainerAdapter) *ContainerAdapter

func NewContainerAdapter(fns ...ContainerFn) *ContainerAdapter {
	c := &ContainerAdapter{}

	for _, f := range fns {
		f(c)
	}

	return c
}

func ContainerAdapterWithFind(ID string) ContainerFn {
	return func(cs *ContainerAdapter) *ContainerAdapter {
		cs.FindFn = func(ctx context.Context, options core.ContainerFindOptions) (core.Container, error) {
			return core.Container{
				ID: ID,
			}, nil
		}

		return cs
	}
}

func ContainerAdapterWithCreate(ID string) ContainerFn {
	return func(cs *ContainerAdapter) *ContainerAdapter {
		cs.CreateFn = func(ctx context.Context, options core.ContainerCreateOptions) (core.Container, error) {
			return core.Container{
				ID: ID,
			}, nil
		}

		return cs
	}
}

func ContainerAdapterWithStart() ContainerFn {
	return func(cs *ContainerAdapter) *ContainerAdapter {
		cs.StartFn = func(ctx context.Context, ID string) error {
			return nil
		}

		return cs
	}
}
