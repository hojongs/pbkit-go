package cmds

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/cache"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/color"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/github"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/pollapo"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/zip"
)

type CmdInstall struct {
	cleanCache     bool
	outDir         string
	pollapoYmlPath string
	zd             ZipDownloader
	loader         PollapoConfigLoader
	cache          cache.Cache
}

func NewCmdInstall(
	cleanCache bool,
	outDir string,
	pollapoYmlPath string,
	// TODO: add cache interface
	zd ZipDownloader,
	loader PollapoConfigLoader,
	cache cache.Cache,
) CmdInstall {
	return CmdInstall{cleanCache, outDir, pollapoYmlPath, zd, loader, cache}
}

type PollapoConfigLoader interface {
	GetPollapoConfig(pollapoYmlPath string) (pollapo.PollapoConfig, error)
}

type PollapoConfigFileLoader struct{}

func (ctx CmdInstall) Install() {
	if ctx.cleanCache {
		fmt.Printf("Clean cache root: %s\n", color.Yellow(ctx.cache.GetRootLocation()))
		ctx.cache.Clean()
	}
	rootCfg, err := ctx.loader.GetPollapoConfig(ctx.pollapoYmlPath)
	if err != nil {
		log.Fatalw("Failed to read file", err, "filename", ctx.pollapoYmlPath)
	}
	ctx.installDepsRecursive(rootCfg)
}

func (_ PollapoConfigFileLoader) GetPollapoConfig(pollapoYmlPath string) (pollapo.PollapoConfig, error) {
	pollapoBytes, err := os.ReadFile(pollapoYmlPath)
	if err != nil {
		return pollapo.PollapoConfig{}, err
	} else {
		return pollapo.ParsePollapo(pollapoBytes), nil
	}
}

type ZipDownloader interface {
	GetZipBin(owner string, repo string, ref string) []byte
}

type GitHubZipDownloader struct {
	Token string
}

func (this GitHubZipDownloader) GetZipBin(owner string, repo string, ref string) []byte {
	// TODO: github authentication with token
	zipUrl := github.GetZipLink(owner, repo, ref)
	fmt.Printf("Downloading %s...", color.Yellow())
	resp, err := http.Get(zipUrl)
	if err != nil {
		log.Fatalw("Failed to HTTP Get", err, "dep", fmt.Sprintf("%s/%s@%v", owner, repo, ref))
	}
	if resp.StatusCode != 200 {
		log.Fatalw("HTTP Response is not OK", nil, "status", resp.StatusCode)
	}
	zipBin, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalw("Failed to Read HTTP Response body", err, "body", zipBin[:1024])
	}
	defer resp.Body.Close()

	return zipBin
}

func (ctx CmdInstall) installDepsRecursive(rootCfg pollapo.PollapoConfig) {
	cacheQueue := []string{}
	cacheQueue = append(cacheQueue, rootCfg.Deps...)
	depsMap := map[string]map[string][]string{} // depsMap[user/repo][ref]=froms
	origin := "<root>"
	for len(cacheQueue) > 0 {
		depTxt := cacheQueue[0]
		cacheQueue = cacheQueue[1:]

		dep, isOk := pollapo.ParseDep(depTxt)
		if !isOk {
			log.Fatalw("Invalid dep", nil, "dep", depTxt)
		}
		if depsMap[dep.Owner+"/"+dep.Repo] == nil {
			depsMap[dep.Owner+"/"+dep.Repo] = map[string][]string{}
		}
		if depsMap[dep.Owner+"/"+dep.Repo][dep.Ref] != nil {
			depsMap[dep.Owner+"/"+dep.Repo][dep.Ref] = append(depsMap[dep.Owner+"/"+dep.Repo][dep.Ref], origin)
		} else {
			depsMap[dep.Owner+"/"+dep.Repo][dep.Ref] = []string{origin}
		}

		zipBin, err := ctx.cache.Get(cacheKeyOf(dep))
		if err != nil || zipBin == nil {
			fmt.Printf("Cache not found of %s\n", color.Yellow(cacheKeyOf(dep)))
			// TODO: github authentication with pollapo login
			zipBin = ctx.zd.GetZipBin(dep.Owner, dep.Repo, dep.Ref)
			fmt.Print("ok\n")
			ctx.cache.Store(cacheKeyOf(dep), zipBin)
		} else {
			fmt.Printf("Use cache of %s.\n", color.Yellow(depTxt))
		}
		cacheOutDir := filepath.Join(ctx.cache.GetRootLocation(), dep.Owner, dep.Repo)
		zip.Unzip(zipBin, cacheOutDir)

		depPollapoYmlPath := filepath.Join(cacheOutDir, "pollapo.yml")
		depCfg, err := ctx.loader.GetPollapoConfig(depPollapoYmlPath)
		if err == nil {
			for _, nestedDep := range depCfg.Deps {
				cacheQueue = append(cacheQueue, nestedDep)
			}
		}
		origin = depTxt
	}

	latestDeps := []string{}
	for repoPath, depRefMap := range depsMap {
		refs := make([]string, 0, len(depRefMap))
		for k := range depRefMap {
			refs = append(refs, k)
		}
		depTxt := fmt.Sprintf("%s@%s", repoPath, latestRef(refs))
		latestDeps = append(latestDeps, depTxt)
	}

	for _, depTxt := range latestDeps {
		dep, isOk := pollapo.ParseDep(depTxt)
		if !isOk {
			log.Fatalw("Failed to parse dep", nil, "dep", depTxt)
		}
		depOutDir := filepath.Join(ctx.outDir, dep.Owner, dep.Repo)
		zipBin, err := ctx.cache.Get(cacheKeyOf(dep))
		if err != nil || zipBin == nil {
			log.Fatalw("Unexpected cache not found. cache has probably been removed during install", err, "dep", dep)
		}
		fmt.Printf("Installing %s...", color.Yellow(dep.String()))
		zip.Unzip(zipBin, depOutDir)
		fmt.Print("ok\n")
	}
}

func cacheKeyOf(dep pollapo.PollapoDep) string {
	return fmt.Sprintf("%v-%v-%v.zip", dep.Owner, dep.Repo, dep.Ref)
}

func latestRef(refs []string) string {
	// TODO: get latest if semver
	// https://github.com/pbkit/pbkit/blob/main/cli/pollapo/rev.ts#L7
	sortedRefs := refs
	sort.Strings(sortedRefs)
	return sortedRefs[0]
}
