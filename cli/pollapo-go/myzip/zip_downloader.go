package myzip

import (
	"archive/zip"
	"io"
	"net/http"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/util"
)

type ZipDownloader interface {
	// returns zip reader with zip data binary
	GetZip(zipUrl string) (*zip.Reader, []byte)
	Flush() error
}

type DefaultZipDownloader struct{}

func NewZipDownloader() ZipDownloader {
	return DefaultZipDownloader{}
}

func (zd DefaultZipDownloader) GetZip(zipUrl string) (*zip.Reader, []byte) {
	resp, err := http.Get(zipUrl)
	if err != nil {
		util.Sugar.Fatalw("Failed to HTTP Get", err, "zipUrl", zipUrl)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		util.Sugar.Fatalw("HTTP Response is not OK", nil, "status", resp.StatusCode)
	}
	zipBin := readAll(resp.Body)
	zipReader := NewZipReader(zipBin)

	return zipReader, zipBin
}

func (zd DefaultZipDownloader) Flush() error { return nil }

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
		util.Sugar.Fatalw("Failed to Read HTTP Response body", err, "body", zipBin[:1024])
	}
	return zipBin
}
