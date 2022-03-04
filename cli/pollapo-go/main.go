package main

import (
	"log"
	"os"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/cmds"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		EnableBashCompletion: true,
		Name:                 "pollapo-go",
		Usage:                "Protobuf dependency installer",
		Version:              "0.3.2",
		Commands: []*cli.Command{
			&cmds.CommandInstall,
			&cmds.CommandLogin,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
