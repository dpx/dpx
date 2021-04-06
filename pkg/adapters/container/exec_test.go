package container_test

import (
	"bufio"
	"bytes"
	"context"
	"net"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/dpx/dpx/pkg/adapters/container"
	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/mock"
	"github.com/dpx/dpx/pkg/mock/assert"
	"github.com/dpx/dpx/pkg/utils"
)

func TestContainer_Exec(t *testing.T) {
	d := newExecMock()
	d.ContainerExecCreateFn = func(ctx context.Context, contID string, config types.ExecConfig) (types.IDResponse, error) {
		assert.Equal(t, contID, "1111")

		return types.IDResponse{}, nil
	}
	d.ContainerExecAttachFn = func(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error) {
		return newExecResponse("ok")
	}

	c := container.New(d)
	o, _ := c.ExecWithOutput(context.Background(), core.ContainerExecOptions{
		ID: "1111",
		Config: &core.CommandConfig{
			TTY:        utils.Bool(true),
			Stdin:      utils.Bool(true),
			WorkingDir: utils.String("/app"),
		},
		Output: true,
	})

	assert.Equal(t, "ok", string(o))
}

func newExecMock() *mock.Docker {
	d := &mock.Docker{}
	d.ContainerExecCreateFn = func(ctx context.Context, contID string, config types.ExecConfig) (types.IDResponse, error) {
		return types.IDResponse{}, nil
	}
	d.ContainerExecAttachFn = func(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error) {
		return types.HijackedResponse{}, nil
	}
	d.ContainerExecInspectFn = func(ctx context.Context, execID string) (types.ContainerExecInspect, error) {
		return types.ContainerExecInspect{ExitCode: 0}, nil
	}

	return d
}

func newExecResponse(data string) (types.HijackedResponse, error) {
	buf := new(bytes.Buffer)
	w := stdcopy.NewStdWriter(buf, stdcopy.Stdout)
	w.Write([]byte(data))

	r := bufio.NewReader(buf)

	_, c := net.Pipe()

	return types.HijackedResponse{Reader: r, Conn: c}, nil
}
