// server_test.go
package server_test

import (
	"Golang-Proxy-Youtube/internal/server"
	proto "Golang-Proxy-Youtube/proto"
	"context"
	"testing"
)

func TestServer_GetThumbnail(t *testing.T) {
	// Создаем новый сервер с моковым кэшем.
	cache := server.NewMemoryCache()
	s := server.NewServer(cache)

	videoID := "test_video_id"
	expectedURL := "https://img.youtube.com/vi/test_video_id/0.jpg"

	// Добавляем URL в кэш для теста.
	cache.Set(videoID, expectedURL)

	// Создаем запрос к серверу для тестирования.
	req := &proto.ThumbnailRequest{VideoId: videoID}
	resp, err := s.GetThumbnail(context.Background(), req)
	if err != nil {
		t.Fatalf("ошибка вызова GetThumbnail: %v", err)
	}

	// Проверяем ответ сервера.
	if resp.ThumbnailUrl != expectedURL {
		t.Errorf("Ожидали %s, получили %s", expectedURL, resp.ThumbnailUrl)
	}
}
