package main

import (
	"Golang-Proxy-Youtube/internal/client"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	async := flag.Bool("async", false, "скачиваются превью асихронно")
	outputDir := flag.String("output-dir", ".", "директория для скачивания")
	flag.Parse()

	videoIDs := flag.Args()
	if len(videoIDs) == 0 {
		fmt.Println("Пожалуйства напишите хотябы одно ID видео.")
		os.Exit(1)
	}
	if err := os.MkdirAll(*outputDir, os.ModePerm); err != nil {
		log.Fatalf("Ошибка в создании директории: %v", err)
	}
	serverAddr := "localhost:50051"
	c, err := client.NewClient(serverAddr)
	if err != nil {
		log.Fatalf("Ошибка создании клиента: %v", err)

	}
	defer c.Close()
	if *async {
		results := c.GetThumbnailAsync(videoIDs)
		for id, url := range results {
			if url == "" {

				fmt.Printf("Ошибка в получении превью видео ID: %s", id)
				continue
			}
			if err := client.DownloadThumbnail(url, id, *outputDir); err != nil {
				log.Printf("Ошибка в скачивании видео с ID %s: %v", id, err)
			}

		}
	} else {
		for _, videoID := range videoIDs {
			url, err := c.GetThumbnail(videoID)
			if err != nil {
				log.Printf("Ошибка в получении превью для %s: %v", videoID, err)
			}
			if err := client.DownloadThumbnail(url, videoID, *outputDir); err != nil {
				fmt.Printf("Ошибка в скачивании для видео с ID %s: %v", videoID, err)
			}
		}
	}
}
