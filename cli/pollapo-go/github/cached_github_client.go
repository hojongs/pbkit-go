package github

import (
	"fmt"
	"path"
	"sync"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/util"
	"github.com/patrickmn/go-cache"
)

type CachedGitHubClient struct {
	Default       GitHubClient
	cache         *cache.Cache
	cacheFilepath string
}

var mu = sync.Mutex{}
var rwmu = sync.RWMutex{}

func NewCachedGitHubClient(cacheDir string, token string) GitHubClient {
	cacheFilepath := path.Join(cacheDir, "github-cache")
	return CachedGitHubClient{
		NewGitHubClient(token),
		util.LoadCache(cacheFilepath, &mu),
		cacheFilepath,
	}
}

func (gc CachedGitHubClient) GetZipLink(owner string, repo string, ref string) (string, error) {
	return gc.Default.GetZipLink(owner, repo, ref)
}

func (gc CachedGitHubClient) GetCommit(owner string, repo string, ref string) (string, error) {
	key := fmt.Sprintf("%v/%v@%v", owner, repo, ref)
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
	util.SaveCache(gc.cache, gc.cacheFilepath, &rwmu)
	return nil
}
