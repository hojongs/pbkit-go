package myzip

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
)

type Unzipper interface {
	Unzip(zipReader *zip.Reader, outDir string)
}

type UnzipperImpl struct{}

func (uz UnzipperImpl) Unzip(zipReader *zip.Reader, outDir string) {
	c := make(chan int)
	l := len(zipReader.File[1:])
	for _, f := range zipReader.File[1:] {
		i := strings.Index(f.Name, "/")
		// log.Infow("Unzip", "filepath", f.Name[i+1:])
		fpath := filepath.Join(outDir, f.Name[i+1:])
		if !strings.HasPrefix(fpath, filepath.Clean(outDir)+string(os.PathSeparator)) {
			return
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			log.Fatalw("Failed to unzip", err)
		}

		// ch <- go ff(f, fpath)
		go func() {
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				log.Fatalw("Failed to unzip", err)
			}
			defer outFile.Close()

			rc, err := f.Open()
			if err != nil {
				log.Fatalw("Failed to unzip", err)
			}
			defer rc.Close()
			_, err = io.Copy(outFile, rc)

			if err != nil {
				log.Fatalw("Failed to unzip", err)
			}

			c <- 0
		}()
	}

	i := 0
	for range c {
		i += 1
		if i == l {
			close(c)
		}
	}
}

func NewZipReader(zipBin []byte) *zip.Reader {
	zipReader, err := zip.NewReader(bytes.NewReader(zipBin), int64(len(zipBin)))
	if err != nil {
		log.Fatalw("Read zip", err)
	}
	return zipReader
}
