package myzip

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/util"
)

type Unzipper interface {
	Unzip(zipReader *zip.Reader, outDir string)
}

type UnzipperImpl struct{}

func GetFileByName(zipReader *zip.Reader, match string) *zip.File {
	for _, f := range zipReader.File[1:] {
		i := strings.Index(f.Name, "/")
		fname := f.Name[i+1:]
		if fname == match {
			return f
		}
	}
	return nil
}

func Open(f *zip.File) (io.ReadCloser, error) {
	return f.Open()
}

func (uz UnzipperImpl) Unzip(zipReader *zip.Reader, outDir string) {
	for _, f := range zipReader.File[1:] {
		// validate file path prefix
		i := strings.Index(f.Name, "/")
		fname := f.Name[i+1:]
		// log.Infow("Unzip", "filepath", fname)
		fpath := filepath.Join(outDir, fname)
		if !strings.HasPrefix(fpath, filepath.Clean(outDir)+string(os.PathSeparator)) {
			util.Sugar.Fatalw("Failed to unzip: invalid path", nil, "path", fpath)
		}

		// mkdir
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			util.Sugar.Fatalw("Failed to unzip", err)
		}

		SaveUnzippedFile(f, fpath)
	}
}

func NewZipReader(zipBin []byte) *zip.Reader {
	zipReader, err := zip.NewReader(bytes.NewReader(zipBin), int64(len(zipBin)))
	if err != nil {
		util.Sugar.Fatalw("Read zip", err)
	}
	return zipReader
}

func SaveUnzippedFile(f *zip.File, fpath string) {
	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		util.Sugar.Fatalw("Failed to unzip", err)
	}
	defer outFile.Close()
	rc, err := f.Open()
	if err != nil {
		util.Sugar.Fatalw("Failed to unzip", err)
	}
	defer rc.Close()
	_, err = io.Copy(outFile, rc)
	if err != nil {
		util.Sugar.Fatalw("Failed to unzip", err)
	}
}
