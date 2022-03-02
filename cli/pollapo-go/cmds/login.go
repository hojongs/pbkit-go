package cmds

import (
	"fmt"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/github"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
	"github.com/urfave/cli/v2"
)

var CommandLogin = cli.Command{
	Name:    "login",
	Aliases: []string{"l"},
	Usage:   "Sign in with GitHub account",
	Action: func(c *cli.Context) error {
		login()
		return nil
	},
}

func login() {
	token := github.GetTokenFromGhHosts()
	if len(token) == 0 {
		log.Infow("Token not found.")
		token = github.TryOauthFlow()
		github.WriteTokenGhHosts(token)
	} else {
		fmt.Println("You're already logged into github.com.")
	}
}
