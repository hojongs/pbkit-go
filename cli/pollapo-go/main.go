package main

import (
	"log"
	"os"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/cache"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/cmds"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:  "pollapo-go",
		Usage: "Protobuf dependency installer",
		Commands: []*cli.Command{
			{
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
				},
				Action: func(c *cli.Context) error {
					cmds.NewCmdInstall(
						c.Bool("clean"),
						c.String("out-dir"),
						c.String("config"),
						cmds.GitHubZipDownloader{Token: c.String("token")},
						cmds.PollapoConfigFileLoader{},
						cache.FileSystemCache{},
					).Install()
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
