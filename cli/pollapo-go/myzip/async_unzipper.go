package myzip

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
)

type ASyncUnzipper struct{}

func (uz ASyncUnzipper) Unzip(zipReader *zip.Reader, outDir string) {
	if zipReader == nil {
		log.Fatalw("zipReader is null", nil)
	}
	wg := sync.WaitGroup{}
	wg.Add(len(zipReader.File[1:]))
	for _, f := range zipReader.File[1:] {
		if f.FileInfo().IsDir() {
			wg.Done()
			continue
		}
		i := strings.Index(f.Name, "/")
		fname := f.Name[i+1:]
		// log.Infow("Unzip", "filepath", fname)
		fpath := filepath.Join(outDir, fname)
		if !strings.HasPrefix(fpath, filepath.Clean(outDir)+string(os.PathSeparator)) {
			return
		}
		go copyFileAsync(f, fpath, &wg)
	}
	wg.Wait()
}

func copyFileAsync(f *zip.File, fpath string, wg *sync.WaitGroup) {
	copyFile(f, fpath)
	wg.Done()
}
