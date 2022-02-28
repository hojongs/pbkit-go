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

func NewFileSystemCache() FileSystemCache {
	return FileSystemCache{rootDir: GetDefaultCacheRoot()}
}

func (cache FileSystemCache) GetRootLocation() string {
	return cache.rootDir
}

func GetDefaultCacheRoot() string {
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
	path := filepath.Join(cache.rootDir, key)
	log.Infow("[Cache] Store", "key", key, "path", path)
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		log.Fatalw("Failed to Write cache file", err, "path", path)
	}
}

func (cache FileSystemCache) Get(key string) ([]byte, error) {
	path := filepath.Join(cache.rootDir, key)
	log.Infow("[Cache] Get", "key", key, "path", path)
	barr, err := os.ReadFile(path)
	if err != nil {
		log.Infow("[Cache] Get Miss", "key", key, "path", path)
		return nil, err
	} else {
		log.Infow("[Cache] Get Hit", "key", key, "path", path)
		return barr, nil
	}
}
