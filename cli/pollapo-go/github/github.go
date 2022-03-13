package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v42/github"
	"github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
)

type GitHubClient interface {
	GetZipLink(owner string, repo string, ref string) (string, error)
	GetCommit(owner string, repo string, ref string) (string, error)
}

type DefaultGitHubClient struct {
	client *github.Client
}

type CachedGitHubClient struct {
	DefaultClient DefaultGitHubClient
	cache         *cache.Cache
}

func NewGitHubClient(token string) DefaultGitHubClient {
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

func NewCachedClient(token string) GitHubClient {
	c := cache.New(5*time.Minute, 5*time.Minute)
	return CachedGitHubClient{NewGitHubClient(token), c}
}

func (gc CachedGitHubClient) GetZipLink(owner string, repo string, ref string) (string, error) {
	return gc.DefaultClient.GetZipLink(owner, repo, ref)
}

func (gc CachedGitHubClient) GetCommit(owner string, repo string, ref string) (string, error) {
	key := cacheKey(owner, repo, ref)
	commit, found := gc.cache.Get(key)
	if !found {
		var err error
		commit, err = gc.DefaultClient.GetCommit(owner, repo, ref)
		if err != nil {
			return "", err
		}
	}
	gc.cache.Set(key, commit, cache.DefaultExpiration)
	return fmt.Sprintf("%v", commit), nil
}

func initClientByToken(token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

func cacheKey(owner string, repo string, ref string) string {
	return fmt.Sprintf("%v/%v@%v", owner, repo, ref)
}
