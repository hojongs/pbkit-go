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
	"github.com/hojongs/pbkit-go/cli/pollapo-go/github"
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
			util.Printf("VERBOSE[Flag]: verbose: %v\n", util.Yellow(c.Bool("verbose")))
			util.Printf("VERBOSE[Flag]: clean: %v\n", util.Yellow(c.Bool("clean")))
			util.Printf("VERBOSE[Flag]: out-dir: %v\n", util.Yellow(c.String("out-dir")))
			util.Printf("VERBOSE[Flag]: token: %v\n", util.Yellow(c.String("token")))
			util.Printf("VERBOSE[Flag]: config: %v\n", util.Yellow(c.String("config")))
		}
		token := c.String("token")
		if token == "" {
			token = github.GetTokenFromGhHosts()
		}

		if c.Bool("clean") {
			util.Printf("Clean cache root: %s\n", util.Yellow(util.GetDefaultCacheRoot()))
			os.RemoveAll(util.GetDefaultCacheRoot())
		}
		onCacheMiss := func(cacheKey string) { util.Printf("Downloading %s...\n", util.Yellow(cacheKey)) }
		onCacheStore := func(cacheKey string) { util.Printf("Store cache %s...\n", util.Yellow(cacheKey)) }
		onCacheHit := func(cacheKey string) { util.Printf("Found cache of %s.\n", util.Yellow(cacheKey)) }
		newCmdInstall(
			c.String("out-dir"),
			c.String("config"),
			github.NewCachedGitHubClient(util.GetDefaultCacheRoot(), token, c.Bool("verbose")),
			myzip.NewCachedZipDownloader(util.GetDefaultCacheRoot(), c.Bool("verbose"), onCacheMiss, onCacheStore, onCacheHit),
			myzip.UnzipperImpl{},
			pollapo.FileConfigLoader{},
			c.Bool("verbose"),
		).Install()
		return nil
	},
}

type cmdInstall struct {
	outDir         string
	pollapoYmlPath string
	gc             github.GitHubClient
	zd             myzip.ZipDownloader
	uz             myzip.Unzipper
	loader         pollapo.ConfigLoader
	verbose        bool
}

func newCmdInstall(
	outDir string,
	pollapoYmlPath string,
	gc github.GitHubClient,
	zd myzip.ZipDownloader,
	uz myzip.Unzipper,
	loader pollapo.ConfigLoader,
	verbose bool,
) cmdInstall {
	return cmdInstall{outDir, pollapoYmlPath, gc, zd, uz, loader, verbose}
}

var logName = "Install"

func (cmd cmdInstall) Install() {
	rootCfg, err := cmd.loader.GetPollapoConfig(cmd.pollapoYmlPath)
	if err != nil {
		util.Printf("%s\n", util.Red("error"))
		absPath, err := filepath.Abs(cmd.pollapoYmlPath)
		if err != nil {
			util.Sugar.Fatalw("Unknown error. Please retry.", err)
		}
		util.Printf("%s not found.\n", util.Red(absPath))
		// TODO: Ask create pollapo.yml
		os.Exit(1)
	}
	util.PrintfVerbose(logName, cmd.verbose, "Clean out directory %s.\n", util.Yellow(cmd.outDir))
	if err := os.RemoveAll(cmd.outDir); err != nil {
		util.Sugar.Fatalw("Remove out dir", err, "outDir", cmd.outDir)
	}
	cmd.installDepsRecursive(&rootCfg)
	cmd.gc.Flush()
	cmd.zd.Flush()
	err = rootCfg.SaveFile(cmd.pollapoYmlPath)
	if err != nil {
		util.Sugar.Fatalw("Failed to re-write %s with updated root.lock: %s", cmd.pollapoYmlPath, err)
	}
	util.Println("Done.")
}

func (cmd cmdInstall) installDepsRecursive(rootCfg *pollapo.PollapoConfig) {
	putDepIntoMap := func(depsMap map[string]map[string][]string, dep pollapo.PollapoDep, origin string) {
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

	depHandleQueue := []pollapo.PollapoDep{}
	for _, dep := range (*rootCfg).GetDeps(cmd.verbose) {
		util.PrintfVerbose(logName, cmd.verbose, "Enqueue %s.\n", util.Yellow(dep))
		depHandleQueue = append(depHandleQueue, dep)
	}
	depsMap := map[string]map[string][]string{} // depsMap[user/repo][ref]=froms
	origin := "<root>"
	for len(depHandleQueue) > 0 {
		// cache zips concurrently
		wg := sync.WaitGroup{}
		wg.Add(len(depHandleQueue))
		for _, dep := range depHandleQueue {
			// cache zip bin of the dep
			go func(dep pollapo.PollapoDep) {
				cmd.getZip(dep)
				wg.Done()
			}(dep)
		}
		wg.Wait()

		queue := []pollapo.PollapoDep{}
		for _, dep := range depHandleQueue {
			// TODO: froms are unused. command 'why' will use it maybe.
			putDepIntoMap(depsMap, dep, origin)
			zipReader := cmd.getZip(dep)

			// read pollapo.yml & enqueue deps
			pollapoFile := myzip.GetFileByName(zipReader, "pollapo.yml")
			if pollapoFile != nil {
				// get pollapo config
				rc, err := pollapoFile.Open()
				if err != nil {
					util.Sugar.Fatalw("Failed to open pollapo file", err)
				}
				bin, err := io.ReadAll(rc)
				rc.Close()
				if err != nil {
					util.Sugar.Fatalw("Failed to read pollapo file", err)
				}
				depCfg := pollapo.ParsePollapo(bin)
				for _, dep := range depCfg.GetDeps(cmd.verbose) {
					queue = append(queue, dep)
					util.PrintfVerbose(logName, cmd.verbose, "Enqueue %s.\n", util.Yellow(dep))
				}
			}

			origin = dep.String()
		}
		depHandleQueue = queue
	}

	latestRef := func(refs RefArray) string {
		sortedRefs := refs
		sort.Sort(sortedRefs)
		return refs[len(refs)-1]
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
	// sort it to keep consistent installation order
	sort.Strings(latestDeps)

	for _, depTxt := range latestDeps {
		dep, isOk := pollapo.ParseDep(depTxt)
		if !isOk {
			util.Sugar.Fatalw("Failed to parse dep", nil, "dep", depTxt)
		}
		commit, err := cmd.gc.GetCommit(dep.Owner, dep.Repo, dep.Ref)
		if err != nil {
			util.Sugar.Fatalw("Failed to get commit: %s", dep)
		}
		if commit[:len(dep.Ref)] != dep.Ref {
			// if dep.Ref is not commit hash
			lockedRef, found := (*rootCfg).GetLock(dep)
			if !found {
				commit, err := cmd.gc.GetCommit(dep.Owner, dep.Repo, dep.Ref)
				if err == nil {
					(*rootCfg).SetLock(dep, commit)
					dep.Ref = commit
				}
			} else {
				dep.Ref = lockedRef
			}
		}
		depOutDir := filepath.Join(cmd.outDir, dep.Owner, dep.Repo)
		zipReader := cmd.getZip(dep)
		util.Printf("Installing %s...", util.Yellow(dep.String()))
		cmd.uz.Unzip(zipReader, depOutDir)
		util.Print("ok\n")
	}
}

func (cmd cmdInstall) getZip(dep pollapo.PollapoDep) *zip.Reader {
	// Use commit is branch or tag rather than commit sha1.
	zipUrl, err := cmd.gc.GetZipLink(dep.Owner, dep.Repo, dep.Ref)
	if err != nil {
		util.Printf("%s\n", util.Red("error"))
		util.Printf("Login required. (%s): %s\n", util.Yellow(dep), err)
		os.Exit(1)
	}
	zipReader, _ := cmd.zd.GetZip(zipUrl)
	return zipReader
}

type RefArray []string

func (refs RefArray) Len() int {
	return len(refs)
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
