package container

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/term"
	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/services/stream"
)

func (c *ContainerAdapter) Exec(ctx context.Context, opts core.ContainerExecOptions) error {
	_, err := exec(ctx, c, opts)

	return err
}

// ExecWithOutput executes command and returns output
func (c *ContainerAdapter) ExecWithOutput(ctx context.Context, opts core.ContainerExecOptions) ([]byte, error) {
	s, err := exec(ctx, c, opts)

	return s.Out.Out.(*bytes.Buffer).Bytes(), err
}

func (c *ContainerAdapter) Inspect(ctx context.Context, ID string) (core.ContainerInspectResult, error) {
	r, err := c.Docker.ContainerExecInspect(ctx, ID)

	return core.ContainerInspectResult{
		ExitCode: r.ExitCode,
		Running:  r.Running,
	}, err
}

func exec(ctx context.Context, c *ContainerAdapter, opts core.ContainerExecOptions) (*stream.Stream, error) {
	config := opts.Config

	exec, err := c.Docker.ContainerExecCreate(ctx, opts.ID, types.ExecConfig{
		// Pseudo TTY
		Tty: *config.TTY,

		// We have to attach these stdin, out and error
		// otherwise we can't read stdout and stderr
		AttachStdin:  *config.Stdin,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          opts.Cmd,
		Env:          config.Envs,
		User:         config.User,
		WorkingDir:   *config.WorkingDir,
	})

	if err != nil {
		panic(err)
	}

	res, err := c.Docker.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	defer res.Close()

	if err != nil {
		panic(err)
	}

	stream := stream.New()
	if opts.Output {
		stream.Out.Out = new(bytes.Buffer)
	}

	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		errCh <- handleStream(ctx, stream, stream.Out.Out, res)
	}()

	if err := <-errCh; err != nil {
		fmt.Printf("Error hijack: %s\n", err)
		return stream, err
	}

	return stream, getExitCode(ctx, c, exec.ID)
}

func handleStream(ctx context.Context, stream *stream.Stream, out io.Writer, res types.HijackedResponse) error {
	// read the output
	outDone := handleOutputStream(out, res.Reader)

	// Write to docker container
	inDone := handleInputStream(stream, res)

	select {
	case err := <-outDone:
		if err != nil {
			return err
		}
		break

	case <-inDone:
		select {
		case err := <-outDone:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}

	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}

	return nil
}

func handleInputStream(stream *stream.Stream, resp types.HijackedResponse) <-chan struct{} {
	var in = make(chan struct{})
	var defaultEscapeKeys = []byte{16, 17}

	go func() {
		ins := ioutils.NewReadCloserWrapper(term.NewEscapeProxy(stream.In.In, defaultEscapeKeys), stream.In.In.Close)
		io.Copy(resp.Conn, ins)

		// We have to close read writer here.
		// otherwise it'll wait for stdin forever.
		if err := resp.CloseWrite(); err != nil {
			fmt.Printf("Couldn't send EOF: %s", err)
		}

		close(in)
	}()

	return in
}

func handleOutputStream(w io.Writer, r io.Reader) <-chan error {
	var errBuf bytes.Buffer
	outDone := make(chan error)

	go func() {
		// Write to bytes.Buffer will block if the input is infinite
		// e.g. `yes | dpx -i cat`
		// if tty use io.Copy else StdCopy
		_, err := stdcopy.StdCopy(w, &errBuf, r)

		if err != nil {
			fmt.Printf("err: %s\n", errBuf.String())
		}
		if errBuf.Len() > 0 {
			fmt.Printf("err: %s\n", errBuf.String())
		}

		outDone <- err
	}()

	return outDone
}

func getExitCode(ctx context.Context, c *ContainerAdapter, ID string) error {
	r, err := c.Inspect(ctx, ID)
	if err != nil {
		return err
	}

	// Propagate exit code
	if r.ExitCode != 0 {
		return &core.ContainerExecErr{Code: r.ExitCode}
	}

	return nil
}

// -e COLUMNS="`tput cols`" -e LINES="`tput lines`" to fix tty size
// when exec /bin/sh
func resizeTty(ctx context.Context, execID string) {

	// tty
	fd := uintptr(os.Stdout.Fd())
	if term.IsTerminal(fd) {
		state, err := term.MakeRaw(fd)
		if err != nil {
			panic(err)
		}
		defer term.RestoreTerminal(fd, state)

		size, err := term.GetWinsize(fd)
		if err != nil {
			panic(err)
		}
		resizeOptions := types.ResizeOptions{
			Height: uint(size.Height),
			Width:  uint(size.Width),
		}

		cli, _ := client.NewClientWithOpts(client.FromEnv)
		if err := cli.ContainerExecResize(ctx, execID, resizeOptions); err != nil {
			panic(err)
		}
	}
}
