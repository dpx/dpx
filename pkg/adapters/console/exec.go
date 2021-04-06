package console

import (
	"github.com/dpx/dpx/pkg/services/command"
	"github.com/dpx/dpx/pkg/services/config"
	"github.com/urfave/cli/v2"
)

// ExecCmd setup `exec` command and flags.
func NewExecCmd(cmd *command.ExecCommand) *cli.Command {
	slices := []stringSlice{
		{cli.NewStringSlice(), &cmd.Options.Envs},
	}

	return &cli.Command{
		Name:    "exec",
		Aliases: []string{"x"},
		Usage:   "Execute a command in container",
		Flags:   withExecFlags(cmd, slices),
		Action: withActions(
			withFirstArgAction(&cmd.Options.Cmd),
			withArgsAction(&cmd.Options.CmdArgs),
			withStringSliceAction(slices),
			withExecuteAction(cmd),
		),
	}
}

func withExecFlags(cmd *command.ExecCommand, slices []stringSlice) []cli.Flag {
	return withFlags(
		withShellFlag(&cmd.Options.Shell),
		withEnvsFlag(slices[0].Src),
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Usage:       "Specify path to " + config.FileName + " config file",
			Destination: &cmd.Options.ConfigPath,
		},
		&cli.StringFlag{
			Name:        "user",
			Aliases:     []string{"u"},
			Usage:       "Set username or UID",
			Destination: &cmd.Options.User,
		},
		&cli.StringFlag{
			Name:        "workdir",
			Aliases:     []string{"w"},
			Usage:       "Set working directory inside the container (by default, it'll use `pwd`)",
			Destination: &cmd.Options.WorkDir,
		},
		&cli.StringFlag{
			Name:        "compose",
			Usage:       "Use docker-compose service name",
			Destination: &cmd.Options.Compose,
		},
		&cli.BoolFlag{
			Name:        "print",
			Usage:       "Print current config",
			Destination: &cmd.Options.Print,
		},
		&cli.BoolFlag{
			Name:        "all",
			Aliases:     []string{"a"},
			Usage:       "Run all commands defined in " + config.FileName,
			Destination: &cmd.Options.All,
		},
		&cli.BoolFlag{
			Name:        "save",
			Aliases:     []string{"s"},
			Usage:       "Save command to " + config.FileName + " file",
			Destination: &cmd.Options.Save,
		},
		&cli.BoolFlag{
			Name:        "tty",
			Aliases:     []string{"t"},
			Usage:       "Allocate a pseudo-TTY",
			Destination: &cmd.Options.Tty,
		},
		&cli.BoolFlag{
			Name:        "no-tty",
			Aliases:     []string{"T"},
			Usage:       "Disable pseudo-TTY",
			Destination: &cmd.Options.NoTty,
		},
		&cli.BoolFlag{
			Name:        "stdin",
			Aliases:     []string{"i"},
			Usage:       "Attach to stdin",
			Destination: &cmd.Options.Stdin,
		},
		&cli.BoolFlag{
			Name:        "no-stdin",
			Aliases:     []string{"I"},
			Usage:       "Disable stdin",
			Destination: &cmd.Options.NoStdin,
		},
	)
}
