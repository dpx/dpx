package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/services/config"
)

var ErrorMissingCommand = errors.New("command name is required (e.g. dpx l go)")

type LinkCommand struct {
	core.ContainerAdapter
	core.ConfigService
	Config  *core.Config
	Options *core.CommandLinkOptions
	Out     io.Writer
}

func (c *LinkCommand) Execute(ctx context.Context) error {

	if len(c.Options.Cmds) == 0 && !c.Options.All {
		return ErrorMissingCommand
	}

	// Find current docker process
	r, err := c.ContainerAdapter.Find(ctx, *config.ContainerFindOptionsFromConfig(c.Config))
	if err != nil {
		return err
	}

	// Make sure dpx-alias file exists
	if _, err = c.ConfigService.CreateAliasFile(); err != nil {
		return err
	}

	cmds := []string{}
	if c.Options.All {
		cmds = allCmdsFromConfig(c.Config)
	} else {
		cmds = c.Options.Cmds
	}

	return linkAll(ctx, c, r.ID, cmds, c.Out)
}

// linkAll links all commands defined in config file
func linkAll(ctx context.Context, c *LinkCommand, ID string, cmds []string, out io.Writer) error {
	for _, cmd := range cmds {
		if err := link(ctx, c, ID, cmd, out); err != nil {
			return err
		}
	}

	return nil
}

// link links single binary
func link(ctx context.Context, c *LinkCommand, ID string, cmd string, out io.Writer) error {
	cmds := buildTestCommand(c.Config.Defaults.Shell, cmd)
	opts := &core.ContainerExecOptions{
		ID:     ID,
		Cmd:    cmds,
		Config: c.Config.Defaults,
	}

	err := c.ContainerAdapter.Exec(ctx, *opts)

	if err != nil {
		return fmt.Errorf("Cannot link %s: %w", cmd, err)
	}

	if _, err = c.ConfigService.CreateBinFile(cmd); err != nil {
		return err
	}

	fmt.Fprintf(out, "LINK: %s (%s)\n", cmd, path.Join(config.BinDir, cmd))

	return nil
}

func allCmdsFromConfig(c *core.Config) []string {
	cmds := []string{}
	for cmd := range c.Commands {
		cmds = append(cmds, cmd)
	}

	return cmds
}

func buildTestCommand(shell, cmd string) []string {
	return []string{
		shell,
		"-c",
		fmt.Sprintf("command -v %s >/dev/null 2>&1", cmd),
	}
}

func NewLinkCmd(c core.ContainerAdapter, cfs core.ConfigService, cf *core.Config, opts *core.CommandLinkOptions, out io.Writer) *LinkCommand {
	return &LinkCommand{
		ContainerAdapter: c,
		ConfigService:    cfs,
		Config:           cf,
		Options:          opts,
		Out:              out,
	}
}
