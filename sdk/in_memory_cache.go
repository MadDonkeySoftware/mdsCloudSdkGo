package sdk

// InMemoryCache Simple in-memory cache
type InMemoryCache struct {
	data map[string]interface{}
}

// NewInMemoryCache Creates a new in-memory cache
func NewInMemoryCache() *InMemoryCache {
	cache := InMemoryCache{
		data: make(map[string]interface{}),
	}

	return &cache
}

// Set Adds an item to the in-memory cache
func (c *InMemoryCache) Set(key string, value interface{}) {
	c.data[key] = value
}

// Get Retrieves an item from the in-memory cache
func (c *InMemoryCache) Get(key string) interface{} {
	return c.data[key]
}

// Remove Removes an item from the in-memory cache if it exists
func (c *InMemoryCache) Remove(key string) {
	delete(c.data, key)
}

// RemoveAll Removes all items from the in-memory cache
func (c *InMemoryCache) RemoveAll() {
	c.data = make(map[string]interface{})
}
