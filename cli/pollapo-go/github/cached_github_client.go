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
	verbose       bool
}

var mu = sync.Mutex{}
var rwmu = sync.RWMutex{}
var logName = "GitHub"

func NewCachedGitHubClient(cacheDir string, token string, verbose bool) GitHubClient {
	cacheFilepath := path.Join(cacheDir, "github-cache")

	util.PrintfVerbose(logName, verbose, "Loading github-cache from %s...\n", util.Yellow(cacheFilepath))
	c := util.LoadCache(cacheFilepath, &mu)
	util.PrintfVerbose(logName, verbose, "Loaded github-cache.\n")
	return CachedGitHubClient{
		NewGitHubClient(token),
		c,
		cacheFilepath,
		verbose,
	}
}

func (gc CachedGitHubClient) GetZipLink(owner string, repo string, ref string) (string, error) {
	key := fmt.Sprintf("zip-link:%v/%v@%v", owner, repo, ref)
	zipLink, found := gc.cache.Get(key)
	if found {
		return fmt.Sprintf("%v", zipLink), nil
	} else {
		// get commit of length 40 to avoid duplicated cache with vary length of refs
		commit, err := gc.getCommit(owner, repo, ref)
		if err != nil {
			return "", err
		}

		zipLink, err := gc.Default.GetZipLink(owner, repo, commit)
		if err != nil {
			return "", err
		}
		gc.cache.Set(key, zipLink, cache.DefaultExpiration)
		return zipLink, nil
	}
}

func (gc CachedGitHubClient) GetCommit(owner string, repo string, ref string) (string, error) {
	return gc.getCommit(owner, repo, ref)
}

func (gc CachedGitHubClient) getCommit(owner string, repo string, ref string) (string, error) {
	// Ref might have variable length, e.g. between 6~40
	// Anyway, The refs in different lengths will be stored separately.
	key := fmt.Sprintf("commit:%v/%v@%v", owner, repo, ref)
	commit, found := gc.cache.Get(key)
	if found {
		return fmt.Sprintf("%v", commit), nil
	} else {
		commit, err := gc.Default.GetCommit(owner, repo, ref)
		if err != nil {
			return "", err
		}
		gc.cache.Set(key, commit, cache.DefaultExpiration)
		key2 := fmt.Sprintf("commit:%v/%v@%v", owner, repo, commit) // the length of commit is 40
		gc.cache.Set(key2, commit, cache.DefaultExpiration)
		return commit, nil
	}
}

func (gc CachedGitHubClient) Flush() error {
	util.PrintfVerbose(logName, gc.verbose, "Save github-cache to %s...\n", util.Yellow(gc.cacheFilepath))
	err := util.SaveCache(gc.cache, gc.cacheFilepath, &rwmu)
	if err != nil {
		util.Printf("%s: %s\n", util.Red("Failed to save github-cache"), err)
	}
	return nil
}
