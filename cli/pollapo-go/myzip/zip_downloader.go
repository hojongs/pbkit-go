package myzip

import (
	"archive/zip"
	"io"
	"net/http"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
)

type ZipDownloader interface {
	// returns zip reader with zip data binary
	GetZip(zipUrl string) (*zip.Reader, []byte)
}

type GitHubZipDownloader struct{}

func NewGitHubZipDownloader() GitHubZipDownloader {
	return GitHubZipDownloader{}
}

func (gzd GitHubZipDownloader) GetZip(zipUrl string) (*zip.Reader, []byte) {
	resp, err := http.Get(zipUrl)
	if err != nil {
		log.Fatalw("Failed to HTTP Get", err, "zipUrl", zipUrl)
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
