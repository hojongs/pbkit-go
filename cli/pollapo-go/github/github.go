package github

import (
	"context"

	"github.com/google/go-github/v42/github"
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

func (gc Client) GetZipLink(owner string, repo string, ref string) (string, error) {
	opts := github.RepositoryContentGetOptions{Ref: ref}
	url, _, err := gc.client.Repositories.GetArchiveLink(context.Background(), owner, repo, github.Zipball, &opts, true)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func initClientByToken(token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}
