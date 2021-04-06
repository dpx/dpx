package config

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/utils"
	"gopkg.in/yaml.v3"
)

type ConfigService struct {
}

var errNotFound = errors.New("cannot find dpx.yml config file")

// Command returns command's config (merged with default settings)
func (c *ConfigService) Command(cmd string, cf *core.Config) *core.CommandConfig {
	d := cf.Defaults

	if config, ok := cf.Commands[cmd]; ok {
		return merge(*d, config)
	}

	return d
}

func (c *ConfigService) Save(cf *core.Config, opts core.ConfigSaveOptions) error {
	// Don't save default settings
	if !opts.SaveDefaults {
		cf.Defaults = nil
	}

	var data bytes.Buffer
	enc := NewEncoder(&data)
	if err := enc.Encode(cf); err != nil {
		return err
	}

	dir := utils.GetCwd()
	file := path.Join(dir, FileName)

	// Set path
	if cf.Path == "" {
		cf.Path = file
	}

	return ioutil.WriteFile(file, data.Bytes(), 0644)
}

func (c *ConfigService) Write(cf *core.Config) error {
	return nil
}

func (c *ConfigService) CreateAliasFile() (core.Path, error) {
	return createAliasFile()
}

func (c *ConfigService) CreateBinFile(cmd string) (core.Path, error) {
	return createBinFile(cmd)
}

func New() *core.Config {
	return core.NewConfig()
}

func NewFromReader(r io.Reader) (*core.Config, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	config := New()
	err = yaml.Unmarshal(data, config)

	return config, err
}

func NewFromString(s string) (*core.Config, error) {
	r := strings.NewReader(s)

	return NewFromReader(r)
}

func NewFromFile(file string) (*core.Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	c, err := NewFromReader(f)
	c.Path = file

	return c, err
}

// Find finds dpx.yml file by searching from current dir and up to parent dirs.
func NewFromWorkDir() (*core.Config, error) {
	dir := utils.GetCwd()

	for i := 0; i < MaxLevel; i++ {
		p := path.Join(dir, FileName)
		if utils.FileExist(p) {
			return NewFromFile(p)
		}

		// reach root level
		if dir == "/" {
			break
		}

		dir = path.Dir(dir)
	}

	return New(), errNotFound
}

// NewOrLoad loads existing config, if not exists creates a new one
func NewOrLoad() *core.Config {
	c, err := NewFromWorkDir()
	if err != nil {
		c = New()
	}

	return c
}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}
