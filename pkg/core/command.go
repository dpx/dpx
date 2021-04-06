package core

import "context"

type Command interface {
	Execute(ctx context.Context) error
}

type CommandInitOptions struct {
	Image string
	Name  string
	Skip  bool
}

type CommandLinkOptions struct {
	All  bool
	Cmds []string
}

type CommandExecOptions struct {
	All        bool
	Print      bool
	Cmd        string
	CmdArgs    []string
	ConfigPath string
	Save       bool
	Shell      string
	Envs       []string
	Compose    string
	User       string
	WorkDir    string
	Stdin      bool
	NoStdin    bool
	Tty        bool
	NoTty      bool
}

type CommandStartOptions struct {
	Name    string
	Envs    []string
	Ports   []string
	Volumes []string
}

type CommandStopOptions struct {
}

type CommandProcessOptions struct {
}

type CommandPathOptions struct {
	Delete bool
}
