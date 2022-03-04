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
	UnzipFilter(zipReader *zip.Reader, outDir string, filter string)
}

type UnzipperImpl struct{}

func (uz UnzipperImpl) UnzipFilter(zipReader *zip.Reader, outDir string, filter string) {
	for _, f := range zipReader.File[1:] {
		// validate file path prefix
		i := strings.Index(f.Name, "/")
		fname := f.Name[i+1:]
		// log.Infow("Unzip", "filepath", fname)
		fpath := filepath.Join(outDir, fname)
		if !strings.HasPrefix(fpath, filepath.Clean(outDir)+string(os.PathSeparator)) {
			log.Fatalw("Failed to unzip: invalid path", nil, "path", fpath)
		}

		// filter filename to unzip
		if filter != "" && fname != filter {
			continue
		}

		// mkdir
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			log.Fatalw("Failed to unzip", err)
		}

		// save unzipped file
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Fatalw("Failed to unzip", err)
		}
		rc, err := f.Open()
		if err != nil {
			log.Fatalw("Failed to unzip", err)
		}
		_, err = io.Copy(outFile, rc)
		// Close the file without defer so that
		// it closes the outfile before the loop
		// moves to the next iteration. this kinda
		// saves an iteration of memory & time in
		// the worst case scenario.
		outFile.Close()
		rc.Close()
		if err != nil {
			log.Fatalw("Failed to unzip", err)
		}
	}
}

func (uz UnzipperImpl) Unzip(zipReader *zip.Reader, outDir string) {
	uz.UnzipFilter(zipReader, outDir, "")
}

func NewZipReader(zipBin []byte) *zip.Reader {
	zipReader, err := zip.NewReader(bytes.NewReader(zipBin), int64(len(zipBin)))
	if err != nil {
		log.Fatalw("Read zip", err)
	}
	return zipReader
}
