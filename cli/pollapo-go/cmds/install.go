package cmds

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/cache"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/color"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/myzip"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/pollapo"
)

type CmdInstall struct {
	cleanCache     bool
	outDir         string
	pollapoYmlPath string
	zd             myzip.ZipDownloader
	uz             myzip.Unzipper
	loader         pollapo.ConfigLoader
	cache          cache.Cache
}

func NewCmdInstall(
	cleanCache bool,
	outDir string,
	pollapoYmlPath string,
	zd myzip.ZipDownloader,
	uz myzip.Unzipper,
	loader pollapo.ConfigLoader,
	cache cache.Cache,
) CmdInstall {
	return CmdInstall{cleanCache, outDir, pollapoYmlPath, zd, uz, loader, cache}
}

func (cmd CmdInstall) Install() {
	if cmd.cleanCache {
		fmt.Printf("Clean cache root: %s\n", color.Yellow(cmd.cache.GetRootLocation()))
		cmd.cache.Clean()
	}
	rootCfg, err := cmd.loader.GetPollapoConfig(cmd.pollapoYmlPath)
	if err != nil {
		log.Fatalw("Failed to read file", err, "filename", cmd.pollapoYmlPath)
	}
	cmd.installDepsRecursive(rootCfg)
}

func (cmd CmdInstall) installDepsRecursive(rootCfg pollapo.PollapoConfig) {
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

		zipBin, err := cmd.cache.Get(cacheKeyOf(dep))
		var zipReader *zip.Reader = nil
		if err != nil || zipBin == nil {
			zipReader = cmd.downloadZip(dep)
		} else {
			fmt.Printf("Use cache of %s.\n", color.Yellow(depTxt))
			zipReader = myzip.NewZipReader(zipBin)
		}
		cacheOutDir := filepath.Join(cmd.cache.GetRootLocation(), dep.Owner, dep.Repo)
		cmd.uz.Unzip(zipReader, cacheOutDir)

		depPollapoYmlPath := filepath.Join(cacheOutDir, "pollapo.yml")
		depCfg, err := cmd.loader.GetPollapoConfig(depPollapoYmlPath)
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
		depOutDir := filepath.Join(cmd.outDir, dep.Owner, dep.Repo)
		zipBin, err := cmd.cache.Get(cacheKeyOf(dep))
		var zipReader *zip.Reader = nil
		if err != nil || zipBin == nil {
			zipReader = cmd.downloadZip(dep)
		} else {
			zipReader = myzip.NewZipReader(zipBin)
		}
		fmt.Printf("Installing %s...", color.Yellow(dep.String()))
		cmd.uz.Unzip(zipReader, depOutDir)
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

func (cmd CmdInstall) downloadZip(dep pollapo.PollapoDep) *zip.Reader {
	fmt.Printf("Cache not found of %s\n", color.Yellow(cacheKeyOf(dep)))
	// TODO: github authentication with pollapo login
	zipReader, zipBin := cmd.zd.GetZip(dep.Owner, dep.Repo, dep.Ref)
	fmt.Print("ok")
	cmd.cache.Store(cacheKeyOf(dep), zipBin)
	fmt.Print("Stored Cache.\n")
	return zipReader
}
