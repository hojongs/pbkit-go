package main

import (
	"os"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/cmds"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		EnableBashCompletion: true,
		Name:                 "pollapo-go",
		Usage:                "Protobuf dependency installer",
		Commands: []*cli.Command{
			&cmds.CommandInstall,
			&cmds.CommandLogin,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Sugar.Fatal(err)
	}
	log.Sugar.Sync() // flushes buffer, if any
}
