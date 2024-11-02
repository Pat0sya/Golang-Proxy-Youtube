package server

import "sync"

type Cache interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

type MemoryCache struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		store: make(map[string]string),
	}
}
func (m *MemoryCache) Get(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if value, found := m.store[key]; found {
		return value, nil
	}
	return "", nil
}

func (m *MemoryCache) Set(key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[key] = value
	return nil
}
