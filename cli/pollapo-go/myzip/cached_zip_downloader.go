package myzip

import (
	"archive/zip"
	"bytes"
	"encoding/gob"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/util"
	"github.com/patrickmn/go-cache"
)

type CachedZipDownloader struct {
	Default   ZipDownloader
	dlChanMap map[string]chan *[]byte // download progress channel
	dlChanMtx *sync.Mutex
	cache     *cache.Cache
	verbose   bool // TODO; replace it with verbose printer (verbose check, print prefix)
}

func NewCachedZipDownloader(cacheDir string, verbose bool) ZipDownloader {
	barr, err := os.ReadFile(path.Join(cacheDir, "zip-cache"))
	var c *cache.Cache
	if err != nil {
		c = cache.New(cache.NoExpiration, cache.NoExpiration)
	} else {
		// Load cache from bytes
		// https://github.com/patrickmn/go-cache/blob/v2.1.0/cache.go#L1002
		dec := gob.NewDecoder(bytes.NewReader(barr))
		mu := sync.Mutex{}
		items := map[string]cache.Item{}
		err := dec.Decode(&items)
		if err == nil {
			mu.Lock()
			for k, v := range items {
				ov, found := items[k] // ov = old value
				if !found || ov.Expired() {
					items[k] = v
				}
			}
			mu.Unlock()
		}
		c = cache.NewFrom(cache.NoExpiration, cache.NoExpiration, items)
	}
	return CachedZipDownloader{
		NewZipDownloader(),
		make(map[string]chan *[]byte),
		&sync.Mutex{},
		c,
		verbose,
	}
}

// cache binary: require enough memory
// alternative: flush cached binary to file system if not enough memory
func (zd CachedZipDownloader) GetZip(zipUrl string) (*zip.Reader, []byte) {
	// parse ref from "GitHub" zipUrl
	i := strings.LastIndex(zipUrl, "?")
	var ref string
	if i == -1 {
		ref = util.Yellow(zipUrl[strings.LastIndex(zipUrl, "/")+1:])
	} else {
		ref = util.Yellow(zipUrl[strings.LastIndex(zipUrl, "/")+1 : i])
	}

	b, found := zd.cache.Get(zipUrl)
	if found {
		if zd.verbose {
			util.Println("[Zip] Cache hit", "url", ref)
		}
		zipBin := b.([]byte)
		r := NewZipReader(zipBin)
		return r, zipBin
	} else {
		if zd.verbose {
			util.Println("[Zip] Cache miss", "url", ref)
		}
		zd.dlChanMtx.Lock()
		if zd.dlChanMap[zipUrl] == nil {
			ch := make(chan *[]byte, 1) // channel should be buffered to avoid blocking
			zd.dlChanMap[zipUrl] = ch
			zd.dlChanMtx.Unlock()
			reader, zipBin := zd.Default.GetZip(zipUrl)
			zd.cache.Set(zipUrl, zipBin, cache.DefaultExpiration)
			if zd.verbose {
				util.Println("[Zip] Cache set", "url", ref)
			}
			ch <- &zipBin // it's done to store cache for the key {zipUrl}
			if zd.verbose {
				util.Println("[Zip] Sent zipBin to ch", "url", ref)
			}
			close(ch)
			return reader, zipBin
		} else {
			ch := zd.dlChanMap[zipUrl]
			zd.dlChanMtx.Unlock()
			if zd.verbose {
				util.Println("[Zip] Wait ch", "url", ref)
			}
			zipBinPtr := <-ch
			if zipBinPtr != nil {
				if zd.verbose {
					util.Println("[Zip] Get zipBin from ch", "url", ref)
				}
				zipBin := *zipBinPtr
				return NewZipReader(zipBin), zipBin
			} else {
				// another receiver of the channel already took it
				if zd.verbose {
					util.Println("[Zip] Get zipBin from cache instead of ch", "url", ref)
				}
				b, found := zd.cache.Get(zipUrl)
				if !found {
					log.Fatalw("Unexpected cache miss", nil)
				}
				zipBin := b.([]byte)
				return NewZipReader(zipBin), zipBin
			}
		}
	}
}

func (zd CachedZipDownloader) Flush() error {
	// TODO: Save cache to file
	return nil
}
