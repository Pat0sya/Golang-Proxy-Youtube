package cache

import (
	"sync"
	"time"
)

type MemoryCache struct {
	data map[string][]byte
	mu   sync.RWMutex
}

// NewMemoryCache инициализирует новый in-memory кэш
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string][]byte),
	}
}

func (m *MemoryCache) Set(key string, value []byte, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	// В продакшене сюда можно добавить логику для очистки по истечению expiration.
	return nil
}

func (m *MemoryCache) Get(key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, exists := m.data[key]
	if !exists {
		return nil, nil // Пустой результат для случаев "ключ не найден"
	}
	return value, nil
}

func (m *MemoryCache) Close() error {
	return nil // В случае с in-memory кэшем делать ничего не нужно
}
