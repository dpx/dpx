package console

import (
	"errors"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/services/config"
	"github.com/urfave/cli/v2"
)

func withActions(actions ...cli.ActionFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		for _, action := range actions {
			if err := action(c); err != nil {
				return err
			}
		}

		return nil
	}
}

func withValidationAction(cf *core.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		if cf.Path == "" {
			return errors.New(config.FileName + " is required")
		}

		if cf.Name == "" && cf.Compose == "" {
			return errors.New("container name or docker-compose service name is required")
		}

		return nil
	}
}

func withArgsAction(dst *[]string) cli.ActionFunc {
	return func(c *cli.Context) error {
		*dst = c.Args().Slice()

		return nil
	}
}

func withFirstArgAction(dst *string) cli.ActionFunc {
	return func(c *cli.Context) error {
		*dst = c.Args().First()

		return nil
	}
}

type stringSlice struct {
	Src *cli.StringSlice
	Dst *[]string
}

func withStringSliceAction(flags []stringSlice) cli.ActionFunc {
	return func(c *cli.Context) error {
		for _, v := range flags {
			*v.Dst = v.Src.Value()
		}

		return nil
	}
}

func withExecuteAction(cmd core.Command) cli.ActionFunc {
	return func(c *cli.Context) error {
		return cmd.Execute(c.Context)
	}
}
