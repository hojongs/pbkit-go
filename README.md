# pbkit-go: Pbkit Rewritten in Go

Original proejct: https://github.com/pbkit/pbkit 

> Note: This project is very very early stage !!!

## Getting Started

### Prerequisite

- Install go: https://go.dev/doc/install

### Installation

```sh
go install -v github.com/hojongs/pbkit-go/cli/pollapo-go@latest
```

### Installation from source

```sh
git clone https://github.com/hojongs/pbkit-go.git
cd pbkit-go
# Ensure your working directory is the root of project
go install ./cli/pollapo-go
```

### Command Example

- help, login, install

There is an aliases for each command

```sh
pollapo-go install
pollapo-go i
```

## Run Test

```sh
mockgen -source ./cli/pollapo-go/myzip/zip_downloader.go -destination ./cli/pollapo-go/myzip/zip_downloader_mock.go -package myzi
mockgen -source ./cli/pollapo-go/myzip/zip.go -destination ./cli/pollapo-go/myzip/zip_mock.go -package myzip
mockgen -source ./cli/pollapo-go/pollapo/pollapo_config_loader.go -destination ./cli/pollapo-go/pollapo/pollapo_config_loader_mock.go -package pollapo
go test -v ./cli/pollapo-go/cmds
```

