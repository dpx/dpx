package command_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/mock"
	"github.com/dpx/dpx/pkg/mock/assert"
	"github.com/dpx/dpx/pkg/services/command"
	"github.com/dpx/dpx/pkg/services/config"
	"github.com/dpx/dpx/pkg/utils"
)

// dpx exec go
func TestCommand_Exec(t *testing.T) {
	ctx := context.Background()
	cfs := mock.NewConfigService()
	c := mock.NewContainerAdapter(
		mock.ContainerAdapterWithFind("1234"),
	)
	c.ExecFn = func(ctx context.Context, options core.ContainerExecOptions) error {
		cmd := []string{
			"sh", "-c", "go -v",
		}

		assert.Equal(t, options.ID, "1234")
		assert.Equal(t, options.Cmd, cmd)

		return nil
	}

	cf := config.New()

	buf := &bytes.Buffer{}
	opts := &core.CommandExecOptions{
		Cmd:     "go",
		CmdArgs: []string{"go", "-v"},
	}

	cmd := command.NewExecCmd(c, cfs, cf, opts, buf)
	cmd.Execute(ctx)

	out := []string{}

	assert.Equal(t, strings.Join(out, "\n"), buf.String())
}

func TestCommand_Exec_WithOptions(t *testing.T) {
	ctx := context.Background()
	cfs := mock.NewConfigService()
	c := mock.NewContainerAdapter(
		mock.ContainerAdapterWithFind("1234"),
	)
	c.ExecFn = func(ctx context.Context, options core.ContainerExecOptions) error {
		cmd := []string{
			"sh", "-c", "go -v",
		}

		assert.Equal(t, options.ID, "1234")
		assert.Equal(t, options.Cmd, cmd)

		return nil
	}

	cf := config.New()
	cf.Commands = map[string]core.CommandConfig{
		"go": {Options: "-v"},
	}

	buf := &bytes.Buffer{}
	opts := &core.CommandExecOptions{
		Cmd:     "go",
		CmdArgs: []string{"go"},
	}

	cmd := command.NewExecCmd(c, cfs, cf, opts, buf)
	cmd.Execute(ctx)

	out := []string{}

	assert.Equal(t, strings.Join(out, "\n"), buf.String())
}

func TestCommand_Exec_WithExitCode(t *testing.T) {
	ctx := context.Background()
	cfs := mock.NewConfigService()
	c := mock.NewContainerAdapter(
		mock.ContainerAdapterWithFind("1234"),
	)
	c.ExecFn = func(ctx context.Context, options core.ContainerExecOptions) error {
		return &core.ContainerExecErr{
			Code: 127,
		}
	}

	cf := config.New()
	cf.Image = "golang"

	buf := &bytes.Buffer{}
	opts := &core.CommandExecOptions{
		Cmd:     "go",
		CmdArgs: []string{"go", "-v"},
	}

	cmd := command.NewExecCmd(c, cfs, cf, opts, buf)
	err := cmd.Execute(ctx)
	out := []string{}

	assert.Equal(t, strings.Join(out, "\n"), buf.String())
	assert.Equal(t, 127, err.(*core.ContainerExecErr).Code)
}

// dpx exec --all
func TestCommand_Exec_All(t *testing.T) {
	ctx := context.Background()
	cfs := mock.NewConfigService()
	c := mock.NewContainerAdapter(
		mock.ContainerAdapterWithFind("1234"),
	)
	c.ExecFn = func(ctx context.Context, options core.ContainerExecOptions) error {
		cmd := []string{
			"sh", "-c", "go -h && go -v && go help mod",
		}

		assert.Equal(t, options.ID, "1234")
		assert.Equal(t, options.Cmd, cmd)

		return nil
	}

	cf := config.New()
	cf.Path = "dpx.yml"
	cf.Runs = []string{
		"go -h",
		"go -v",
		"go help mod",
	}

	buf := &bytes.Buffer{}
	opts := &core.CommandExecOptions{
		All: true,
	}

	cmd := command.NewExecCmd(c, cfs, cf, opts, buf)
	cmd.Execute(ctx)

	out := []string{}

	assert.Equal(t, strings.Join(out, "\n"), buf.String())
}

// dpx --print go
func TestCommand_Exec_Print(t *testing.T) {
	ctx := context.Background()
	cfs := mock.NewConfigService()
	c := mock.NewContainerAdapter(
		mock.ContainerAdapterWithFind("1234"),
	)
	cf := &core.Config{
		Path: "dpx.yml",
		Defaults: &core.CommandConfig{
			TTY:  utils.Bool(true),
			User: "root",
			Envs: []string{
				"TEST=true",
			},
		},
		Commands: map[string]core.CommandConfig{
			"go": {
				User: "guest",
				Envs: []string{
					"USER=dpx",
				},
			},
		},
	}

	buf := &bytes.Buffer{}
	opts := &core.CommandExecOptions{
		Print: true,
		Cmd:   "go",
	}

	cmd := command.NewExecCmd(c, cfs, cf, opts, buf)
	cmd.Execute(ctx)

	out := `tty: true
envs:
  - TEST=true
  - USER=dpx
user: guest
workdir: .
`

	assert.Equal(t, out, buf.String())
}
