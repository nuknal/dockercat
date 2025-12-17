# DockerCat - Terminal UI for Docker

<p align="center">
  <img src="images/screenshot.png" alt="DockerCat" />
</p>

## About

DockerCat is a terminal UI (TUI) for Docker. It focuses on day-to-day container and image operations, plus quick cleanup to free disk space.

UIKit: [tview](https://github.com/rivo/tview)

**Do not use it in production environments.**

## Features

- Panels: Containers, Images, Volumes, Networks, Cleanup, Info, Tasks
- Inspect: press `Enter` to view JSON details in the right pane
- Containers: start/stop/restart, pause/unpause, logs (`Ctrl+l`), stats (`Ctrl+s`), batch actions (select with `m`, then `U/S/D`)
- Images: remove, remove dangling (`Ctrl+d`), layer history (`s`), pull (`p`), tag (`t`), push (`P`)
- Volumes/Networks: inspect, remove, prune unused (`Ctrl+d`)
- Cleanup: system prune (`a`) or prune by type (`c/i/n/v`)
- Filter list: press `/` in the current panel
- Background tasks: long operations are queued and shown in the Tasks panel (cancel with `c`)

## Keybindings (Quick Start)

- Switch panels: `h` / `l`, `←` / `→`, `Tab`
- Filter: `/`
- Inspect: `Enter`
- Quit: `q`
- Scroll detail view: `Ctrl+j` / `Ctrl+k`

## Supported OS

- macOS

## Requirements

- Go 1.13+
- Docker Engine 18.06.1+

## Installation

### Homebrew

```
brew tap nuknal/taps
brew install dockercat
```

### From source

```
git clone https://github.com/nuknal/dockercat.git
cd dockercat
GO111MODULE=on go install
```

### Build binaries (macOS/Linux)

```
make build
ls -la bin/
```

## Usage

```
dockercat [flags]

  -H, --host          docker host, e.g. unix:///var/run/docker.sock
  -A, --api-version   docker api version, e.g. 1.39
  -R, --refresh       refresh interval, e.g. 5s
```

If `DOCKER_HOST` is set, dockercat uses Docker environment variables (including TLS settings via `DOCKER_TLS_VERIFY` / `DOCKER_CERT_PATH`).

## TODO

- k8s
- docker-compose
