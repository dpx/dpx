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
)

func TestCommand_Link(t *testing.T) {
	ctx := context.Background()
	cfs := mock.NewConfigService()
	c := mock.NewContainerAdapter(
		mock.ContainerAdapterWithFind("1234"),
	)
	c.ExecFn = func(ctx context.Context, options core.ContainerExecOptions) error {
		cmd := []string{
			"sh", "-c", "command -v ruby >/dev/null 2>&1",
		}

		assert.Equal(t, options.ID, "1234")
		assert.Equal(t, options.Cmd, cmd)

		return nil
	}

	cf := config.New()
	buf := &bytes.Buffer{}
	opts := &core.CommandLinkOptions{
		Cmds: []string{"ruby"},
	}

	cmd := command.NewLinkCmd(c, cfs, cf, opts, buf)
	cmd.Execute(ctx)

	out := []string{
		"LINK: ruby (.dpx/bin/ruby)",
	}

	assert.Equal(t, strings.Join(out, "\n")+"\n", buf.String())
}

func TestCommand_Link_Multiple(t *testing.T) {
	ctx := context.Background()
	cfs := mock.NewConfigService()
	c := mock.NewContainerAdapter(
		mock.ContainerAdapterWithFind("1234"),
	)
	c.ExecFn = func(ctx context.Context, options core.ContainerExecOptions) error {
		assert.Equal(t, options.ID, "1234")

		return nil
	}

	cf := config.New()
	buf := &bytes.Buffer{}
	opts := &core.CommandLinkOptions{
		Cmds: []string{"ruby", "gem"},
	}

	cmd := command.NewLinkCmd(c, cfs, cf, opts, buf)
	cmd.Execute(ctx)

	out := []string{
		"LINK: ruby (.dpx/bin/ruby)",
		"LINK: gem (.dpx/bin/gem)",
	}

	assert.Equal(t, strings.Join(out, "\n")+"\n", buf.String())
}

func TestCommand_Link_All(t *testing.T) {
	ctx := context.Background()
	cfs := mock.NewConfigService()
	c := mock.NewContainerAdapter(
		mock.ContainerAdapterWithFind("1234"),
	)
	c.ExecFn = func(ctx context.Context, options core.ContainerExecOptions) error {
		assert.Equal(t, options.ID, "1234")

		return nil
	}

	cf := config.New()
	cf.Commands = map[string]core.CommandConfig{
		"ruby": {},
		"rust": {},
	}

	buf := &bytes.Buffer{}
	opts := &core.CommandLinkOptions{
		All: true,
	}

	cmd := command.NewLinkCmd(c, cfs, cf, opts, buf)
	cmd.Execute(ctx)

	s := buf.String()
	assert.True(t, strings.Contains(s, "LINK: ruby"))
	assert.True(t, strings.Contains(s, "LINK: rust"))
}

func TestCommand_Link_WithoutCommand(t *testing.T) {
	ctx := context.Background()
	cfs := mock.NewConfigService()
	c := mock.NewContainerAdapter(
		mock.ContainerAdapterWithFind("1234"),
	)

	cf := config.New()
	buf := &bytes.Buffer{}
	opts := &core.CommandLinkOptions{
		Cmds: []string{},
	}

	cmd := command.NewLinkCmd(c, cfs, cf, opts, buf)
	err := cmd.Execute(ctx)

	assert.Error(t, err, command.ErrorMissingCommand)
}
