package server

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "Golang-Proxy-Youtube/proto"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedThumbnailServiceServer
	cache Cache
}

func NewServer(cache Cache) *Server {
	return &Server{cache: cache}
}

func (s *Server) GetThumbnail(ctx context.Context, req *pb.ThumbnailRequest) (*pb.ThumbnailResponse, error) {
	videoID := req.GetVideoId()
	thumbnail, err := s.cache.Get(videoID)
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

func FetchThumbnail(videoID string) (string, error) {
	thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", videoID)
	return thumbnailURL, nil
}
