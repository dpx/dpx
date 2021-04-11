# `dpx` - Docker Process Executer

Run any *executable* inside container on your machine in few seconds (e.g. [go](https://hub.docker.com/_/golang), [node](https://hub.docker.com/_/node), [ruby](https://hub.docker.com/_/ruby), [python](https://hub.docker.com/_/python), [elixir](https://hub.docker.com/_/elixir), etc.).

![](./docs/dpx-go-demo.gif)

## Getting started

1. Intialize new container from any docker images (e.g. `golang`).
> This will create `dpx.yml` config file in current directory.

```sh
dpx init golang:1.16-alpine
```

2. Set `$PATH` variable.

```sh
eval $(dpx path) 
```

3. Link executable file and have fun!

```sh-session
# Create a link inside `.dpx/bin` directory
dpx link go

# Run any `go` command
go mod init app
```

## Installation

### macOS

`dpx` is available on macOS via [Homebrew](https://brew.sh).

```
brew install dpx/dpx/dpx
```

### Others

Packaged binaries can be downloaded from the [releases page](http://github.com/dpx/dpx/releases).

### Usage

#### `dpx exec [cmd]`

Execute a program inside container (equivalent to `docker exec -it [container] [cmd]`).

#### Default configurations

By default, all commands will inherit default configuration from `defaults:` section in `dpx.yml` file.

```yml
defaults:
  envs:
    - MODE=test
  user: docker
  workdir: /app
```

This will make every command that gets executed via `dpx exec [cmd]` inherits those settings.

#### Per-command configurations

Each command can also be configured with pre-define options.

**`envs`** set `ENV` before running a command.

```yml
commands:
  node:
    envs:
      - NODE_ENV=test
```

`dpx exec node` will be evaluated to `NODE_ENV=test node` in container.

**`options`** add options to your command.

```yml
commands:
  ps:
    options: aux
```

`dpx exec ps` will add `aux` options to `ps` command (e.g. `ps aux`).

**`workdir`** set current working dir for a command.

```yml
commands:
  ls:
    workdir: /app
```

`dpx exec ls` will be evaulated to `$ /app ls`

### Working with code editor, etc.

* **Visual Studio Code** - Once you've setup `$PATH` variable. Launch `vscode` from current directory with `code .` command.
> This will allow `vscode` to search for executable file under `$PATH` variable.

### Commands

```sh-session
COMMANDS:
   init, i   Setup docker image and create dpx.yml config file
   start, s  Start container from config
   stop      Stop a running container
   exec, x   Execute a command in container
   link, l   Link to an executable/binary inside container
   ps        Print current container ID
   path      Print $PATH variable
   help, h   Shows a list of commands or help for one command
```
