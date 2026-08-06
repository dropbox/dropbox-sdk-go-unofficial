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

	result, err := downloader.Download(
		ctx,
		requiredEnv("DROPBOX_DOWNLOAD_PATH"),
		filetransfer.File(requiredEnv("LOCAL_DOWNLOAD_PATH")),
		filetransfer.DownloadOptions{},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Downloaded %s (%d bytes)\n", result.Metadata.PathDisplay, result.Metadata.Size)
}

func requiredEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(name + " is required")
	}
	return value
}
