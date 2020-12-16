package sdk

import (
	"testing"
)

func assertInt(t *testing.T, value int, expected int, message string) {
	if value != expected {
		t.Errorf("%s, got: %d, expected: %d", message, value, expected)
	}
}

func assertString(t *testing.T, value string, expected string, message string) {
	if value != expected {
		t.Errorf("%s, got: %s, expected: %s", message, value, expected)
	}
}

func TestSet(t *testing.T) {
	cache := NewInMemoryCache()

	cache.Set("key1", "value1")
	keyCount := len(cache.data)
	assertInt(t, keyCount, 1, "Key count incorrect")
}

func TestMultipleSet(t *testing.T) {
	cache := NewInMemoryCache()

	cache.Set("key1", "value1")
	cache.Set("key1", "value2")
	keyCount := len(cache.data)
	assertInt(t, keyCount, 1, "Key count incorrect")
	value := cache.Get("key1").(string)
	assertString(t, value, "value2", "Key value incorrect")
}

func TestGet(t *testing.T) {
	cache := NewInMemoryCache()

	cache.Set("key1", "value1")
	value := cache.Get("key1").(string)
	assertString(t, value, "value1", "Key value incorrect")
}

func TestRemove(t *testing.T) {
	cache := NewInMemoryCache()

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Remove("key2")
	keyCount := len(cache.data)
	assertInt(t, keyCount, 1, "Key count incorrect")
	value := cache.Get("key1").(string)
	assertString(t, value, "value1", "Key value incorrect")
	value2 := cache.Get("key2")
	if value2 != nil {
		t.Errorf("Expected nil but found value %s", value2)
	}
}

func TestRemoveAll(t *testing.T) {
	cache := NewInMemoryCache()

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.RemoveAll()
	keyCount := len(cache.data)
	assertInt(t, keyCount, 0, "Key count incorrect")
}
