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

## Auto completion

### Zsh

Run this command once to download zsh_autocomplete, and add a line to your .zshrc

```sh
mkdir -p $HOME/.config/pollapo-go
curl https://raw.githubusercontent.com/urfave/cli/master/autocomplete/zsh_autocomplete > $HOME/.config/pollapo-go/zsh_autocomplete


echo 'PROG=pollapo-go' >> $HOME/.zshrc
echo '_CLI_ZSH_AUTOCOMPLETE_HACK=1' >> $HOME/.zshrc
echo '. $HOME/.config/pollapo-go/zsh_autocomplete' >> $HOME/.zshrc
```

For more detail
- Bash: https://github.com/urfave/cli/blob/master/docs/v2/manual.md#distribution-and-persistent-autocompletion
- Zsh: https://github.com/urfave/cli/blob/master/docs/v2/manual.md#zsh-support

## Run Test

```sh
mockgen -source ./cli/pollapo-go/myzip/zip_downloader.go -destination ./cli/pollapo-go/myzip/zip_downloader_mock.go -package myzi
mockgen -source ./cli/pollapo-go/myzip/zip.go -destination ./cli/pollapo-go/myzip/zip_mock.go -package myzip
mockgen -source ./cli/pollapo-go/pollapo/pollapo_config_loader.go -destination ./cli/pollapo-go/pollapo/pollapo_config_loader_mock.go -package pollapo
go test -v ./cli/pollapo-go/cmds
```

