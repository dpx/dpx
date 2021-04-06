package command

import (
	"context"
	"fmt"
	"io"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/services/config"
)

type ProcessCommand struct {
	core.ContainerAdapter
	Config  *core.Config
	Options *core.CommandProcessOptions
	Out     io.Writer
}

func (c *ProcessCommand) Execute(ctx context.Context) error {

	// find current docker process
	r, err := c.ContainerAdapter.Find(ctx, *config.ContainerFindOptionsFromConfig(c.Config))
	if err != nil {
		return err
	}

	// print process id
	fmt.Fprintln(c.Out, config.FormatProcessId(r.ID))

	return nil
}

func NewProcessCmd(c core.ContainerAdapter, cf *core.Config, opts *core.CommandProcessOptions, out io.Writer) *ProcessCommand {
	return &ProcessCommand{
		ContainerAdapter: c,
		Config:           cf,
		Options:          opts,
		Out:              out,
	}
}
