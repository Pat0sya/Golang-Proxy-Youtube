// cache_test.go
package cache_test

import (
	"Golang-Proxy-Youtube/pkg/cache"
	"testing"
	"time"
)

func TestRedisCache_SetAndGet(t *testing.T) {
	cache := cache.NewMemoryCache() // Используем in-memory кэш для тестирования.
	defer cache.Close()

	key := "test_key"
	value := []byte("test_value")
	expiration := 1 * time.Hour

	// Устанавливаем значение в кэш.
	if err := cache.Set(key, value, expiration); err != nil {
		t.Fatalf("не удалось установить значение в кэш: %v", err)
	}

	// Извлекаем значение из кэша.
	cachedValue, err := cache.Get(key)
	if err != nil {
		t.Fatalf("ошибка получения значения из кэша: %v", err)
	}

	// Проверяем, что значение соответствует ожидаемому.
	if string(cachedValue) != string(value) {
		t.Errorf("Ожидали %s, получили %s", value, cachedValue)
	}
}
