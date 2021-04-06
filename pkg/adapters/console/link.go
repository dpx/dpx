package console

import (
	"github.com/dpx/dpx/pkg/services/command"
	"github.com/dpx/dpx/pkg/services/config"
	"github.com/urfave/cli/v2"
)

// NewLinkCmd setup `link` command and flags.
// links a binary inside container and creates a symlink file in .dpx/bin/
//
// dpx link go
func NewLinkCmd(cmd *command.LinkCommand) *cli.Command {
	return &cli.Command{
		Name:    "link",
		Aliases: []string{"l"},
		Usage:   "Link to an executable/binary inside container",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "all",
				Usage:       "Link all binaries defined in " + config.FileName,
				Destination: &cmd.Options.All,
			},
		},
		Action: withActions(
			withValidationAction(cmd.Config),
			withArgsAction(&cmd.Options.Cmds),
			withExecuteAction(cmd),
		),
	}
}
