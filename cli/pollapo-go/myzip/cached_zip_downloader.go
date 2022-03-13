package myzip

import (
	"archive/zip"
	"net/url"
	"path"
	"sync"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
	"github.com/hojongs/pbkit-go/cli/pollapo-go/util"
	"github.com/patrickmn/go-cache"
)

type CachedZipDownloader struct {
	Default       ZipDownloader
	dlChanMap     map[string]chan *[]byte // download progress channel
	dlChanMtx     *sync.Mutex
	cache         *cache.Cache
	cacheFilepath string
	verbose       bool // TODO; replace it with verbose printer (verbose check, print prefix)
	onCacheMiss   func(ref string)
	onCacheStore  func(ref string)
	onCacheHit    func(ref string)
}

var mu = sync.Mutex{}
var rwmu = sync.RWMutex{}
var logName = "Zip"

func NewCachedZipDownloader(cacheDir string, verbose bool, onCacheMiss func(ref string), onCacheStore func(ref string), onCacheHit func(ref string)) ZipDownloader {
	cacheFilepath := path.Join(cacheDir, "zip-cache")
	util.PrintfVerbose(logName, verbose, "Loading zip-cache from %s...\n", util.Yellow(cacheFilepath))
	c := util.LoadCache(cacheFilepath, &mu)
	util.PrintfVerbose(logName, verbose, "Loaded zip-cache.\n")
	return CachedZipDownloader{
		NewZipDownloader(),
		make(map[string]chan *[]byte),
		&sync.Mutex{},
		c,
		cacheFilepath,
		verbose,
		onCacheMiss,
		onCacheStore,
		onCacheHit,
	}
}

// cache binary: require enough memory
// alternative: flush cached binary to file system if not enough memory
func (zd CachedZipDownloader) GetZip(zipUrl string) (*zip.Reader, []byte) {
	// parse ref from "GitHub" zipUrl
	u, err := url.Parse(zipUrl)
	cacheKey := u.Path
	if err != nil {
		log.Sugar.Fatalw("Failed to parse URL", err, "u.Path", u.Path)
	}

	b, found := zd.cache.Get(cacheKey)
	if found {
		// Ref is too long...
		// if zd.onCacheHit != nil {
		// 	zd.onCacheHit(u.Path)
		// }
		zipBin := b.([]byte)
		r := NewZipReader(zipBin)
		return r, zipBin
	} else {
		util.Printf("VERBOSE[Zip]: Cache miss %s\n", util.Yellow(u.Path))
		zd.dlChanMtx.Lock()
		if zd.dlChanMap[zipUrl] == nil {
			ch := make(chan *[]byte, 1) // channel should be buffered to avoid blocking
			zd.dlChanMap[zipUrl] = ch
			zd.dlChanMtx.Unlock()
			if zd.onCacheMiss != nil {
				zd.onCacheMiss(u.Path)
			}
			reader, zipBin := zd.Default.GetZip(zipUrl)
			if zd.onCacheStore != nil {
				zd.onCacheStore(u.Path)
			}
			zd.cache.Set(cacheKey, zipBin, cache.DefaultExpiration)
			ch <- &zipBin // it's done to store cache for the key {zipUrl}
			util.PrintfVerbose(logName, zd.verbose, "Sent zipBin to ch %s\n", util.Yellow(u.Path))
			close(ch)
			return reader, zipBin
		} else {
			ch := zd.dlChanMap[zipUrl]
			zd.dlChanMtx.Unlock()
			util.PrintfVerbose(logName, zd.verbose, "VERBOSE[Zip]: Wait ch %s\n", util.Yellow(u.Path))
			zipBinPtr := <-ch
			if zipBinPtr != nil {
				util.PrintfVerbose(logName, zd.verbose, "VERBOSE[Zip]: Get zipBin from ch %s\n", util.Yellow(u.Path))
				zipBin := *zipBinPtr
				return NewZipReader(zipBin), zipBin
			} else {
				// another receiver of the channel already took it
				util.PrintfVerbose(logName, zd.verbose, "VERBOSE[Zip]: Get zipBin from cache instead of ch %s\n", util.Yellow(u.Path))
				b, found := zd.cache.Get(cacheKey)
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
	util.PrintfVerbose(logName, zd.verbose, "Save zip-cache to %s...\n", util.Yellow(zd.cacheFilepath))
	err := util.SaveCache(zd.cache, zd.cacheFilepath, &rwmu)
	if err != nil {
		util.Printf("%s: %s\n", util.Red("Failed to save zip-cache"), err)
	}
	return nil
}
