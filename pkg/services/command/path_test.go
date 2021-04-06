package command_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/mock/assert"
	"github.com/dpx/dpx/pkg/services/command"
	"github.com/dpx/dpx/pkg/utils"
)

func TestCommand_Path(t *testing.T) {
	ctx := context.Background()
	opts := &core.CommandPathOptions{}
	cf := &core.Config{
		Path: "/app/dpx.yml",
	}
	buf := &bytes.Buffer{}

	cmd := command.NewPathCmd(cf, opts, buf)
	cmd.Execute(ctx)

	out := `PATH=/app/.dpx/bin:$PATH`

	assert.Equal(t, out+"\n", buf.String())
}

func TestCommand_Path_Delete(t *testing.T) {
	ctx := context.Background()
	opts := &core.CommandPathOptions{
		Delete: true,
	}
	cf := &core.Config{
		Path: "/app/dpx.yml",
	}
	buf := &bytes.Buffer{}

	tmp := os.Getenv("PATH")
	sep := string(os.PathListSeparator)

	// temporary change PATH variable
	os.Setenv("PATH", tmp+sep+"/app/.dpx/bin"+sep+"/usr/local/bin")

	cmd := command.NewPathCmd(cf, opts, buf)
	cmd.Execute(ctx)

	out := tmp + sep + "/usr/local/bin"

	assert.Equal(t, out+"\n", buf.String())

	// restore
	os.Setenv("PATH", tmp)
}

func TestCommand_Path_DeleteTests(t *testing.T) {
	opts := &core.CommandPathOptions{}
	buf := &bytes.Buffer{}
	cf := &core.Config{}

	binPath := "/app/.dpx/bin"
	testcases := map[string]string{
		utils.JoinPath(binPath, "/usr/local/bin"):                  "/usr/local/bin",
		utils.JoinPath("/usr/local/bin", binPath):                  "/usr/local/bin",
		utils.JoinPath("/usr/local/bin", binPath, "/usr/bin"):      "/usr/local/bin:/usr/bin",
		utils.JoinPath(binPath, "/usr/local/bin", binPath):         "/usr/local/bin",
		utils.JoinPath("/bin", binPath, "/usr/local/bin", binPath): "/bin:/usr/local/bin",
	}

	cmd := command.NewPathCmd(cf, opts, buf)
	for path, expected := range testcases {
		r := cmd.Delete(path, binPath)

		assert.Equal(t, expected, r)
	}
}
