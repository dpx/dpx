package mock

import (
	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/services/config"
)

type ConfigService struct {
	CommandFn         func(cmd string, cf *core.Config) *core.CommandConfig
	SaveFn            func(cf *core.Config, opts core.ConfigSaveOptions) error
	CreateAliasFileFn func() (core.Path, error)
	CreateBinFileFn   func(cmd string) (core.Path, error)
}

func (c *ConfigService) Command(cmd string, cf *core.Config) *core.CommandConfig {
	return c.CommandFn(cmd, cf)
}

func (c *ConfigService) Save(cf *core.Config, opts core.ConfigSaveOptions) error {
	return c.SaveFn(cf, opts)
}

func (c *ConfigService) CreateAliasFile() (core.Path, error) {
	return c.CreateAliasFileFn()
}

func (c *ConfigService) CreateBinFile(cmd string) (core.Path, error) {
	return c.CreateBinFileFn(cmd)
}

func NewConfigService() *ConfigService {
	cfs := config.NewConfigService()

	return &ConfigService{
		CommandFn: func(cmd string, cf *core.Config) *core.CommandConfig {
			return cfs.Command(cmd, cf)
		},
		SaveFn: func(cf *core.Config, opts core.ConfigSaveOptions) error {
			return nil
		},
		CreateAliasFileFn: func() (core.Path, error) {
			return "dpx.yml", nil
		},
		CreateBinFileFn: func(cmd string) (core.Path, error) {
			return "test", nil
		},
	}
}
