// Пакет client реализует клиента для работы с gRPC сервером, который предоставляет
// сервис для получения превью (thumbnail) видео с YouTube, с возможностью кэширования изображений в Redis.
package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"Golang-Proxy-Youtube/pkg/cache"
	pb "Golang-Proxy-Youtube/proto"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.ThumbnailServiceClient
	cache  *cache.RedisCache
}

// NewClient создаёт нового клиента для gRPC-сервиса
func NewClient(serverAddr string) (*Client, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials())) //Что уже deped?! Прочитать!
	if err != nil {
		return nil, fmt.Errorf("Ошибка при подключении к gRPC серверу: %w", err)
	}

	client := pb.NewThumbnailServiceClient(conn)
	redisAddr := "localhost:6379"
	redisCache := cache.NewRedisCache(redisAddr, "", 0)

	return &Client{conn: conn, client: client, cache: redisCache}, nil

}

// Close закрывает соединения gRPC и Redis
func (c *Client) Close() {
	c.conn.Close()
	c.cache.Close()
}

// GetThumbnail запрашивает URL превью у сервиса
func (c *Client) GetThumbnail(videoID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pb.ThumbnailRequest{VideoId: videoID}
	res, err := c.client.GetThumbnail(ctx, req)
	if err != nil {
		return "", fmt.Errorf("Ошибка в получении превью: %w", err)

	}
	return res.GetThumbnailUrl(), nil
}

// GetThumbnailAsync параллельно запрашивает URL-ы превьюшек для нескольких видео
func (c *Client) GetThumbnailAsync(videoIDs []string) map[string]string {
	var wg sync.WaitGroup
	results := make(map[string]string)
	//Рутины!
	mu := &sync.Mutex{}
	for _, videoID := range videoIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			url, err := c.GetThumbnail(id)
			mu.Lock()
			if err != nil {
				results[id] = fmt.Sprintf("Ошибка: %v", err)
			} else {
				results[id] = url
			}
			mu.Unlock()
		}(videoID)
	}
	wg.Wait()
	return results

}

/* This is code of the early implementation, feel free to edit and use
 */
//	func DownloadThumbnail(url, videoID, outputDir string) error {
//		response, err := http.Get(url)
//		if err != nil {
//			return fmt.Errorf("Ошибка в скачивании превью: %w", err)
//		}
//		defer response.Body.Close()
//		if response.StatusCode != http.StatusOK {
//			return fmt.Errorf("Ошибка в загрузке превью, статус: %d", response.StatusCode)
//
//		}
//		filename := filepath.Join(outputDir, fmt.Sprintf("%s.jpg", videoID))
//		file, err := os.Create(filename)
//		defer file.Close()
//		_, err = io.Copy(file, response.Body)
//		if err != nil {
//			return fmt.Errorf("Ошибка в загрузке превью в файл: %w", err)
//		}
//
//		fmt.Printf("Превью сохранен в %s\n", filename)
//		return nil
//	}
func (c *Client) DownloadThumbnail(videoID, outputDir string) error {
	cacheKey := fmt.Sprintf("thumbnail:%s", videoID)
	cachedData, err := c.cache.Get(cacheKey)
	if err == nil && cachedData != nil {
		filename := filepath.Join(outputDir, fmt.Sprintf("%s.jpg", videoID))
		return os.WriteFile(filename, cachedData, 0644)
	}
	url, err := c.GetThumbnail(videoID)
	if err != nil {
		return fmt.Errorf("ошибка в получении превью: %w", err)
	}
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Ошибка в скачивании превью: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Ошибка с скачивании превью, статус код: %d", response.StatusCode)
	}
	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Ошибка в получении данных: %w", err)

	}
	if err := c.cache.Set(cacheKey, imageData, 24*time.Hour); err != nil {
		return fmt.Errorf("Ошибка в кэшировании: %w", err)
	}
	filename := filepath.Join(outputDir, fmt.Sprintf("%s.jpg", videoID))
	return os.WriteFile(filename, imageData, 0644)
}
