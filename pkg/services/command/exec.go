package command

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/services/config"
)

type ExecCommand struct {
	core.ContainerAdapter
	core.ConfigService
	Config  *core.Config
	Options *core.CommandExecOptions
	Out     io.Writer
}

func (c *ExecCommand) Execute(ctx context.Context) error {
	parseExecConfig(c)

	// execute all commands defined in `runs`
	if c.Options.All {
		return exec(ctx, c, c.Config, core.NewCommandConfig(), c.Config.Runs, " && ", c.Out)
	}

	if c.Options.Cmd == "" {
		return fmt.Errorf("command is required")
	}

	// execute single command e.g. `dpx go -v`
	cmdConfig := c.ConfigService.Command(c.Options.Cmd, c.Config)

	// print command's config
	if c.Options.Print {
		return execPrintConfig(cmdConfig, c.Out)
	}

	if err := exec(ctx, c, c.Config, cmdConfig, c.Options.CmdArgs, " ", c.Out); err != nil {
		return err
	}

	// save command
	if c.Options.Save {
		c.Config.Runs = append(c.Config.Runs, strings.Join(c.Options.CmdArgs, " && "))
		c.ConfigService.Save(c.Config, core.ConfigSaveOptions{})
	}

	return nil
}

func exec(ctx context.Context, c *ExecCommand, cf *core.Config, cmdConfig *core.CommandConfig, cmds []string, sep string, out io.Writer) error {
	if cf.Image == "" && cf.Name == "" {
		return fmt.Errorf("image or name must be set")
	}

	// find current docker process
	r, err := c.ContainerAdapter.Find(ctx, *config.ContainerFindOptionsFromConfig(cf))
	if err != nil {
		return err
	}

	// run all commands
	opts := &core.ContainerExecOptions{
		ID:     r.ID,
		Cmd:    prepareExecCmds(cmdConfig, cmds, sep),
		Config: cmdConfig,
	}

	return c.ContainerAdapter.Exec(ctx, *opts)
}

func execPrintConfig(cmdConfig *core.CommandConfig, out io.Writer) error {
	enc := config.NewEncoder(out)

	return enc.Encode(cmdConfig)
}

func prepareExecCmds(cmdConfig *core.CommandConfig, cmds []string, sep string) []string {
	prefix := []string{cmdConfig.Shell, "-c"}

	// merge with options
	if cmdConfig.Options != "" {
		cmds = append([]string{cmds[0], cmdConfig.Options}, cmds[1:]...)
	}

	return append(prefix, strings.Join(cmds, sep))
}

func parseExecConfig(c *ExecCommand) {
	withParsers(
		parseConfigPathOption(c.Options.ConfigPath),
		parseStdinOptions(c.Options.Stdin, c.Options.NoStdin),
		parseTtyOptions(c.Options.Tty, c.Options.NoTty),
		parseEnvsOption(&c.Config.Defaults.Envs, c.Options.Envs),
		parseShellOption(c.Options.Shell),
		parseUserOption(c.Options.User),
		parseWorkDirOption(c.Options.WorkDir, c.Config.Path),
	)(c.Config)
}

func NewExecCmd(c core.ContainerAdapter, cfs core.ConfigService, cf *core.Config, opts *core.CommandExecOptions, out io.Writer) *ExecCommand {
	return &ExecCommand{
		ContainerAdapter: c,
		ConfigService:    cfs,
		Config:           cf,
		Options:          opts,
		Out:              out,
	}
}
