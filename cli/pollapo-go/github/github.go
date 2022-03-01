package github

import (
	"context"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

type GitHubClient struct {
	client *github.Client
}

func NewGitHubClient(token string) GitHubClient {
	var client *github.Client = nil
	if len(token) > 0 {
		client = initClientByToken(token)
	} else {
		client = github.NewClient(nil)
	}
	return GitHubClient{client}
}

func (gc GitHubClient) GetZipLink(owner string, repo string, ref string) string {
	opts := github.RepositoryContentGetOptions{Ref: ref}
	url, _, _ := gc.client.Repositories.GetArchiveLink(context.Background(), owner, repo, github.Zipball, &opts, true)
	return url.String()
}

func initClientByToken(token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}
