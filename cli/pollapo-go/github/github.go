package github

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v42/github"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/mycolor"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
}

func NewClient(token string) Client {
	var client *github.Client = nil
	if len(token) > 0 {
		client = initClientByToken(token)
	} else {
		client = github.NewClient(nil)
	}
	return Client{client}
}

func (gc Client) GetZipLink(owner string, repo string, ref string) string {
	opts := github.RepositoryContentGetOptions{Ref: ref}
	url, _, err := gc.client.Repositories.GetArchiveLink(context.Background(), owner, repo, github.Zipball, &opts, true)
	if err != nil {
		fmt.Printf("%s\n", mycolor.Red("error"))
		fmt.Printf("Login required. (%s/%s@%s)\n", owner, repo, ref)
		os.Exit(1)
	}
	return url.String()
}

func initClientByToken(token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}
