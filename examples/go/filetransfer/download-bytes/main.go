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
	target := filetransfer.Bytes()

	result, err := downloader.Download(
		ctx,
		requiredEnv("DROPBOX_DOWNLOAD_PATH"),
		target,
		filetransfer.DownloadOptions{},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Downloaded %s into memory (%d bytes)\n", result.Metadata.PathDisplay, len(target.Bytes()))
}

func requiredEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(name + " is required")
	}
	return value
}
