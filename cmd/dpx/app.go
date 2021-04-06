package main

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/dpx/dpx/pkg/adapters/console"
	"github.com/dpx/dpx/pkg/adapters/container"
	"github.com/dpx/dpx/pkg/adapters/image"
	"github.com/dpx/dpx/pkg/core"
	"github.com/dpx/dpx/pkg/services/command"
	"github.com/dpx/dpx/pkg/services/config"
	"github.com/dpx/dpx/pkg/services/stream"
	"github.com/dpx/dpx/pkg/version"
	"github.com/urfave/cli/v2"
)

// App represents cli application.
type App struct {
	core.ContainerAdapter
	core.ImageAdapter
	core.ConfigService

	config  *core.Config
	stream  *stream.Stream
	context context.Context
}

// NewApp return cli's main application.
func NewApp() *cli.App {
	app := new()

	return &cli.App{
		Version: version.Version,
		Name:    config.Name,
		Usage:   "Docker process executor",
		Action:  app.defaultAction,
		Commands: []*cli.Command{
			newInitCmd(app),
			newStartCmd(app),
			newStopCmd(app),
			newExecCmd(app),
			newLinkCmd(app),
			newProcessCmd(app),
			newPathCmd(app),
		},
		UseShortOptionHandling: true,
	}
}

func (m *App) defaultAction(c *cli.Context) error {
	return cli.ShowAppHelp(c)
}

func new() *App {
	dk, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// prepare Main
	return &App{
		ConfigService: config.NewConfigService(),
		ContainerAdapter: container.NewWithOpts(
			container.WithDocker(dk),
		),
		ImageAdapter: image.New(
			image.WithDocker(dk),
		),
		config:  config.NewOrLoad(),
		stream:  stream.New(),
		context: context.Background(),
	}
}

func newInitCmd(app *App) *cli.Command {
	return console.NewInitCmd(command.NewInitCmd(
		app.ImageAdapter,
		app.ContainerAdapter,
		app.ConfigService,
		&core.CommandInitOptions{},
		app.config,
		app.stream.Out,
	), newStartCmd(app))
}

func newStartCmd(app *App) *cli.Command {
	return console.NewStartCmd(command.NewStartCmd(
		app.ContainerAdapter,
		app.config,
		&core.CommandStartOptions{},
		app.stream.Out,
	))
}

func newStopCmd(app *App) *cli.Command {
	return console.NewStopCmd(command.NewStopCmd(
		app.ContainerAdapter,
		app.config,
		&core.CommandStopOptions{},
		app.stream.Out,
	))
}

func newProcessCmd(app *App) *cli.Command {
	return console.NewProcessCmd(command.NewProcessCmd(
		app.ContainerAdapter,
		app.config,
		&core.CommandProcessOptions{},
		app.stream.Out,
	))
}

func newPathCmd(app *App) *cli.Command {
	return console.NewPathCmd(command.NewPathCmd(
		app.config,
		&core.CommandPathOptions{},
		app.stream.Out,
	))
}

func newLinkCmd(app *App) *cli.Command {
	return console.NewLinkCmd(command.NewLinkCmd(
		app.ContainerAdapter,
		app.ConfigService,
		app.config,
		&core.CommandLinkOptions{},
		app.stream.Out,
	))
}

func newExecCmd(app *App) *cli.Command {
	return console.NewExecCmd(command.NewExecCmd(
		app.ContainerAdapter,
		app.ConfigService,
		app.config,
		&core.CommandExecOptions{},
		app.stream.Out,
	))
}
