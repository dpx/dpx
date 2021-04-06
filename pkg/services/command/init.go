package command

import (
	"context"
	"fmt"
	"io"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/services/config"
)

type InitCommand struct {
	core.ImageAdapter
	core.ContainerAdapter
	core.ConfigService
	Options *core.CommandInitOptions
	Config  *core.Config
	Out     io.Writer
}

func (c *InitCommand) Execute(ctx context.Context) error {
	parseInitConfig(c)

	if c.Options.Image == "" {
		return fmt.Errorf("IMAGE name is required")
	}

	// check if local image exists
	found, err := c.ImageAdapter.Find(ctx, c.Config.Image)
	if err != nil {
		return err
	}

	// pull image from registry
	if len(found) == 0 {
		fmt.Fprintf(c.Out, "Cannot find %s image locally\n", c.Config.Image)
		if err = c.ImageAdapter.Pull(ctx, c.Config.Image); err != nil {
			return err
		}
	}

	_, err = c.ContainerAdapter.Create(ctx, *config.ContainerCreateOptionsFromConfig(c.Config))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Out, "Creating a container %s\n", c.Config.Name)

	// Write config file
	if !c.Options.Skip {
		return c.ConfigService.Save(c.Config, core.ConfigSaveOptions{})
	}

	return nil
}

func NewInitCmd(i core.ImageAdapter, c core.ContainerAdapter, cfs core.ConfigService, opts *core.CommandInitOptions, cf *core.Config, out io.Writer) *InitCommand {
	return &InitCommand{
		ImageAdapter:     i,
		ContainerAdapter: c,
		ConfigService:    cfs,
		Config:           cf,
		Options:          opts,
		Out:              out,
	}
}

func parseInitConfig(c *InitCommand) {
	withParsers(
		parseImageOption(&c.Config.Image, c.Options.Image),
		parseNameOption(&c.Config.Name, c.Options.Name, c.Options.Image, c.Config.Path),
	)(c.Config)
}
