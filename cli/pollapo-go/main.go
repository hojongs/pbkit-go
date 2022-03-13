package main

import (
	"os"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/cmds"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/util"
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
		util.Sugar.Fatal(err)
	}
	util.Sugar.Sync() // flushes buffer, if any
}
