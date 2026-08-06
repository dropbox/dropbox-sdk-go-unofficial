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
	uploader := filetransfer.NewUploader(client)

	source, err := filetransfer.FileUpload(requiredEnv("LOCAL_UPLOAD_PATH"))
	if err != nil {
		panic(err)
	}

	_, err = uploader.Upload(
		ctx,
		source,
		files.NewCommitInfo(requiredEnv("DROPBOX_UPLOAD_PATH")),
		filetransfer.UploadOptions{
			Progress: func(progress filetransfer.UploadProgress) {
				fmt.Printf("\rUploaded %d of %d bytes", progress.BytesCommitted, progress.TotalBytes)
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
