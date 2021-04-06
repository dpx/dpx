package core

import "github.com/dpx/dpx/pkg/utils"

// Runtime settings
const (
	DefaultShell     = "sh"
	DefaultTTY       = false
	DefaultStdin     = true
	DefaultAutoTTY   = true
	DefaultAutoStdin = true
)

// Config represents global configuration
//
// Loading orders during init & start:
// 	 - core.NewConfig()  // without default settings (CommandConfig)
//	 - command line args
// Loading orders during exec:
//	 - core.NewConfigWithDefaults()
//	 - dpx.yml
//	 - runtime defaults
//	 - command line args
// cli args > bin (dpx.yml) > global (dpx.yml) > default
type Config struct {
	Image   string   `yaml:"image"`
	Name    string   `yaml:"name"`
	Compose string   `yaml:"compose,omitempty"`
	Path    string   `yaml:"-"`
	Runs    []string `yaml:"runs,omitempty"`
	Envs    []string `yaml:"envs,omitempty"`
	Ports   []string `yaml:"ports,omitempty"`
	Volumes []string `yaml:"volumes,omitempty"`

	// Command's config
	Defaults *CommandConfig           `yaml:"defaults,omitempty"`
	Commands map[string]CommandConfig `yaml:"commands,omitempty"`
}

// CommandConfig represents command's configuration
type CommandConfig struct {
	Compose    string   `yaml:"compose,omitempty"`
	TTY        *bool    `yaml:"tty,omitempty"`
	Stdin      *bool    `yaml:"stdin,omitempty"`
	Envs       []string `yaml:"envs,omitempty"`
	User       string   `yaml:"user,omitempty"`
	WorkingDir *string  `yaml:"workdir,omitempty"`
	Shell      string   `yaml:"shell,omitempty"`
	Options    string   `yaml:"options,omitempty"`
}

type Path string
type ConfigService interface {
	Command(cmd string, cf *Config) *CommandConfig
	Save(cf *Config, opts ConfigSaveOptions) error
	CreateAliasFile() (Path, error)
	CreateBinFile(cmd string) (Path, error)
}

type ConfigSaveOptions struct {
	SaveDefaults bool
}

func NewConfig() *Config {
	return &Config{
		Defaults: NewCommandConfig(),
	}
}

func NewConfigWithDefaults() *Config {
	return &Config{
		Defaults: NewCommandConfig(),
	}
}

func NewCommandConfig() *CommandConfig {
	return &CommandConfig{
		Shell:      DefaultShell,
		TTY:        utils.Bool(DefaultTTY),
		Stdin:      utils.Bool(DefaultStdin),
		WorkingDir: utils.String(""),
	}
}
