package stream

import (
	"io"

	"github.com/docker/docker/pkg/term"
)

type Stream struct {
	In  *in
	Out *out
	Err *err
}

type base struct {
	fd         uintptr
	isTerminal bool
	state      *term.State
}

type in struct {
	base
	In io.ReadCloser
}

type out struct {
	base
	Out io.Writer
}

type err struct {
	base
	Err io.Writer
}

func New() *Stream {
	in, out, err := term.StdStreams()

	return &Stream{
		newIn(in),
		newOut(out),
		newErr(err),
	}
}

func NewFromStd(in io.ReadCloser, out io.Writer, err io.Writer) *Stream {
	return &Stream{
		newIn(in),
		newOut(out),
		newErr(err),
	}
}

func (s *Stream) RestoreTerminal() {
	s.In.Restore()
	s.Out.Restore()
}

func (b *base) Restore() {
	if b.state != nil {
		term.RestoreTerminal(b.fd, b.state)
	}
}

func (b *base) GetSize() (*term.Winsize, error) {
	return term.GetWinsize(b.fd)
}

func (o *out) Write(b []byte) (int, error) {
	return o.Out.Write(b)
}

func (o *err) Write(b []byte) (int, error) {
	return o.Err.Write(b)
}

func newIn(stdin io.ReadCloser) *in {
	fd, isTerm := term.GetFdInfo(stdin)

	return &in{base: base{fd: fd, isTerminal: isTerm}, In: stdin}
}

func newOut(stdout io.Writer) *out {
	fd, isTerm := term.GetFdInfo(stdout)

	return &out{base: base{fd: fd, isTerminal: isTerm}, Out: stdout}
}

func newErr(stderr io.Writer) *err {
	fd, isTerm := term.GetFdInfo(stderr)

	return &err{base: base{fd: fd, isTerminal: isTerm}, Err: stderr}
}
