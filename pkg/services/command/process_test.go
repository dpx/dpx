package command_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/mock"
	"github.com/dpx/dpx/pkg/mock/assert"
	"github.com/dpx/dpx/pkg/services/command"
	"github.com/dpx/dpx/pkg/services/config"
)

func TestCommand_Process(t *testing.T) {
	ctx := context.Background()
	c := mock.NewContainerAdapter(
		mock.ContainerAdapterWithFind("1234567890123456789"),
		mock.ContainerAdapterWithStart(),
	)
	opts := &core.CommandProcessOptions{}
	cf := &core.Config{}
	buf := &bytes.Buffer{}

	cmd := command.NewProcessCmd(c, cf, opts, buf)
	cmd.Execute(ctx)

	out := "123456789012\n"

	assert.Equal(t, out, buf.String())
	assert.Equal(t, config.ProcessIdLength+1, len(out))
}
