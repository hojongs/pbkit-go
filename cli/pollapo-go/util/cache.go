package util

import (
	"bytes"
	"encoding/gob"
	"os"
	"sync"

	"github.com/hojongs/pbkit-go/cli/pollapo-go/log"
	"github.com/patrickmn/go-cache"
)

func LoadCache(cacheFilepath string, mu *sync.Mutex) *cache.Cache {
	barr, err := os.ReadFile(cacheFilepath)
	var c *cache.Cache
	if err != nil {
		c = cache.New(cache.NoExpiration, cache.NoExpiration)
	} else {
		// Load cache from bytes
		// https://github.com/patrickmn/go-cache/blob/v2.1.0/cache.go#L1002
		dec := gob.NewDecoder(bytes.NewReader(barr))
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
	return c
}

func SaveCache(cache *cache.Cache, cacheFilepath string, mu *sync.RWMutex) error {
	// Save cache items to file
	// https://github.com/patrickmn/go-cache/blob/v2.1.0/cache.go#L963
	f, err := os.Create(cacheFilepath)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(f)
	defer func() {
		if x := recover(); x != nil {
			log.Sugar.Fatal("Error registering item types with Gob library")
		}
	}()
	items := cache.Items()
	mu.RLock()
	defer mu.RUnlock()
	for _, v := range items {
		gob.Register(v.Object)
	}
	err = enc.Encode(&items)
	if err != nil {
		return err
	}
	return nil
}
