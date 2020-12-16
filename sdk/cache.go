package sdk

// BaseCache Simple in-memory cache
type BaseCache interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Remove(key string)
	RemoveAll()
}
