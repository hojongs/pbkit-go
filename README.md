# pbkit-go: Pbkit Rewritten in Go

Original proejct: https://github.com/pbkit/pbkit 

> Note: This project is very very early stage !!!

## Getting Started

### Prerequisite

- Install go: https://go.dev/doc/install

### Installation

```sh
# Ensure your working directory is the root of project
go install ./cli/pollapo-go
```

### Command Example

- install
  ```sh
  pollapo-go help
  pollapo-go i
  pollapo-go install
  ```

## Run Test

```sh
mockgen -source ./cli/pollapo-go/myzip/zip_downloader.go -destination ./cli/pollapo-go/myzip/zip_downloader_mock.go -package myzi
mockgen -source ./cli/pollapo-go/myzip/zip.go -destination ./cli/pollapo-go/myzip/zip_mock.go -package myzip
mockgen -source ./cli/pollapo-go/pollapo/pollapo_config_loader.go -destination ./cli/pollapo-go/pollapo/pollapo_config_loader_mock.go -package pollapo
go test -v ./cli/pollapo-go/cmds
```
