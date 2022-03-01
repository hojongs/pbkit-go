# pbkit-go: Pbkit Rewritten in Go

Original proejct: https://github.com/pbkit/pbkit 

> Note: This project is very very early stage !!!

# Pollapo: Protobuf Dependency Installer

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

### Commands

#### Help

![image](https://user-images.githubusercontent.com/15096588/156098794-4babe731-5c16-4742-83cc-db707b66afae.png)

#### Install

![image](https://user-images.githubusercontent.com/15096588/156098974-922c4269-2b4a-4d27-a0f0-b0818aa94bd1.png)

## Run Test

```sh
mockgen -source ./cli/pollapo-go/myzip/zip_downloader.go -destination ./cli/pollapo-go/myzip/zip_downloader_mock.go -package myzi
mockgen -source ./cli/pollapo-go/myzip/zip.go -destination ./cli/pollapo-go/myzip/zip_mock.go -package myzip
mockgen -source ./cli/pollapo-go/pollapo/pollapo_config_loader.go -destination ./cli/pollapo-go/pollapo/pollapo_config_loader_mock.go -package pollapo
go test -v ./cli/pollapo-go/cmds
```

