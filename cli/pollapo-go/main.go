package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/cache"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/cmds"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/github"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/mycolor"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/myzip"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/pollapo"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:    "pollapo-go",
		Usage:   "Protobuf dependency installer",
		Version: "0.2.0",
		Commands: []*cli.Command{
			{ // TODO: move to install file
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "Install dependencies.",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "clean",
						Aliases: []string{"c"},
						Usage:   "Clean cache directory before install",
						Value:   false,
					},
					&cli.StringFlag{
						Name:    "out-dir",
						Aliases: []string{"o"},
						Usage:   "Out directory",
						Value:   ".pollapo",
					},
					&cli.StringFlag{
						Name:    "token",
						Aliases: []string{"t"},
						Usage:   "GitHub OAuth token",
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"C"},
						Usage:   "Pollapo yml path",
						Value:   "pollapo.yml",
					},
					&cli.BoolFlag{
						Name:    "sync-unzip",
						Aliases: []string{"s"},
						Usage:   "Unzip synchronously. Use this flag if you get the error 'too many open files'",
						Value:   false,
					},
				},
				Action: func(c *cli.Context) error {
					var token string
					if len(c.String("token")) > 0 {
						token = c.String("token")
					} else {
						token = github.GetTokenFromGhHosts()
					}
					gc := github.NewClient(token)
					var uz myzip.Unzipper
					if c.Bool("sync-unzip") {
						fmt.Printf("%s\n", mycolor.Yellow("Sync unzip mode"))
						uz = myzip.SyncUnzipper{}
					} else {
						uz = myzip.ASyncUnzipper{}
					}
					cmds.NewCmdInstall(
						c.Bool("clean"),
						c.String("out-dir"),
						c.String("config"),
						myzip.NewGitHubZipDownloader(gc),
						uz,
						pollapo.FileConfigLoader{},
						cache.NewFileSystemCache(),
					).Install()
					return nil
				},
			},
			{ // TODO: move to login file
				Name:    "login",
				Aliases: []string{"l"},
				Usage:   "Sign in with GitHub account",
				Action: func(c *cli.Context) error {
					cmds.Login()
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
