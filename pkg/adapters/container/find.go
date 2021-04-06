package container

import (
	"context"

	"github.com/dpx/dpx/pkg/core"
)

// FindAll finds current docker process (including stopped ones) and return container id
func (c *ContainerAdapter) FindAll(ctx context.Context, opts core.ContainerFindOptions) ([]core.Container, error) {
	f := newFinder(opts.Compose, c.Finder)

	return f.FindAll(ctx, c.Docker, opts)
}

// Find finds current docker process and return container id
func (c *ContainerAdapter) Find(ctx context.Context, opts core.ContainerFindOptions) (core.Container, error) {
	r := core.Container{}

	cs, err := c.FindAll(ctx, opts)
	if err != nil {
		return r, err
	}

	if len(cs) == 0 {
		return r, core.ErrContainerNotFound
	}

	return cs[0], nil
}
