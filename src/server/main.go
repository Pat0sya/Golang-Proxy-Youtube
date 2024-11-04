package main

import (
	"Golang-Proxy-Youtube/internal/server"
	"log"
	"os"
)

func main() {
	port := "localhost:50051"
	cache := server.NewMemoryCache()
	if err := server.Start(port, cache); err != nil {
		log.Fatalf("Ошибка запуска gRPC сервера: %v", err)
		os.Exit(1)
	}
}
