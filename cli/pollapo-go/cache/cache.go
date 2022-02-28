package cache

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
)

// TODO: cache interface define & impl the interface & pass the impl to install.

func getCacheRoot() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalw("UserHomeDir", err)
	}
	return fmt.Sprintf("%v/.cache/pollapo-go", homeDir)
}

var CacheRoot = getCacheRoot()

func Clean() {
	os.RemoveAll(CacheRoot)
}

func Store(key string, data []byte) {
	if _, err := os.Stat(CacheRoot); os.IsNotExist(err) {
		os.MkdirAll(CacheRoot, 0755)
	}
	dst := filepath.Join(CacheRoot, key)
	err := os.WriteFile(dst, data, 0644)
	if err != nil {
		log.Fatalw("Failed to Write cache file", err, "dst", dst)
	}
}

func Get(key string) ([]byte, error) {
	barr, err := os.ReadFile(filepath.Join(CacheRoot, key))
	if err != nil {
		return nil, err
	} else {
		return barr, nil
	}
}
