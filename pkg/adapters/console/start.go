package console

import (
	"github.com/dpx/dpx/pkg/services/command"
	"github.com/urfave/cli/v2"
)

// NewStartCmd setup `start` command and flags.
//
// starts new container from existing image.
// (Assuming user has run `dpx i` previously)
//
// dpx start
func NewStartCmd(cmd *command.StartCommand) *cli.Command {
	slices := []stringSlice{
		{cli.NewStringSlice(), &cmd.Options.Envs},
		{cli.NewStringSlice(), &cmd.Options.Ports},
		{cli.NewStringSlice(), &cmd.Options.Volumes},
	}

	return &cli.Command{
		Name:    "start",
		Aliases: []string{"s"},
		Usage:   "Start container from config",
		Flags: withFlags(
			&cli.StringFlag{
				Name:        "name",
				Usage:       "Set container name",
				Destination: &cmd.Options.Name,
			},
			withEnvsFlag(slices[0].Src),
			withPortsFlag(slices[1].Src),
			withVolumesFlag(slices[2].Src),
		),
		Action: withActions(
			withValidationAction(cmd.Config),
			withStringSliceAction(slices),
			withExecuteAction(cmd),
		),
	}
}
