package console

import (
	"github.com/dpx/dpx/pkg/services/command"
	"github.com/urfave/cli/v2"
)

// NewPathCmd setup `path` command and flags.
//
// Path setup $PATH environment variable.
//
// dpx path
func NewPathCmd(cmd *command.PathCommand) *cli.Command {
	return &cli.Command{
		Name:  "path",
		Usage: "Print $PATH variable",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "delete",
				Aliases:     []string{"d"},
				Usage:       "Remove .dpx/bin from $PATH variable",
				Destination: &cmd.Options.Delete,
			},
		},
		Action: withActions(
			withValidationAction(cmd.Config),
			withExecuteAction(cmd),
		),
	}
}
