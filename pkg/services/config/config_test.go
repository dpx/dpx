package config_test

import (
	"testing"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/mock/assert"
	"github.com/dpx/dpx/pkg/services/config"
	"github.com/dpx/dpx/pkg/utils"
)

func TestConfig_New(t *testing.T) {
	a := defaultConfig()
	c := config.New()

	assertConfig(t, a.Defaults, c.Defaults)
}

func TestConfig_NewFromReader_FalsyValues(t *testing.T) {
	a := defaultConfig()
	a.Defaults.WorkingDir = utils.String("")

	c, _ := config.NewFromString(`
defaults:
  stdin: true
  tty: false
  workdir: ""
  envs:
`)

	assertConfig(t, a.Defaults, c.Defaults)
}

func TestConfig_CommandConfig(t *testing.T) {
	a := defaultConfig()
	c, _ := config.NewFromString(`
commands:
  go:
`)
	cfs := config.NewConfigService()

	assertConfig(t, a.Defaults, cfs.Command("go", c))
}

func TestConfig_CommandConfigInheritsDefault(t *testing.T) {
	a := map[string]*core.CommandConfig{
		"go": {
			Stdin: utils.Bool(false),
			TTY:   utils.Bool(false),
			Envs: []string{
				"TEST=true",
			},
			WorkingDir: utils.String("/app"),
			User:       "dpx",
		},
	}

	c, _ := config.NewFromString(`
defaults:
  stdin: false
  tty: false
  envs:
    - TEST=true
  workdir: /app
  user: dpx
commands:
  go:
`)

	cfs := config.NewConfigService()

	assertConfig(t, a["go"], cfs.Command("go", c))
}

func TestConfig_CommandConfigOverridesDefault(t *testing.T) {
	a := map[string]*core.CommandConfig{
		"go": {
			Stdin: utils.Bool(true),
			TTY:   utils.Bool(true),
			Envs: []string{
				"TEST=true",
				"NAME=dpx",
			},
			WorkingDir: utils.String("/app2"),
			User:       "root",
		},
	}

	c, _ := config.NewFromString(`
defaults:
  stdin: false
  tty: false
  envs:
    - TEST=true
  workdir: /app
  user: dpx
commands:
  go:
    stdin: true
    tty: true
    envs:
      - NAME=dpx
    workdir: /app2
    user: root
`)
	cfs := config.NewConfigService()

	assertConfig(t, a["go"], cfs.Command("go", c))
}

func assertConfig(t *testing.T, e *core.CommandConfig, c *core.CommandConfig) {
	if e.Stdin != nil {
		t.Run("stdin", func(t *testing.T) {
			assert.Bool(t, e.Stdin, c.Stdin)
		})
	}

	if e.TTY != nil {
		t.Run("tty", func(t *testing.T) {
			assert.Bool(t, e.TTY, c.TTY)
		})
	}

	if e.WorkingDir != nil {
		t.Run("workdir", func(t *testing.T) {
			assert.Equal(t, *e.WorkingDir, *c.WorkingDir)
		})
	}

	t.Run("envs", func(t *testing.T) {
		assert.Equal(t, e.Envs, c.Envs)
	})
}

func defaultConfig() *core.Config {
	c := config.New()
	c.Defaults = &core.CommandConfig{
		TTY:        utils.Bool(false),
		Stdin:      utils.Bool(true),
		WorkingDir: utils.String(""),
	}

	return c
}
