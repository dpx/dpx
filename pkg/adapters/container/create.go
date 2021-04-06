package container

import (
	"context"
	"runtime"

	"github.com/docker/docker/api/types"
	cont "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/volume/mounts"
	"github.com/docker/go-connections/nat"
	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/utils"
)

// Create creates container
func (c *ContainerAdapter) Create(ctx context.Context, opts core.ContainerCreateOptions) (core.Container, error) {
	r := core.Container{}

	// parse ports
	ports, hostPorts, err := parsePorts(opts.Ports)
	if err != nil {
		return r, err
	}

	// parse volumes
	mnts, err := parseVolumes(opts.Volumes)
	if err != nil {
		return r, err
	}

	dir := utils.GetCwd()
	result, err := c.Docker.ContainerCreate(ctx,
		&cont.Config{
			Image:        opts.Image,
			Tty:          true,
			Env:          opts.Envs,
			OpenStdin:    true,
			WorkingDir:   dir,
			ExposedPorts: ports,
		},
		// https://stackoverflow.com/a/48472934
		&cont.HostConfig{
			Mounts:       mnts,
			PortBindings: hostPorts,
		},
		&network.NetworkingConfig{},
		opts.Name,
	)

	return core.Container{ID: result.ID}, err
}

// Start starts container
func (c *ContainerAdapter) Start(ctx context.Context, ID string) error {
	return c.Docker.ContainerStart(ctx, ID, types.ContainerStartOptions{})
}

// Start starts container
func (c *ContainerAdapter) Stop(ctx context.Context, ID string) error {
	return c.Docker.ContainerStop(ctx, ID, nil)
}

func defaultVolume() mount.Mount {
	dir := utils.GetCwd()

	return mount.Mount{
		Type:        mount.TypeBind,
		Source:      dir,
		Target:      dir,
		Consistency: mount.ConsistencyCached,
	}
}

func parsePorts(ports []string) (map[nat.Port]struct{}, map[nat.Port][]nat.PortBinding, error) {
	return nat.ParsePortSpecs(ports)
}

func parseVolumes(vols []string) ([]mount.Mount, error) {
	parser := mounts.NewParser(runtime.GOOS)

	// default volume
	mnts := []mount.Mount{}
	mnts = append(mnts, defaultVolume())

	var err error
	var mp *mounts.MountPoint

	for _, vol := range vols {
		mp, err = parser.ParseMountRaw(vol, "local")
		if err != nil {
			return mnts, err
		}

		mnts = append(mnts, mp.Spec)
	}

	return mnts, err
}
