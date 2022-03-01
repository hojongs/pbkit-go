package myzip

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/github"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/mycolor"
)

type ZipDownloader interface {
	// returns zip reader with zip data binary
	GetZip(owner string, repo string, ref string) (*zip.Reader, []byte)
}

type GitHubZipDownloader struct {
	client github.Client
}

func NewGitHubZipDownloader(client github.Client) GitHubZipDownloader {
	return GitHubZipDownloader{client}
}

func (gzd GitHubZipDownloader) GetZip(owner string, repo string, ref string) (*zip.Reader, []byte) {
	// TODO: github authentication with token
	zipUrl := gzd.client.GetZipLink(owner, repo, ref)
	fmt.Printf("Downloading %s...", mycolor.Yellow(owner+"/"+repo+"@"+ref))
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
	zipReader := NewZipReader(zipBin)

	return zipReader, zipBin
}
