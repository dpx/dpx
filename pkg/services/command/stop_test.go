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

func TestCommand_Stop_ContainerExists(t *testing.T) {
	ctx := context.Background()
	id := "123456789012"
	c := mock.NewContainerAdapter(
		mock.ContainerAdapterWithFind(id),
	)
	c.StopFn = func(ctx context.Context, ID string) error {
		assert.Equal(t, id, ID)

		return nil
	}

	opts := &core.CommandStopOptions{}
	cf := &core.Config{}
	buf := &bytes.Buffer{}

	cmd := command.NewStopCmd(c, cf, opts, buf)
	cmd.Execute(ctx)

	out := "Stopping a container 123456789012\n"
	assert.Equal(t, out, buf.String())
}

func TestCommand_Stop_ContainerNotExists(t *testing.T) {
	ctx := context.Background()
	c := mock.NewContainerAdapter()
	c.FindFn = func(ctx context.Context, options core.ContainerFindOptions) (core.Container, error) {
		return core.Container{}, core.ErrContainerNotFound
	}

	opts := &core.CommandStopOptions{}
	cf := &core.Config{}
	buf := &bytes.Buffer{}

	cmd := command.NewStopCmd(c, cf, opts, buf)
	err := cmd.Execute(ctx)

	assert.Equal(t, err, core.ErrContainerNotFound)
}
