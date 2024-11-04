package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"

	pb "Golang-Proxy-Youtube/proto"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.ThumbnailServiceClient
}

func NewClient(serverAddr string) (*Client, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure()) //Что уже deped?! Прочитать!
	if err != nil {
		return nil, fmt.Errorf("Ошибка при подключении к gRPC серверу: %w", err)
	}
	client := pb.NewThumbnailServiceClient(conn)
	return &Client{conn: conn, client: client}, nil

}
func (c *Client) Close() {
	c.conn.Close()
}
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
func (c *Client) GetThumbnailAsync(videoIDs []string) map[string]string {
	var wg sync.WaitGroup
	results := make(map[string]string)
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
