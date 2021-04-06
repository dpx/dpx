package console

import (
	"github.com/dpx/dpx/pkg/services/command"
	"github.com/urfave/cli/v2"
)

// ProcessCmd setup `ps` command and flags.
//
// prints current process id.
func NewProcessCmd(cmd *command.ProcessCommand) *cli.Command {
	return &cli.Command{
		Name:  "ps",
		Usage: "Print current container ID",
		Action: withActions(
			withValidationAction(cmd.Config),
			withExecuteAction(cmd),
		),
	}
}
