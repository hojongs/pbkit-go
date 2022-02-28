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

func Install(
	clean bool,
	outDir string,
	token string,
	pollapoYmlPath string,
) {
	log.Infow("Params", "clean", clean, "outDir", outDir, "token", token, "config", pollapoYmlPath)

	if clean {
		fmt.Printf("Clean cache root: %s\n", color.Yellow(cache.CacheRoot))
		cache.Clean()
	}

	pollapoBytes, err := os.ReadFile(pollapoYmlPath)
	if err != nil {
		log.Fatalw("Failed to read file", "filename", pollapoYmlPath, "cause", err.Error())
	}

	cfg := pollapo.ParsePollapo(pollapoBytes)
	log.Infow("LoadPollapoYml", "pollapoYml", cfg)
	// install deps in cfg
	q := []string{}
	q = append(q, cfg.Deps...)
	for len(q) > 0 {
		depTxt := q[0]
		q = q[1:]

		dep, isOk := pollapo.ParseDep(depTxt)
		if !isOk {
			log.Fatalw("Invalid dep", nil, "dep", depTxt)
		}

		cacheKey := fmt.Sprintf("%v-%v-%v.zip", dep.Owner, dep.Repo, dep.Ref)
		zipBin, err := cache.Get(cacheKey)
		if err != nil {
			// TODO: color print
			fmt.Printf("Cache not found of %s\n", color.Yellow(cacheKey))
			// TODO: github authentication with pollapo login
			// TODO: github authentication with token
			zipUrl := github.GetZipLink(dep)
			fmt.Printf("Downloading %s...", color.Yellow(depTxt))
			resp, err := http.Get(zipUrl)
			if err != nil {
				log.Fatalw("Failed to HTTP Get", err, "dep", dep)
			}
			if resp.StatusCode != 200 {
				log.Fatalw("HTTP Response is not OK", nil, "status", resp.StatusCode)
			}
			zipBin, err = io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalw("Failed to Read HTTP Response body", err, "body", zipBin[:1024])
			}
			defer resp.Body.Close()
			fmt.Print("ok\n")
		} else {
			fmt.Printf("Use cache of %s.\n", color.Yellow(depTxt))
		}

		depOutDir := filepath.Join(outDir, dep.Owner, dep.Repo)
		fmt.Printf("Installing %s...", color.Yellow(depTxt))
		zip.Unzip(zipBin, depOutDir)
		fmt.Print("ok\n")

		// TODO: enqueue nested deps

		cache.Store(cacheKey, zipBin)
	}

	// getToken
	// backoff (validateToken)
	// cacheDir
	// cacheDeps
	// lockTable
	// analyzeDeps(cacheDir, pollapoYml)
	// *emptyDir
	// *recursive installDep
	// stringify sanitizeDeps
	// writeFile
	//
}
