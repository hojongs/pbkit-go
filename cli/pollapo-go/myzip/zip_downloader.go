package myzip

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/github"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/util"
)

type ZipDownloader interface {
	// returns zip reader with zip data binary
	GetZip(owner string, repo string, ref string) (*zip.Reader, []byte)
}

type GitHubZipDownloader struct {
	// TODO: remove
	client github.GitHubClient
}

func NewGitHubZipDownloader(client github.GitHubClient) GitHubZipDownloader {
	return GitHubZipDownloader{client}
}

func (gzd GitHubZipDownloader) GetZip(owner string, repo string, ref string) (*zip.Reader, []byte) {
	zipUrl, err := gzd.client.GetZipLink(owner, repo, ref)
	if err != nil {
		util.Printf("%s\n", util.Red("error"))
		util.Printf("Login required. (%s/%s@%s)\n", owner, repo, ref)
		os.Exit(1)
	}
	resp, err := http.Get(zipUrl)
	if err != nil {
		log.Fatalw("Failed to HTTP Get", err, "dep", fmt.Sprintf("%s/%s@%v", owner, repo, ref))
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalw("HTTP Response is not OK", nil, "status", resp.StatusCode)
	}
	zipBin := readAll(resp.Body)
	zipReader := NewZipReader(zipBin)

	return zipReader, zipBin
}

func readAll(reader io.Reader) []byte {
	// TODO: print progress if verbose for slow network
	// currSize := int64(0)
	// pr := &ProgressReader{reader, func(r int64) {
	// 	currSize += r
	// 	if r > 0 {
	// 		fmt.Println(currSize)
	// 	} else {
	// 		fmt.Println("Downloaded")
	// 	}
	// }}
	zipBin, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalw("Failed to Read HTTP Response body", err, "body", zipBin[:1024])
	}
	return zipBin
}
