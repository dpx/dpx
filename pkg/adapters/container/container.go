package container

import (
	"github.com/dpx/dpx/pkg/core"
)

type ContainerAdapter struct {
	Finder
	Docker core.Docker
}

func New(d core.Docker) *ContainerAdapter {
	return &ContainerAdapter{
		Docker: d,
	}
}

type ContainerFn func(*ContainerAdapter)

func NewWithOpts(fns ...ContainerFn) *ContainerAdapter {
	c := &ContainerAdapter{}

	for _, fn := range fns {
		fn(c)
	}

	return c
}

func WithDocker(d core.Docker) ContainerFn {
	return func(ca *ContainerAdapter) {
		ca.Docker = d
	}
}

func WithFinder(f Finder) ContainerFn {
	return func(ca *ContainerAdapter) {
		ca.Finder = f
	}
}
