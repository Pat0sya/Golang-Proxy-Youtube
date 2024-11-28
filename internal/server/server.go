package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "Golang-Proxy-Youtube/proto"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedThumbnailServiceServer
	cache Cache
}

// NewServer создаёт новый сервер с переданным кэшем
func NewServer(cache Cache) *Server {
	return &Server{cache: cache}
}

// GetThumbnail возвращает URL превью, с проверкой кэша, естественно
func (s *Server) GetThumbnail(ctx context.Context, req *pb.ThumbnailRequest) (*pb.ThumbnailResponse, error) {
	videoID := req.GetVideoId()
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	// Ищем в кэше превьюшку.
	thumbnail, err := s.cache.Get(videoID)
	//Хренова туча проверок
	if err != nil {
		return nil, fmt.Errorf("Ошибка в получении хеша: %w", err)
	}
	if thumbnail != "" {
		log.Printf("Попадание в кэш для ID: %s", videoID)
		return &pb.ThumbnailResponse{ThumbnailUrl: thumbnail}, nil
	}
	thumbnail, err = FetchThumbnail(videoID)
	if err != nil {
		return nil, fmt.Errorf("Ошибка ловли картинки: %w", err)

	}
	err = s.cache.Set(videoID, thumbnail)
	if err != nil {
		log.Printf("Ошибка сохранения кэша: %v", err)
	}
	return &pb.ThumbnailResponse{ThumbnailUrl: thumbnail}, nil
}

// Start запускает gRPC сервер на заданном порту
func Start(port string, cache Cache) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("Не получен ответ: %w", err)

	}
	grpcServer := grpc.NewServer()
	pb.RegisterThumbnailServiceServer(grpcServer, NewServer(cache))
	log.Printf("gRPC сервер ждет запросов на порте: %s", port)
	return grpcServer.Serve(listener)
}

// FetchThumbnail строит URL к превьюшке,
func FetchThumbnail(videoID string) (string, error) {
	thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", videoID)
	return thumbnailURL, nil
}
