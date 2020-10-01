# dockercat - Another Terminal UI for Docker

## About dockercat

Mix UI and functions of [lazydocker](https://github.com/jesseduffield/lazydocker) and [docui](https://github.com/skanehira/docui). Focus on managing existed images and containers.

Ijust make it so simple to find out information that I care about most: conteiner logs|stats|env, image size|unused, volume mountpoint. It's easy to release disk space for host.

UIKit: [tview](https://github.com/rivo/tview)

**It's NOT a good idea to use it in your production environment.**

## Supported OS

- Mac

## Required Tools

- Go Ver.1.13+
- Docker Engine Ver.18.06.1+
- Git

## Installation

### From Source

```
$ git clone https://github.com/nuknal/dockercat.git
$ cd dockercat/
$ GO111MODULE=on go install
```

### Homebrew

```
$ brew tap nuknal/taps
$ brew install dockercat
```

## TODO

- k8s
- docker-compose
- colorful
