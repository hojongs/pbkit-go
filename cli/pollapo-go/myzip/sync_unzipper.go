package myzip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
)

type SyncUnzipper struct{}

func (uz SyncUnzipper) Unzip(zipReader *zip.Reader, outDir string) {
	if zipReader == nil {
		log.Fatalw("zipReader is null", nil)
	}
	for _, f := range zipReader.File[1:] {
		if f.FileInfo().IsDir() {
			continue
		}
		i := strings.Index(f.Name, "/")
		fname := f.Name[i+1:]
		// log.Infow("Unzip", "filepath", fname)
		fpath := filepath.Join(outDir, fname)
		if !strings.HasPrefix(fpath, filepath.Clean(outDir)+string(os.PathSeparator)) {
			return
		}
		copyFile(f, fpath)
	}
}

func copyFile(f *zip.File, fpath string) {
	if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
		log.Fatalw("Failed to unzip", err)
	}
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
}
