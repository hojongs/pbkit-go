package myzip

import (
	"archive/zip"

	"github.com/patrickmn/go-cache"
)

type CachedZipDownloader struct {
	Default ZipDownloader
	dlChan  map[string]chan interface{} // download progress channel
	cache   *cache.Cache
}

func NewCachedZipDownloader() ZipDownloader {
	return CachedZipDownloader{}
}

// cache binary: require enough memory
// alternative: flush cached binary to file system
func (zd CachedZipDownloader) GetZip(zipUrl string) (*zip.Reader, []byte) {
	b, found := zd.cache.Get(zipUrl)
	if found {
		bb := b.([]byte)
		r := NewZipReader(bb)
		return r, bb
	} else {
		r, b := zd.Default.GetZip(zipUrl)
		zd.cache.Set(zipUrl, b, cache.DefaultExpiration)
		return r, b
	}
}

func (zd CachedZipDownloader) Flush() error {
	// TODO
	return nil
}
