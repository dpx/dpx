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
)

func TestCommand_Init_LocalImageNotExists(t *testing.T) {
	ctx := context.Background()
	img := "go"
	i := mock.NewImageAdapter(
		mock.ImageAdapterWithPull(img),
	)
	c := mock.NewContainerAdapter(
		mock.ContainerAdapterWithCreate("1234"),
	)
	cfs := mock.NewConfigService()
	cf := &core.Config{}
	opts := &core.CommandInitOptions{
		Image: "go",
	}

	o := &bytes.Buffer{}
	cmd := command.NewInitCmd(i, c, cfs, opts, cf, o)
	cmd.Execute(ctx)

	out := []string{
		"Cannot find go image locally",
		"Creating a container dpx-go",
	}
	assert.Equal(t, strings.Join(out, "\n")+"\n", o.String())
}

func TestCommand_Init_LocalImageExists(t *testing.T) {
	ctx := context.Background()
	img := "go"
	i := mock.NewImageAdapter(
		mock.ImageAdapterWithFind(img),
	)
	c := mock.NewContainerAdapter()
	c.CreateFn = func(ctx context.Context, options core.ContainerCreateOptions) (core.Container, error) {
		assert.Equal(t, "go", options.Image)
		assert.Equal(t, "dpx-go", options.Name)

		return core.Container{ID: "1234"}, nil
	}
	cfs := mock.NewConfigService()
	cf := &core.Config{}
	opts := &core.CommandInitOptions{
		Image: "go",
	}

	o := &bytes.Buffer{}

	cmd := command.NewInitCmd(i, c, cfs, opts, cf, o)
	cmd.Execute(ctx)

	out := []string{
		"Creating a container dpx-go",
	}
	assert.Equal(t, strings.Join(out, "\n")+"\n", o.String())
}

func TestCommand_Init_LocalImageExistsWithName(t *testing.T) {
	ctx := context.Background()
	img := "go"
	i := mock.NewImageAdapter(
		mock.ImageAdapterWithFind(img),
	)
	c := mock.NewContainerAdapter()
	c.CreateFn = func(ctx context.Context, options core.ContainerCreateOptions) (core.Container, error) {
		assert.Equal(t, "go", options.Image)
		assert.Equal(t, "dpx-go", options.Name)

		return core.Container{ID: "1234"}, nil
	}

	cfs := mock.NewConfigService()
	cf := &core.Config{}
	opts := &core.CommandInitOptions{
		Image: "go",
		Name:  "dpx-go",
	}

	o := &bytes.Buffer{}

	cmd := command.NewInitCmd(i, c, cfs, opts, cf, o)
	cmd.Execute(ctx)

	out := []string{
		"Creating a container dpx-go",
	}
	assert.Equal(t, strings.Join(out, "\n")+"\n", o.String())
}
