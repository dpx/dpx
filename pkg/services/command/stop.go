package command

import (
	"context"
	"fmt"
	"io"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/services/config"
)

type StopCommand struct {
	core.ContainerAdapter
	Config  *core.Config
	Options *core.CommandStopOptions
	Out     io.Writer
}

func (c *StopCommand) Execute(ctx context.Context) error {
	// Check if container exists
	r, err := c.ContainerAdapter.Find(ctx, *config.ContainerFindOptionsFromConfig(c.Config))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Out, "Stopping a container %s\n", config.FormatProcessId(r.ID))

	return c.ContainerAdapter.Stop(ctx, r.ID)
}

func NewStopCmd(c core.ContainerAdapter, cf *core.Config, opts *core.CommandStopOptions, out io.Writer) *StopCommand {
	return &StopCommand{
		ContainerAdapter: c,
		Config:           cf,
		Options:          opts,
		Out:              out,
	}
}
