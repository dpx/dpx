package console

import (
	"github.com/dpx/dpx/pkg/core"
	"github.com/urfave/cli/v2"
)

func withFlags(flags ...cli.Flag) []cli.Flag {
	return flags
}

func withEnvsFlag(v *cli.StringSlice) cli.Flag {
	return &cli.StringSliceFlag{
		Name:        "env",
		Aliases:     []string{"e"},
		Usage:       "Set environment variables `ENV`",
		Destination: v,
	}
}

func withPortsFlag(v *cli.StringSlice) cli.Flag {
	return &cli.StringSliceFlag{
		Name:        "port",
		Aliases:     []string{"p"},
		Usage:       "Bind port `PORT`",
		Destination: v,
	}
}

func withVolumesFlag(v *cli.StringSlice) cli.Flag {
	return &cli.StringSliceFlag{
		Name:        "volumes",
		Aliases:     []string{"v"},
		Usage:       "Volume `VOLUME`",
		Destination: v,
	}
}

func withShellFlag(v *string) cli.Flag {
	return &cli.StringFlag{
		Name:        "shell",
		Usage:       "Set default shell",
		Value:       core.DefaultShell,
		Destination: v,
	}
}
