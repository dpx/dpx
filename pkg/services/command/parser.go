package command

import (
	"path"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/services/config"
	"github.com/dpx/dpx/pkg/utils"
)

type parseFn func(cf *core.Config)

func withParsers(fns ...parseFn) parseFn {
	return func(cf *core.Config) {
		for _, f := range fns {
			f(cf)
		}
	}
}

func parseConfigPathOption(path string) parseFn {
	return func(cf *core.Config) {
		if path != "" {
			cfg, err := config.NewFromFile(path)

			if err == nil {
				// override
				*cf = *cfg
			}
		}
	}
}

func parseShellOption(sh string) parseFn {
	return func(cf *core.Config) {
		if sh != "" {
			cf.Defaults.Shell = sh
		}
	}
}

func parseStringOption(dst *string, src string) parseFn {
	return func(cf *core.Config) {
		if src != "" {
			*dst = src
		}
	}
}

var parseImageOption = parseStringOption

func parseNameOption(dst *string, name, image, path string) parseFn {
	return func(cf *core.Config) {
		if name == "" {
			name = utils.GetContainerName(image, resolveWorkingDirPath(cf.Path))
		}

		*dst = name
	}
}

func parseStringSliceOption(dst *[]string, src []string) parseFn {
	return func(cf *core.Config) {
		*dst = append(*dst, src...)
	}
}

var parsePortsOption = parseStringSliceOption
var parseEnvsOption = parseStringSliceOption
var parseVolumesOption = parseStringSliceOption

func parseTtyOptions(tty bool, noTty bool) parseFn {
	return func(cf *core.Config) {
		if tty {
			cf.Defaults.TTY = utils.Bool(true)
		} else if noTty {
			cf.Defaults.TTY = utils.Bool(false)
		} else if cf.Defaults.TTY == nil {
			isStdin := utils.IsStdIn()
			isTTY := utils.IsTTY()

			cf.Defaults.TTY = utils.Bool(isTTY && !isStdin)
		}
	}
}

func parseStdinOptions(stdin bool, noStdin bool) parseFn {
	return func(cf *core.Config) {
		if stdin {
			cf.Defaults.Stdin = utils.Bool(true)
		} else if noStdin {
			cf.Defaults.Stdin = utils.Bool(false)
		}
	}
}

func parseWorkDirOption(pwd string, configPath string) parseFn {
	return func(cf *core.Config) {
		if pwd != "" {
			cf.Defaults.WorkingDir = utils.String(pwd)
		} else if cf.Defaults.WorkingDir == nil && configPath != "" {
			cf.Defaults.WorkingDir = utils.String(resolveWorkingDirPath(configPath))
		}
	}
}

func parseUserOption(user string) parseFn {
	return func(cf *core.Config) {
		if user != "" {
			cf.Defaults.User = user
		}
	}
}

func resolveWorkingDirPath(configPath string) string {
	return path.Dir(configPath)
}
