package cache

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
)

type FileSystemCache struct {
	rootDir string
}

func NewCache() FileSystemCache {
	return FileSystemCache{
		rootDir: initCacheRoot(),
	}
}

func (cache FileSystemCache) GetRootLocation() string {
	return cache.rootDir
}

func initCacheRoot() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalw("UserHomeDir", err)
	}
	return fmt.Sprintf("%v/.cache/pollapo-go", homeDir)
}

func (cache FileSystemCache) Clean() {
	os.RemoveAll(cache.rootDir)
}

func (cache FileSystemCache) Store(key string, data []byte) {
	if _, err := os.Stat(cache.rootDir); os.IsNotExist(err) {
		os.MkdirAll(cache.rootDir, 0755)
	}
	dst := filepath.Join(cache.rootDir, key)
	err := os.WriteFile(dst, data, 0644)
	if err != nil {
		log.Fatalw("Failed to Write cache file", err, "dst", dst)
	}
}

func (cache FileSystemCache) Get(key string) ([]byte, error) {
	barr, err := os.ReadFile(filepath.Join(cache.rootDir, key))
	if err != nil {
		return nil, err
	} else {
		return barr, nil
	}
}
