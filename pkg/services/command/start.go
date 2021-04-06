package command

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/services/config"
)

type StartCommand struct {
	core.ContainerAdapter
	Config  *core.Config
	Options *core.CommandStartOptions
	Out     io.Writer
}

func (c *StartCommand) Execute(ctx context.Context) error {
	parseStartConfig(c)

	// Check if container exists
	r, err := c.ContainerAdapter.Find(ctx, *config.ContainerFindOptionsFromConfig(c.Config))
	if err != nil && !errors.Is(err, core.ErrContainerNotFound) {
		return err
	}

	// Can't find the container. Let's create a new one.
	if r.ID == "" {
		r, err = c.ContainerAdapter.Create(ctx, *config.ContainerCreateOptionsFromConfig(c.Config))
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(c.Out, "Starting a container %s\n", config.FormatProcessId(r.ID))

	return c.ContainerAdapter.Start(ctx, r.ID)
}

func parseStartConfig(c *StartCommand) {
	withParsers(
		parseEnvsOption(&c.Config.Envs, c.Options.Envs),
		parsePortsOption(&c.Config.Ports, c.Options.Ports),
		parseVolumesOption(&c.Config.Volumes, c.Options.Volumes),
	)(c.Config)
}

func NewStartCmd(c core.ContainerAdapter, cf *core.Config, opts *core.CommandStartOptions, out io.Writer) *StartCommand {
	return &StartCommand{
		ContainerAdapter: c,
		Config:           cf,
		Options:          opts,
		Out:              out,
	}
}
