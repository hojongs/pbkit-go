package myzip

import (
	"archive/zip"
	"bytes"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
)

type Unzipper interface {
	Unzip(zipReader *zip.Reader, outDir string)
}

func NewZipReader(zipBin []byte) *zip.Reader {
	zipReader, err := zip.NewReader(bytes.NewReader(zipBin), int64(len(zipBin)))
	if err != nil {
		log.Fatalw("Read zip", err)
	}
	return zipReader
}
