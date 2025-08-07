# pman
pman is a command-line tool to keep track of all your side projects.

## Why?
I needed something to keep track of all my side projects.

## Install using the go package manager

```shell
go install github.com/theredditbandit/pman@latest
```

## Usage

```
A cli project manager

Usage:
  pman [flags]
  pman [command]

Available Commands:
  add         Adds a project directory to the index
  alias       Sets the alias for a project, whose name might be too big
  completion  Generate the autocompletion script for the specified shell
  delete      Deletes a project from the index database. This does not delete the project from the filesystem
  help        Help about any command
  i           Launches pman in interactive mode
  info        The info command pretty prints the README.md file present at the root of the specified project.
  init        Takes exactly 1 argument, a directory name, and initializes it as a project directory.
  ls          List all indexed projects along with their status
  reset       Deletes the current indexed projects, run pman init to reindex the projects
  set         Set the status of a project
  status      Get the status of a project

Flags:
  -h, --help      help for pman
  -v, --version   version for pman

Use "pman [command] --help" for more information about a command.
```

## How does it work

### init
When you run `pman init .` in any directory, it will look for subdirectories that contain a README.md or a .git folder and consider it as a project directory.
![image](https://github.com/theredditbandit/pman/assets/85390033/b9d6fcd3-41ca-4bd2-aa32-9b3e9bff1be8)

### set, status, info and filter
Set the status of any project using `pman set <projectName>`
![image](https://github.com/theredditbandit/pman/assets/85390033/1c9658ab-4280-435e-8d30-52963f656cc6)

Get the status of any project individually using `pman status <projectName>`
![image](https://github.com/theredditbandit/pman/assets/85390033/5466c077-4886-40db-b486-261738f06b4a)

Filter the results by status while listing projects using `pman ls --f <status>`
![image](https://github.com/theredditbandit/pman/assets/85390033/f8311d11-7fda-48f2-a634-daaf4ded90f2)

Print the README contents of a project using `pman info <projectName>`
![image](https://github.com/theredditbandit/pman/assets/85390033/6eabda18-479e-445b-8a6a-7b5b370e3e49)

### Interactive mode
Launch pman in interactive mode using `pman i`
![image](https://github.com/theredditbandit/pman/assets/85390033/c9f9f836-d1b3-45c0-a36f-eda7ed842e1a)


## Running Tests

To run the tests for this project, use the following command:

```shell
go test -v ./...
```

### Coverage

To generate a coverage report, use the following command:

```shell
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=theredditbandit/pman&type=Date)](https://star-history.com/#theredditbandit/pman&Date)
