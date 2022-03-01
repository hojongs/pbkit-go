package cmds

import (
	"fmt"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/github"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
)

func Login() {
	token := github.GetTokenFromGhHosts()
	if len(token) == 0 {
		log.Infow("Token not found.")
		token = github.TryOauthFlow()
		github.WriteTokenGhHosts(token)
	} else {
		fmt.Println("You're already logged into github.com.")
	}
}
