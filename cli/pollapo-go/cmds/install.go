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
			fmt.Printf("Arguments are not required.\n")
			fmt.Printf("Given arguments: count %v, values %v\n", mycolor.Yellow(c.Args().Len()), mycolor.Yellow(c.Args()))
			os.Exit(1)
		}
		if c.Bool("verbose") {
			fmt.Printf("Flag verbose: %v\n", mycolor.Yellow(c.Bool("verbose")))
			fmt.Printf("Flag clean: %v\n", mycolor.Yellow(c.Bool("clean")))
			fmt.Printf("Flag out-dir: %v\n", mycolor.Yellow(c.String("out-dir")))
			if c.String("token") != "" {
				fmt.Printf("Flag token: %v\n", mycolor.Yellow(c.String("token")))
			}
			fmt.Printf("Flag config: %v\n", mycolor.Yellow(c.String("config")))
		}
		token := c.String("token")
		if token == "" {
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
			c.Bool("verbose"),
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
	verbose        bool
}

func newCmdInstall(
	cleanCache bool,
	outDir string,
	pollapoYmlPath string,
	zd myzip.ZipDownloader,
	uz myzip.Unzipper,
	loader pollapo.ConfigLoader,
	cache cache.Cache,
	verbose bool,
) cmdInstall {
	return cmdInstall{cleanCache, outDir, pollapoYmlPath, zd, uz, loader, cache, verbose}
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
		fmt.Printf("%s not found.\n", mycolor.Red(absPath))
		// TODO: Create absPath?
		os.Exit(1)
	}
	cmd.printfIfVerbose("Clean out directory %s.\n", mycolor.Yellow(cmd.outDir))
	if err := os.RemoveAll(cmd.outDir); err != nil {
		log.Fatalw("Remove out dir", err, "outDir", cmd.outDir)
	}
	cmd.installDepsRecursive(rootCfg)
	fmt.Println("Done.")
}

func (cmd cmdInstall) installDepsRecursive(rootCfg pollapo.PollapoConfig) {
	depHandleQueue := []string{}
	for _, dep := range rootCfg.Deps {
		cmd.printfIfVerbose("Enqueue %s.\n", mycolor.Yellow(dep))
	}
	depHandleQueue = append(depHandleQueue, rootCfg.Deps...)
	depsMap := map[string]map[string][]string{} // depsMap[user/repo][ref]=froms
	origin := "<root>"
	for len(depHandleQueue) > 0 {
		depTxt := depHandleQueue[0]
		depHandleQueue = depHandleQueue[1:]

		dep, isOk := pollapo.ParseDep(depTxt)
		if !isOk {
			log.Fatalw("Invalid dep", nil, "dep", depTxt)
		}

		// TODO: froms are unused. command 'why' will use it maybe.
		putDepIntoMap(depsMap, dep, origin)

		// get dependency zip
		zipBin, err := cmd.cache.Get(cacheKeyOf(dep))
		var zipReader *zip.Reader = nil
		if err != nil || zipBin == nil {
			zipReader = cmd.downloadZip(dep)
		} else {
			fmt.Printf("Use cache of %s.\n", mycolor.Yellow(depTxt))
			zipReader = myzip.NewZipReader(zipBin)
		}

		// unzip pollapo.yml
		cacheOutDir := filepath.Join(cmd.cache.GetRootLocation(), dep.Owner, dep.Repo)
		cmd.uz.UnzipFilter(zipReader, cacheOutDir, "pollapo.yml")

		// get pollapo config
		depPollapoYmlPath := filepath.Join(cacheOutDir, "pollapo.yml")
		depCfg, err := cmd.loader.GetPollapoConfig(depPollapoYmlPath)
		if err != nil {
			cmd.printfIfVerbose("pollapo.yml not found %s\n", mycolor.Yellow(depPollapoYmlPath))
		} else {
			for _, dep := range depCfg.Deps {
				cmd.printfIfVerbose("Enqueue %s.\n", mycolor.Yellow(dep))
			}
			depHandleQueue = append(depHandleQueue, depCfg.Deps...)
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

func (cmd cmdInstall) printfIfVerbose(format string, a ...interface{}) (n int, err error) {
	if cmd.verbose {
		return fmt.Printf(format, a...)
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
