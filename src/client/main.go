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
	flag.Parse()

	videoIDs := flag.Args()
	if len(videoIDs) == 0 {
		fmt.Println("Пожалуйства напишите хотябы одно ID видео.")
		os.Exit(1)
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
			fmt.Printf("VideoID: %s, Thumbnail URL: %s\n", id, url)

		}
	} else {
		for _, videoID := range videoIDs {
			url, err := c.GetThumbnail(videoID)
			if err != nil {
				log.Printf("Ошибка в получении превью для %s: %v", videoID, err)
			} else {
				fmt.Printf("VideoID: %s, Thumbnail URL: %s\n", videoID, url)
			}
		}
	}
}
