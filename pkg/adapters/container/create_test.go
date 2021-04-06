package container_test

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	cont "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/dpx/dpx/pkg/adapters/container"
	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/mock"
	"github.com/dpx/dpx/pkg/mock/assert"
	"github.com/dpx/dpx/pkg/utils"
)

func TestContainer_Create_WithDefaultOptions(t *testing.T) {
	d := &mock.Docker{}
	d.ContainerCreateFn = func(ctx context.Context, config *cont.Config, hostConfig *cont.HostConfig, networkingConfig *network.NetworkingConfig, name string) (cont.ContainerCreateCreatedBody, error) {
		h := defaultHostConfig()
		c := defaultConfig()

		expectHostConfig(t, h, hostConfig)
		expectConfig(t, c, config)

		assert.Equal(t, "", name)

		return cont.ContainerCreateCreatedBody{ID: "1111"}, nil
	}

	c := container.New(d)
	r, _ := c.Create(context.Background(), core.ContainerCreateOptions{
		Image: "test",
	})

	assert.Equal(t, "1111", r.ID)
}

func TestContainer_Create_WithPortOptions(t *testing.T) {
	d := &mock.Docker{}
	d.ContainerCreateFn = func(ctx context.Context, config *cont.Config, hostConfig *cont.HostConfig, networkingConfig *network.NetworkingConfig, name string) (cont.ContainerCreateCreatedBody, error) {
		h := defaultHostConfig()
		h.PortBindings = nat.PortMap{"5000/tcp": []nat.PortBinding{
			{HostIP: "", HostPort: "4000"},
		}}

		c := defaultConfig()
		c.ExposedPorts = nat.PortSet{
			"5000/tcp": struct{}{},
		}

		expectHostConfig(t, h, hostConfig)
		expectConfig(t, c, config)

		assert.Equal(t, "", name)

		return cont.ContainerCreateCreatedBody{}, nil
	}

	c := container.New(d)
	c.Create(context.Background(), core.ContainerCreateOptions{
		Image: "test",
		Ports: []string{
			"4000:5000",
		},
	})
}

func TestContainer_Create_WithVolumeOptions(t *testing.T) {
	d := &mock.Docker{}
	d.ContainerCreateFn = func(ctx context.Context, config *cont.Config, hostConfig *cont.HostConfig, networkingConfig *network.NetworkingConfig, name string) (cont.ContainerCreateCreatedBody, error) {
		h := defaultHostConfig()
		h.Mounts = []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: "named",
				Target: "/dev/vol1",
			},
		}

		c := defaultConfig()

		expectHostConfig(t, h, hostConfig)
		expectConfig(t, c, config)

		assert.Equal(t, "", name)

		return cont.ContainerCreateCreatedBody{}, nil
	}

	c := container.New(d)
	c.Create(context.Background(), core.ContainerCreateOptions{
		Image: "test",
		Volumes: []string{
			"named:/dev/vol1",
		},
	})
}

func TestContainer_Create_WithNetworkOptions(t *testing.T) {
}

func TestContainer_Create_WithEnvOptions(t *testing.T) {
}

func TestContainer_Start(t *testing.T) {
	d := &mock.Docker{}
	d.ContainerStartFn = func(ctx context.Context, cID string, options types.ContainerStartOptions) error {
		assert.Equal(t, "1111", cID)

		return nil
	}

	c := container.New(d)
	c.Start(context.Background(), "1111")
}

func expectConfig(t *testing.T, a *cont.Config, b *cont.Config) {
	t.Run("contaner config", func(t *testing.T) {
		t.Run("image", func(t *testing.T) {
			assert.Equal(t, a.Image, b.Image)
		})

		assert.Equal(t, a.Tty, b.Tty)
		assert.Equal(t, a.ExposedPorts, b.ExposedPorts)
		assert.Equal(t, a.OpenStdin, b.OpenStdin)
		assert.Equal(t, a.Tty, b.Tty)
	})
}

func defaultConfig() *cont.Config {
	return &cont.Config{
		ExposedPorts: nat.PortSet{},
		Tty:          true,
		OpenStdin:    true,
		WorkingDir:   utils.GetCwd(),
		Image:        "test",
	}
}

func defaultHostConfig() *cont.HostConfig {
	return &cont.HostConfig{
		PortBindings: nat.PortMap{},
		Mounts:       []mount.Mount{},
	}
}

func expectHostConfig(t *testing.T, a *cont.HostConfig, b *cont.HostConfig) {
	t.Run("host config", func(t *testing.T) {
		assert.Equal(t, a.PortBindings, b.PortBindings)

		dir := utils.GetCwd()

		// default volume
		dm := mount.Mount{
			Type:   mount.TypeBind,
			Source: dir,
			Target: dir,
		}

		expectHostMount(t, dm, b.Mounts[0])

		if len(a.Mounts) != 0 {
			t.Run("mounts", func(t *testing.T) {
				for i, m1 := range a.Mounts {
					m2 := b.Mounts[i+1]

					expectHostMount(t, m1, m2)
				}
			})
		}
	})
}

func expectHostMount(t *testing.T, m1 mount.Mount, m2 mount.Mount) {
	assert.Equal(t, m1.Type, m2.Type)
	assert.Equal(t, m1.Source, m2.Source)
	assert.Equal(t, m1.Target, m2.Target)
}
