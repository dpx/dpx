package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/services/config"
)

type PathCommand struct {
	Config  *core.Config
	Options *core.CommandPathOptions
	Out     io.Writer
}

func (c *PathCommand) Execute(ctx context.Context) error {
	p := c.Config.Path

	binDir := path.Join(path.Dir(p), config.BinDir)
	var r string

	if c.Options.Delete {
		r = c.Delete(os.Getenv("PATH"), binDir)
	} else {
		r = fmt.Sprintf(config.PathData, binDir)
	}

	fmt.Fprintln(c.Out, r)

	return nil
}

var maxDeleteLevel = 5

// Delete removes .dpx/bin path from PATH variable
func (c *PathCommand) Delete(pathEnv, binPath string) string {
	sep := string(os.PathListSeparator)
	reg := regexp.MustCompile("(?i)(?:^" + binPath + sep + "|" + sep + binPath + "$|" + sep + "+" + binPath + ")")
	r := reg.ReplaceAllLiteralString(pathEnv, "")

	return r
}

func NewPathCmd(cf *core.Config, opts *core.CommandPathOptions, out io.Writer) *PathCommand {
	return &PathCommand{
		Config:  cf,
		Options: opts,
		Out:     out,
	}
}
