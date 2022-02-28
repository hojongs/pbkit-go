package cache

type Cache interface {
	GetRootLocation() string
	Clean()
	Store(key string, data []byte)
	Get(key string) ([]byte, error)
}

type EmptyCache struct{}

func (cache EmptyCache) GetRootLocation() string        { return "" }
func (cache EmptyCache) Clean()                         {}
func (cache EmptyCache) Store(key string, data []byte)  {}
func (cache EmptyCache) Get(key string) ([]byte, error) { return nil, nil }
