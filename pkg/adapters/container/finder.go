package container

import (
	"context"
	x "os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/dpx/dpx/pkg/core"
)

const composeCommand = "docker-compose"

type Finder interface {
	FindAll(ctx context.Context, d core.Docker, opts core.ContainerFindOptions) ([]core.Container, error)
}

type ContainerFinder struct{}
type ComposeFinder struct {
	Exec func(ctx context.Context, name string, args ...string) ([]byte, error)
}

// Find finds container process by using docker-compose
func (c ComposeFinder) FindAll(ctx context.Context, d core.Docker, opts core.ContainerFindOptions) ([]core.Container, error) {
	r, err := c.Exec(ctx, composeCommand, "ps", "-q", opts.Compose)
	if err != nil {
		return []core.Container{}, err
	}

	return []core.Container{
		{ID: string(r)},
	}, nil
}

// Find finds container process
func (c ContainerFinder) FindAll(ctx context.Context, d core.Docker, opts core.ContainerFindOptions) ([]core.Container, error) {
	filterArgs := filters.NewArgs()

	if opts.Name != "" {
		filterArgs.Add("name", "^" + opts.Name + "$")
	}

	if opts.Image != "" {
		filterArgs.Add("ancestor", opts.Image)
	}

	cs, err := d.ContainerList(ctx, types.ContainerListOptions{
		Filters: filterArgs,
		All:     opts.All,
	})

	rs := []core.Container{}

	if err != nil {
		return rs, err
	}

	// maps to core.Container
	for _, i := range cs {
		rs = append(rs, core.Container{ID: i.ID})
	}

	return rs, nil
}

func newFinder(compose string, f Finder) Finder {
	if f != nil {
		return f
	}

	if compose != "" {
		return newComposeFinder()
	}

	return &ContainerFinder{}
}

func newComposeFinder() Finder {
	c := &ComposeFinder{}
	c.Exec = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		return x.CommandContext(ctx, name, args...).Output()
	}

	return c
}
