package command_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/mock"
	"github.com/dpx/dpx/pkg/mock/assert"
	"github.com/dpx/dpx/pkg/services/command"
)

func TestCommand_Start_ContainerExists(t *testing.T) {
	ctx := context.Background()
	c := mock.NewContainerAdapter(
		mock.ContainerAdapterWithFind("12345678901234"),
		mock.ContainerAdapterWithStart(),
	)
	opts := &core.CommandStartOptions{}
	cf := &core.Config{}
	buf := &bytes.Buffer{}

	cmd := command.NewStartCmd(c, cf, opts, buf)
	cmd.Execute(ctx)

	out := "Starting a container 123456789012\n"
	assert.Equal(t, out, buf.String())
}

func TestCommand_Start_ContainerNotExists(t *testing.T) {
	ctx := context.Background()
	c := mock.NewContainerAdapter(
		mock.ContainerAdapterWithFind(""),
		mock.ContainerAdapterWithStart(),
	)
	c.CreateFn = func(ctx context.Context, options core.ContainerCreateOptions) (core.Container, error) {
		return core.Container{ID: "12345678901234"}, nil
	}

	opts := &core.CommandStartOptions{}
	cf := &core.Config{}
	buf := &bytes.Buffer{}

	cmd := command.NewStartCmd(c, cf, opts, buf)
	cmd.Execute(ctx)

	out := "Starting a container 123456789012\n"
	assert.Equal(t, out, buf.String())
}
