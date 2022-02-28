package zip

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
)

// TODO: Unzipper interface define & impl the interface & pass the impl to install.

func Unzip(barr []byte, outDir string) {
	os.MkdirAll(outDir, 0755)
	zipReader, err := zip.NewReader(bytes.NewReader(barr), int64(len(barr)))
	if err != nil {
		log.Fatalw("Read zip", err)
	}

	for _, file := range zipReader.File {
		filename := filepath.Base(file.Name)
		// log.Infow("Reading file", "filename", filename)
		fileBarr, err := readFileInZip(file)
		if err != nil {
			log.Fatalw("Failed to Read file in zip", err)
		}
		dst := filepath.Join(outDir, filename)
		// TODO: if the directory or file already exists
		err = os.WriteFile(dst, fileBarr, 0644)
		if err != nil {
			log.Fatalw("Failed to Write file from zip", err, "dst", dst)
		}
	}
}

func readFileInZip(zf *zip.File) ([]byte, error) {
	r, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	barr, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return barr, nil
}
