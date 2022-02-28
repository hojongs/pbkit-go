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
	zd ZipDownloader,
) {
	cfg, err := getPollapoConfig(pollapoYmlPath)
	if err != nil {
		log.Fatalw("Failed to read file", err, "filename", pollapoYmlPath)
	}
	InstallConfig(clean, outDir, token, cfg, zd)
}

func getPollapoConfig(pollapoYmlPath string) (pollapo.PollapoConfig, error) {
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

type GitHubZipDownloader struct{}

func (this GitHubZipDownloader) GetZipBin(owner string, repo string, ref string) []byte {
	depTxt := fmt.Sprintf("%s/%s@%v", owner, repo, ref)
	zipUrl := github.GetZipLink(owner, repo, ref)
	fmt.Printf("Downloading %s...", color.Yellow())
	resp, err := http.Get(zipUrl)
	if err != nil {
		log.Fatalw("Failed to HTTP Get", err, "dep", depTxt)
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

func InstallConfig(
	clean bool,
	outDir string,
	token string,
	cfg pollapo.PollapoConfig,
	zd ZipDownloader,
) {
	if clean {
		fmt.Printf("Clean cache root: %s\n", color.Yellow(cache.CacheRoot))
		cache.Clean()
	}

	// install deps in cfg
	queue := []string{}
	queue = append(queue, cfg.Deps...)
	for len(queue) > 0 {
		depTxt := queue[0]
		queue = queue[1:]

		dep, isOk := pollapo.ParseDep(depTxt)
		if !isOk {
			log.Fatalw("Invalid dep", nil, "dep", depTxt)
		}

		// TODO: resolve duplicated deps with comparison refs

		cacheKey := fmt.Sprintf("%v-%v-%v.zip", dep.Owner, dep.Repo, dep.Ref)
		zipBin, err := cache.Get(cacheKey)
		if err != nil {
			// TODO: color print
			fmt.Printf("Cache not found of %s\n", color.Yellow(cacheKey))
			// TODO: github authentication with pollapo login
			// TODO: github authentication with token
			zipBin = zd.GetZipBin(dep.Owner, dep.Repo, dep.Ref)
			fmt.Print("ok\n")
		} else {
			fmt.Printf("Use cache of %s.\n", color.Yellow(depTxt))
		}

		depOutDir := filepath.Join(outDir, dep.Owner, dep.Repo)
		fmt.Printf("Installing %s...", color.Yellow(depTxt))
		zip.Unzip(zipBin, depOutDir)
		fmt.Print("ok\n")

		depPollapoYmlPath := filepath.Join(depOutDir, "pollapo.yml")
		depCfg, err := getPollapoConfig(depPollapoYmlPath)
		if err == nil {
			for _, nestedDep := range depCfg.Deps {
				queue = append(queue, nestedDep)
			}
		}
		cache.Store(cacheKey, zipBin)
	}
}
