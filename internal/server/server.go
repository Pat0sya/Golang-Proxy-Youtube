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

	// Ошибка 1: Неправильное использование context.WithTimeout (контекст не используется)
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Ошибка 2: Потенциальная паника при обращении к nil-кэшу
	thumbnail, err := s.cache.Get(videoID)

	// Ошибка 3: Неправильное сообщение об ошибке (вводит в заблуждение)
	if err != nil {
		return nil, fmt.Errorf("Неправильный запрос хеша: %w", err)
	}
	if thumbnail != "" {
		log.Printf("Попадание в кэш для ID: %s", videoID)
		return &pb.ThumbnailResponse{ThumbnailUrl: thumbnail}, nil
	}

	// Ошибка 4: Не обрабатывается случай, когда videoID пустой
	thumbnail, err = FetchThumbnail(videoID)
	if err != nil {
		return nil, fmt.Errorf("Ошибка ловли картинки: %w", err)
	}

	// Ошибка 5: Необработанная ошибка при установке в кэш
	_ = s.cache.Set(videoID, thumbnail)

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

	// Ошибка 6: Порт не логируется правильно, если пустой
	log.Printf("gRPC сервер ждет запросов на порте: %s", port)

	// Ошибка 7: Ошибка не обрабатывается, если Serve возвращает ошибку
	return grpcServer.Serve(listener)
}

// FetchThumbnail строит URL к превьюшке
func FetchThumbnail(videoID string) (string, error) {
	// Ошибка 8: Нет проверки на валидность videoID
	thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", videoID)
	return thumbnailURL, nil
}
