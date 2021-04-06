package console

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/services/command"
	"github.com/dpx/dpx/pkg/services/config"
	"github.com/urfave/cli/v2"
)

// InitCmd setup `init` command and flags.
// Init initializes and setup dpx.yml file.
//
// dpx init golang
func NewInitCmd(cmd *command.InitCommand, startCmd *cli.Command) *cli.Command {
	return &cli.Command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "Setup docker image and create dpx.yml config file",
		UsageText: `dpx init [options] IMAGE

EXAMPLE:
   dpx init golang             (use latest golang image)
   dpx init golang:1.14-alpine (use golang:1.14-alpine)`,
		Flags: withFlags(
			&cli.StringFlag{
				Name:        "name",
				Usage:       "Set container name",
				Destination: &cmd.Options.Name,
			},
		),
		Action: withActions(
			withFirstArgAction(&cmd.Options.Image),
			withInitConfigAction(cmd.Config, &cmd.Options.Skip, cmd.Out),
			withExecuteAction(cmd),
			startCmd.Action,
		),
	}
}

func withInitConfigAction(cf *core.Config, skip *bool, out io.Writer) cli.ActionFunc {
	return func(c *cli.Context) error {
		// config doesn't exist
		if cf.Path == "" {
			return nil
		}

		// config already exists
		fmt.Fprintf(out, "%s already exists. Do you want to continue? (y - overwrite, s - skip, other - abort)\n", config.FileName)
		r, _, err := bufio.NewReader(os.Stdin).ReadRune()
		if err != nil {
			return err
		}

		switch unicode.ToUpper(r) {
		case 'Y':
			return nil
		case 'S':
			*skip = true
			return nil
		}

		return fmt.Errorf("Abort")
	}
}
