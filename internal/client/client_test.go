package client_test

import (
	"Golang-Proxy-Youtube/internal/client"
	"Golang-Proxy-Youtube/pkg/cache"
	"testing"
	"time"
)

func TestClient_GetThumbnail(t *testing.T) {
	mockCache := cache.NewMemoryCache()

	c := &client.Client{}

	videoID := "test_video_id"
	expectedURL := "https://img.youtube.com/vi/test_video_id/0.jpg"
	err := mockCache.Set("thumbnail:"+videoID, []byte(expectedURL), 24*time.Hour)
	if err != nil {
		t.Fatalf("не удалось сохранить в кэш: %v", err)
	}

	url, err := c.GetThumbnail(videoID)
	if err != nil {
		t.Fatalf("ошибка получения превью: %v", err)
	}

	if url != expectedURL {
		t.Errorf("Ожидали %s, получили %s", expectedURL, url)
	}
}
