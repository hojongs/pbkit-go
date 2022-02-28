package cmds

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/cache"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/color"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/github"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/pollapo"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/zip"
)

type Cmd interface {
	Install()
}

type CmdContextInstall struct {
	clean          bool
	outDir         string
	pollapoYmlPath string
	zd             ZipDownloader
	loader         PollapoConfigLoader
}

func NewCmdContextInstall(
	clean bool,
	outDir string,
	pollapoYmlPath string,
	// TODO: add cache interface
	zd ZipDownloader,
	loader PollapoConfigLoader,
) CmdContextInstall {
	return CmdContextInstall{clean, outDir, pollapoYmlPath, zd, loader}
}

type PollapoConfigLoader interface {
	GetPollapoConfig(pollapoYmlPath string) (pollapo.PollapoConfig, error)
}

type PollapoConfigFileLoader struct{}

func (ctx CmdContextInstall) Install() {
	if ctx.clean {
		fmt.Printf("Clean cache root: %s\n", color.Yellow(cache.CacheRoot))
		cache.Clean()
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

func (ctx CmdContextInstall) installDepsRecursive(rootCfg pollapo.PollapoConfig) {
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

		zipBin, err := cache.Get(cacheKeyOf(dep))
		if err != nil {
			fmt.Printf("Cache not found of %s\n", color.Yellow(cacheKeyOf(dep)))
			// TODO: github authentication with pollapo login
			zipBin = ctx.zd.GetZipBin(dep.Owner, dep.Repo, dep.Ref)
			fmt.Print("ok\n")
			cache.Store(cacheKeyOf(dep), zipBin)
		} else {
			fmt.Printf("Use cache of %s.\n", color.Yellow(depTxt))
		}
		cacheOutDir := filepath.Join(cache.CacheRoot, dep.Owner, dep.Repo)
		// TODO: if already unzipped zipBin?
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

	log.Infow("", "latestDeps", latestDeps)
	for _, depTxt := range latestDeps {
		dep, isOk := pollapo.ParseDep(depTxt)
		if !isOk {
			log.Fatalw("Failed to parse dep", nil, "dep", depTxt)
		}
		depOutDir := filepath.Join(ctx.outDir, dep.Owner, dep.Repo)
		zipBin, err := cache.Get(cacheKeyOf(dep))
		if err != nil {
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
	// TODO: get latest
	return refs[0]
}
