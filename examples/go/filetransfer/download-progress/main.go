package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/filetransfer"
)

func main() {
	ctx := context.Background()
	client := files.NewContext(dropbox.Config{Token: requiredEnv("DROPBOX_ACCESS_TOKEN")})
	downloader := filetransfer.NewDownloader(client)

	_, err := downloader.Download(
		ctx,
		requiredEnv("DROPBOX_DOWNLOAD_PATH"),
		filetransfer.File(requiredEnv("LOCAL_DOWNLOAD_PATH")),
		filetransfer.DownloadOptions{
			Progress: func(progress filetransfer.DownloadProgress) {
				fmt.Printf("\rDownloaded %d of %d bytes", progress.BytesCommitted, progress.TotalBytes)
			},
		},
	)
	fmt.Println()
	if err != nil {
		panic(err)
	}
}

func requiredEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(name + " is required")
	}
	return value
}
