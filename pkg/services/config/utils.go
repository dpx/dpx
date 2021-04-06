package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/utils"
	"gopkg.in/yaml.v3"
)

func createAliasFile() (core.Path, error) {
	fpath := path.Join(BinDir, AliasFile)

	_, err := os.Stat(fpath)
	if os.IsExist(err) {
		// do nothing, if it already exists
		return core.Path(fpath), nil
	}

	err = os.MkdirAll(BinDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	ioutil.WriteFile(fpath, []byte(AliasData), AliasFileMode)

	return core.Path(fpath), nil
}

func createBinFile(bin string) (core.Path, error) {
	fpath := path.Join(BinDir, bin)

	err := os.MkdirAll(BinDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	if utils.FileExist(fpath) {
		return "", fmt.Errorf("file already exists: %s", fpath)
	}

	os.Symlink(AliasFile, fpath)

	return core.Path(fpath), nil
}

func merge(a core.CommandConfig, b core.CommandConfig) *core.CommandConfig {
	as := reflect.ValueOf(&a).Elem()
	bs := reflect.ValueOf(&b).Elem()

	n := as.NumField()
	for i := 0; i < n; i++ {
		av := as.Field(i)
		bv := bs.Field(i)

		switch av.Kind() {
		case reflect.Slice:
			bv.Set(reflect.AppendSlice(av, bv))
		case reflect.Ptr:
			if bv.IsNil() {
				bv.Set(av)
			}
		default:
			if bv.Interface() == nil || bv.IsZero() {
				bv.Set(av)
			}
		}
	}

	return &b
}

func NewEncoder(buf io.Writer) *yaml.Encoder {
	enc := yaml.NewEncoder(buf)
	enc.SetIndent(IndentLevel)

	return enc
}

func FormatProcessId(ID string) string {
	return ID[:ProcessIdLength]
}

func ContainerFindOptionsFromConfig(c *core.Config) *core.ContainerFindOptions {
	return &core.ContainerFindOptions{
		Image:   c.Image,
		Name:    c.Name,
		Compose: c.Compose,
		All:     true,
	}
}

func ContainerCreateOptionsFromConfig(c *core.Config) *core.ContainerCreateOptions {
	var workingDir string
	if c.Defaults != nil {
		workingDir = *c.Defaults.WorkingDir
	}

	return &core.ContainerCreateOptions{
		Image:      c.Image,
		Name:       c.Name,
		Ports:      c.Ports,
		Envs:       c.Envs,
		Volumes:    c.Volumes,
		WorkingDir: workingDir,
	}
}

// When `-w` option is not set. We'll assign default value from `pwd`.
//
// This is imporant. Otherwise commands that need dir context
// like `golint` won't work.
//
// For example. Code editor might execute `golint` from some dir
// (e.g. $ /app/pkg/config golint).
// 	if bin.WorkingDir == "." {
// 		bin.WorkingDir, _ = os.Getwd()
// 	}

// 	return bin
// }
