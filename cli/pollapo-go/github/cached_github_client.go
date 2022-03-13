package github

import (
	"fmt"

	"github.com/patrickmn/go-cache"
)

type CachedGitHubClient struct {
	Default GitHubClient
	cache   *cache.Cache
}

func NewCachedGitHubClient(token string) GitHubClient {
	// TODO: Load cache from file
	return CachedGitHubClient{
		NewGitHubClient(token),
		cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

func (gc CachedGitHubClient) GetZipLink(owner string, repo string, ref string) (string, error) {
	return gc.Default.GetZipLink(owner, repo, ref)
}

func (gc CachedGitHubClient) GetCommit(owner string, repo string, ref string) (string, error) {
	key := cacheKey(owner, repo, ref)
	commit, found := gc.cache.Get(key)
	if found {
		return fmt.Sprintf("%v", commit), nil
	} else {
		commit, err := gc.Default.GetCommit(owner, repo, ref)
		if err != nil {
			return "", err
		}
		gc.cache.Set(key, commit, cache.DefaultExpiration)
		return commit, nil
	}
}

func (gc CachedGitHubClient) Flush() error {
	// TODO: Save cache to file
	return nil
}
