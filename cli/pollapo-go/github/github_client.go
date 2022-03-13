package github

import (
	"context"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

type GitHubClient interface {
	GetZipLink(owner string, repo string, ref string) (string, error)
	GetCommit(owner string, repo string, ref string) (string, error)
	Flush() error
}

type DefaultGitHubClient struct {
	client *github.Client
}

func NewGitHubClient(token string) GitHubClient {
	var client *github.Client = nil
	if len(token) > 0 {
		client = initClientByToken(token)
	} else {
		client = github.NewClient(nil)
	}
	return DefaultGitHubClient{client}
}

func (gc DefaultGitHubClient) GetZipLink(owner string, repo string, ref string) (string, error) {
	opts := github.RepositoryContentGetOptions{Ref: ref}
	url, _, err := gc.client.Repositories.GetArchiveLink(context.Background(), owner, repo, github.Zipball, &opts, true)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (gc DefaultGitHubClient) GetCommit(owner string, repo string, ref string) (string, error) {
	commit, _, err := gc.client.Repositories.GetCommitSHA1(context.Background(), owner, repo, ref, "")
	if err != nil {
		return "", err
	}
	return commit, nil
}

func (gc DefaultGitHubClient) Flush() error { return nil }

func initClientByToken(token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}
