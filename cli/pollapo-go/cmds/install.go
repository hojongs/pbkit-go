package cmds

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/Masterminds/semver"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/cache"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/github"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/myzip"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/pollapo"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/util"
	"github.com/urfave/cli/v2"
)

var CommandInstall = cli.Command{
	Name:    "install",
	Aliases: []string{"i"},
	Usage:   "Install dependencies.",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Print verbose logs",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "clean",
			Aliases: []string{"c", "clean-cache"},
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
			Name:  "config",
			Usage: "Pollapo yml path",
			Value: "pollapo.yml",
		},
	},
	Action: func(c *cli.Context) error {
		if c.Args().Len() >= 2 {
			util.Printf("Arguments are not required.\n")
			util.Printf("Given arguments: count %v, values %v\n", util.Yellow(c.Args().Len()), util.Yellow(c.Args()))
			os.Exit(1)
		}
		if c.Bool("verbose") {
			util.Printf("Flag verbose: %v\n", util.Yellow(c.Bool("verbose")))
			util.Printf("Flag clean: %v\n", util.Yellow(c.Bool("clean")))
			util.Printf("Flag out-dir: %v\n", util.Yellow(c.String("out-dir")))
			if c.String("token") != "" {
				util.Printf("Flag token: %v\n", util.Yellow(c.String("token")))
			}
			util.Printf("Flag config: %v\n", util.Yellow(c.String("config")))
		}
		token := c.String("token")
		if token == "" {
			token = github.GetTokenFromGhHosts()
		}
		gc := github.NewCachedGitHubClient(token)
		newCmdInstall(
			c.Bool("clean"),
			c.String("out-dir"),
			c.String("config"),
			gc,
			myzip.NewZipDownloader(),
			myzip.UnzipperImpl{},
			pollapo.FileConfigLoader{},
			cache.NewFileSystemCache(),
			c.Bool("verbose"),
		).Install()
		return nil
	},
}

type cmdInstall struct {
	cleanCache     bool
	outDir         string
	pollapoYmlPath string
	gc             github.GitHubClient
	zd             myzip.ZipDownloader
	uz             myzip.Unzipper
	loader         pollapo.ConfigLoader
	cache          cache.Cache
	verbose        bool
}

func newCmdInstall(
	cleanCache bool,
	outDir string,
	pollapoYmlPath string,
	gc github.GitHubClient,
	zd myzip.ZipDownloader,
	uz myzip.Unzipper,
	loader pollapo.ConfigLoader,
	cache cache.Cache,
	verbose bool,
) cmdInstall {
	return cmdInstall{cleanCache, outDir, pollapoYmlPath, gc, zd, uz, loader, cache, verbose}
}

func (cmd cmdInstall) Install() {
	if cmd.cleanCache {
		util.Printf("Clean cache root: %s\n", util.Yellow(cmd.cache.GetRootLocation()))
		cmd.cache.Clean()
	}
	rootCfg, err := cmd.loader.GetPollapoConfig(cmd.pollapoYmlPath)
	if err != nil {
		util.Printf("%s\n", util.Red("error"))
		absPath, err := filepath.Abs(cmd.pollapoYmlPath)
		if err != nil {
			log.Fatalw("Unknown error. Please retry.", err)
		}
		util.Printf("%s not found.\n", util.Red(absPath))
		// TODO: Ask create pollapo.yml
		os.Exit(1)
	}
	cmd.printfIfVerbose("Clean out directory %s.\n", util.Yellow(cmd.outDir))
	if err := os.RemoveAll(cmd.outDir); err != nil {
		log.Fatalw("Remove out dir", err, "outDir", cmd.outDir)
	}
	cmd.installDepsRecursive(rootCfg)
	// TODO: call gc.Flush()
	util.Println("Done.")
}

func (cmd cmdInstall) installDepsRecursive(rootCfg pollapo.PollapoConfig) {
	depHandleQueue := []pollapo.PollapoDep{}
	for _, dep := range rootCfg.Deps {
		cmd.printfIfVerbose("Enqueue %s.\n", util.Yellow(dep))
	}
	for _, depTxt := range rootCfg.Deps {
		dep, isOk := pollapo.ParseDep(depTxt)
		if !isOk {
			log.Fatalw("Invalid dep", nil, "dep", depTxt)
		}
		depHandleQueue = append(depHandleQueue, dep)
	}
	depsMap := map[string]map[string][]string{} // depsMap[user/repo][ref]=froms
	origin := "<root>"
	for len(depHandleQueue) > 0 {
		// cache zips concurrently
		wg := sync.WaitGroup{}
		wg.Add(len(depHandleQueue))
		for _, dep := range depHandleQueue {
			go cmd.cacheZipIfMiss(dep, &wg)
		}
		wg.Wait()

		queue := []pollapo.PollapoDep{}
		for _, dep := range depHandleQueue {
			// TODO: froms are unused. command 'why' will use it maybe.
			putDepIntoMap(depsMap, dep, origin)

			// get dependency zip (the dep cached)
			var zipReader *zip.Reader = nil
			zipBin, _ := cmd.cache.Get(cacheKeyOf(dep, "zip"))
			zipReader = myzip.NewZipReader(zipBin)

			// read pollapo.yml & enqueue deps
			pollapoFile := myzip.GetFileByName(zipReader, "pollapo.yml")
			if pollapoFile != nil {
				// get pollapo config
				rc, err := pollapoFile.Open()
				if err != nil {
					log.Fatalw("Failed to open pollapo file", err)
				}
				bin, err := io.ReadAll(rc)
				rc.Close()
				if err != nil {
					log.Fatalw("Failed to read pollapo file", err)
				}
				depCfg := pollapo.ParsePollapo(bin)
				for _, depTxt := range depCfg.Deps {
					dep, isOk := pollapo.ParseDep(depTxt)
					if !isOk {
						log.Fatalw("Invalid dep", nil, "dep", depTxt)
					}
					queue = append(queue, dep)
					cmd.printfIfVerbose("Enqueue %s.\n", util.Yellow(dep))
				}
			}

			origin = dep.String()
		}
		depHandleQueue = queue
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
		// TODO: 2 layer cache: in-memory, fs
		zipBin, err := cmd.cache.Get(cacheKeyOf(dep, "zip"))
		var zipReader *zip.Reader = nil
		if err != nil || zipBin == nil {
			zipReader = cmd.getAndCacheZip(dep)
		} else {
			zipReader = myzip.NewZipReader(zipBin)
		}
		util.Printf("Installing %s...", util.Yellow(dep.String()))
		cmd.uz.Unzip(zipReader, depOutDir)
		util.Print("ok\n")
	}
}

func cacheKeyOf(dep pollapo.PollapoDep, fileExt string) string {
	return fmt.Sprintf("%v-%v-%v.%v", dep.Owner, dep.Repo, dep.Ref, fileExt)
}

type RefArray []string

func (refa RefArray) Len() int {
	return len(refa)
}

func (refs RefArray) Less(i, j int) bool {
	aa, erra := semver.NewVersion(refs[i])
	bb, errb := semver.NewVersion(refs[j])
	if erra == nil && errb != nil { // aa is only semver
		return false
	}
	if erra != nil && errb == nil { // bb is only semver
		return true
	}
	if erra == nil && errb == nil { // both aa and bb are semvers
		return aa.Compare(bb) < 0
	}
	return refs[i] < refs[j]
}

func (refs RefArray) Swap(i, j int) {
	refs[i], refs[j] = refs[j], refs[i]
}

func latestRef(refs RefArray) string {
	sortedRefs := refs
	sort.Sort(sortedRefs)
	return refs[len(refs)-1]
}

func (cmd cmdInstall) getAndCacheZip(dep pollapo.PollapoDep) *zip.Reader {
	// log.Infow("Cache not found", "dep", util.Yellow(cacheKeyOf(dep)))
	zipUrl, err := cmd.gc.GetZipLink(dep.Owner, dep.Repo, dep.Ref)
	if err != nil {
		util.Printf("%s\n", util.Red("error"))
		util.Printf("Login required. (%s)\n", dep)
		os.Exit(1)
	}
	zipReader, zipBin := cmd.zd.GetZip(zipUrl)
	cmd.cache.Store(cacheKeyOf(dep, "zip"), zipBin)
	return zipReader
}

func (cmd cmdInstall) printfIfVerbose(format string, a ...interface{}) (n int, err error) {
	if cmd.verbose {
		return util.Printf(format, a...)
	} else {
		return 0, nil
	}
}

func putDepIntoMap(depsMap map[string]map[string][]string, dep pollapo.PollapoDep, origin string) {
	f := func(dep pollapo.PollapoDep) string { return dep.Owner + "/" + dep.Repo }
	if depsMap[f(dep)] == nil {
		depsMap[f(dep)] = map[string][]string{}
	}
	if depsMap[f(dep)][dep.Ref] != nil {
		depsMap[f(dep)][dep.Ref] = append(depsMap[f(dep)][dep.Ref], origin)
	} else {
		depsMap[f(dep)][dep.Ref] = []string{origin}
	}
}

func (cmd cmdInstall) cacheZipIfMiss(dep pollapo.PollapoDep, wg *sync.WaitGroup) {
	if _, err := cmd.cache.Get(cacheKeyOf(dep, "zip")); err != nil {
		util.Printf("Downloading %s...\n", util.Yellow(dep.String()))
		cmd.getAndCacheZip(dep)
		util.Printf("Stored cache %s\n", util.Yellow(dep.String()))
	} else {
		cmd.printfIfVerbose("Found cache of %s.\n", util.Yellow(dep.String()))
	}
	wg.Done()
}
