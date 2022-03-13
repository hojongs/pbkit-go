package github

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type CachedGitHubClient struct {
	DefaultClient DefaultGitHubClient
	cache         *cache.Cache
}

func NewCachedGitHubClient(token string) GitHubClient {
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
