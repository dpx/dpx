package container_test

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/dpx/dpx/pkg/adapters/container"
	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/mock"
	"github.com/dpx/dpx/pkg/mock/assert"
)

func TestContainer_Find_WithDefaultOptions(t *testing.T) {
	d := &mock.Docker{}
	d.ContainerListFn = func(ctx context.Context, opts types.ContainerListOptions) ([]types.Container, error) {
		o := defaultFindOptions()

		expectFindOptions(t, o, opts)

		return []types.Container{}, nil
	}

	c := container.New(d)
	c.Find(context.Background(), core.ContainerFindOptions{})
}

func TestContainer_Find_WithNameAndImageOptions(t *testing.T) {
	d := &mock.Docker{}
	d.ContainerListFn = func(ctx context.Context, opts types.ContainerListOptions) ([]types.Container, error) {
		o := defaultFindOptions()
		o.Filters.Add("name", "^dpx$")
		o.Filters.Add("ancestor", "golang")

		expectFindOptions(t, o, opts)

		return []types.Container{}, nil
	}

	c := container.New(d)
	c.Find(context.Background(), core.ContainerFindOptions{
		Name:  "dpx",
		Image: "golang",
	})
}

func TestContainer_Find_WithCompose(t *testing.T) {
	d := &mock.Docker{}
	f := &container.ComposeFinder{}
	f.Exec = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		assert.Equal(t, "docker-compose", name)
		assert.Equal(t, []string{"ps", "-q", "app"}, args)

		return []byte{}, nil
	}

	c := container.NewWithOpts(
		container.WithDocker(d),
		container.WithFinder(f),
	)
	c.Find(context.Background(), core.ContainerFindOptions{
		Compose: "app",
	})
}

func defaultFindOptions() types.ContainerListOptions {
	return types.ContainerListOptions{
		Filters: filters.NewArgs(),
	}
}

func expectFindOptions(t *testing.T, a types.ContainerListOptions, b types.ContainerListOptions) {
	t.Run("find", func(t *testing.T) {
		assert.Equal(t, a.Filters, b.Filters)
	})
}
