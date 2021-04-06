package console

import (
	"github.com/dpx/dpx/pkg/services/command"
	"github.com/urfave/cli/v2"
)

// NewStopCmd stops a running container.
//
// dpx stop
func NewStopCmd(cmd *command.StopCommand) *cli.Command {
	return &cli.Command{
		Name:  "stop",
		Usage: "Stop a running container",
		Action: withActions(
			withValidationAction(cmd.Config),
			withExecuteAction(cmd),
		),
	}
}
