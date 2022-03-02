package cmds

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/cache"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/github"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/mycolor"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/myzip"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/pollapo"
	"github.com/urfave/cli/v2"
)

var CommandInstall = cli.Command{
	Name:    "install",
	Aliases: []string{"i"},
	Usage:   "Install dependencies.",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "clean",
			Aliases: []string{"c"},
			Usage:   "Clean cache directory before install",
			Value:   false,
		},
		&cli.StringFlag{
			Name:    "out-dir",
			Aliases: []string{"o"},
			Usage:   "Out directory",
			Value:   ".pollapo",
		},
		&cli.StringFlag{
			Name:    "token",
			Aliases: []string{"t"},
			Usage:   "GitHub OAuth token",
		},
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"C"},
			Usage:   "Pollapo yml path",
			Value:   "pollapo.yml",
		},
	},
	Action: func(c *cli.Context) error {
		var token string
		if len(c.String("token")) > 0 {
			token = c.String("token")
		} else {
			token = github.GetTokenFromGhHosts()
		}
		gc := github.NewClient(token)
		newCmdInstall(
			c.Bool("clean"),
			c.String("out-dir"),
			c.String("config"),
			myzip.NewGitHubZipDownloader(gc),
			myzip.UnzipperImpl{},
			pollapo.FileConfigLoader{},
			cache.NewFileSystemCache(),
		).Install()
		return nil
	},
}

type cmdInstall struct {
	cleanCache     bool
	outDir         string
	pollapoYmlPath string
	zd             myzip.ZipDownloader
	uz             myzip.Unzipper
	loader         pollapo.ConfigLoader
	cache          cache.Cache
}

func newCmdInstall(
	cleanCache bool,
	outDir string,
	pollapoYmlPath string,
	zd myzip.ZipDownloader,
	uz myzip.Unzipper,
	loader pollapo.ConfigLoader,
	cache cache.Cache,
) cmdInstall {
	return cmdInstall{cleanCache, outDir, pollapoYmlPath, zd, uz, loader, cache}
}

func (cmd cmdInstall) Install() {
	if cmd.cleanCache {
		fmt.Printf("Clean cache root: %s\n", mycolor.Yellow(cmd.cache.GetRootLocation()))
		cmd.cache.Clean()
	}
	rootCfg, err := cmd.loader.GetPollapoConfig(cmd.pollapoYmlPath)
	if err != nil {
		fmt.Printf("%s\n", mycolor.Red("error"))
		absPath, err := filepath.Abs(cmd.pollapoYmlPath)
		if err != nil {
			log.Fatalw("Unknown error. Please retry.", err)
		}
		fmt.Printf("\"%s\" not found.\n", mycolor.Red(absPath))
		// TODO: Create absPath?
		os.Exit(1)
	}
	if err := os.RemoveAll(cmd.outDir); err != nil {
		log.Fatalw("Remove out dir", err, "outDir", cmd.outDir)
	}
	cmd.installDepsRecursive(rootCfg)
	fmt.Println("Done.")
}

func (cmd cmdInstall) installDepsRecursive(rootCfg pollapo.PollapoConfig) {
	cacheQueue := []string{}
	cacheQueue = append(cacheQueue, rootCfg.Deps...)
	depsMap := map[string]map[string][]string{} // depsMap[user/repo][ref]=froms
	origin := "<root>"
	for len(cacheQueue) > 0 {
		depTxt := cacheQueue[0]
		cacheQueue = cacheQueue[1:]

		dep, isOk := pollapo.ParseDep(depTxt)

		f := func(dep pollapo.PollapoDep) string { return dep.Owner + "/" + dep.Repo }
		if !isOk {
			log.Fatalw("Invalid dep", nil, "dep", depTxt)
		}
		if depsMap[f(dep)] == nil {
			depsMap[f(dep)] = map[string][]string{}
		}
		if depsMap[f(dep)][dep.Ref] != nil {
			depsMap[f(dep)][dep.Ref] = append(depsMap[f(dep)][dep.Ref], origin)
		} else {
			depsMap[f(dep)][dep.Ref] = []string{origin}
		}

		zipBin, err := cmd.cache.Get(cacheKeyOf(dep))
		var zipReader *zip.Reader = nil
		if err != nil || zipBin == nil {
			zipReader = cmd.downloadZip(dep)
		} else {
			fmt.Printf("Use cache of %s.\n", mycolor.Yellow(depTxt))
			zipReader = myzip.NewZipReader(zipBin)
		}
		cacheOutDir := filepath.Join(cmd.cache.GetRootLocation(), dep.Owner, dep.Repo)
		cmd.uz.UnzipFilter(zipReader, cacheOutDir, "pollapo.yml")

		depPollapoYmlPath := filepath.Join(cacheOutDir, "pollapo.yml")
		depCfg, err := cmd.loader.GetPollapoConfig(depPollapoYmlPath)
		if err == nil {
			cacheQueue = append(cacheQueue, depCfg.Deps...)
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
	sort.Strings(latestDeps)

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
		fmt.Printf("Installing %s...", mycolor.Yellow(dep.String()))
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

func (cmd cmdInstall) downloadZip(dep pollapo.PollapoDep) *zip.Reader {
	// log.Infow("Cache not found", "dep", mycolor.Yellow(cacheKeyOf(dep)))
	zipReader, zipBin := cmd.zd.GetZip(dep.Owner, dep.Repo, dep.Ref)
	fmt.Print("ok.")
	cmd.cache.Store(cacheKeyOf(dep), zipBin)
	fmt.Print(" Stored Cache.\n")
	return zipReader
}
