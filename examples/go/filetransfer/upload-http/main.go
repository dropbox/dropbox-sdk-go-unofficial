package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/filetransfer"
)

func main() {
	ctx := context.Background()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requiredEnv("SOURCE_URL"), nil)
	if err != nil {
		panic(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		response.Body.Close()
		panic(fmt.Sprintf("source returned %s", response.Status))
	}

	source, err := filetransfer.ReaderUpload(response.Body)
	if err != nil {
		response.Body.Close()
		panic(err)
	}

	client := files.NewContext(dropbox.Config{Token: requiredEnv("DROPBOX_ACCESS_TOKEN")})
	uploader := filetransfer.NewUploader(client)
	result, err := uploader.Upload(
		ctx,
		source,
		files.NewCommitInfo(requiredEnv("DROPBOX_UPLOAD_PATH")),
		filetransfer.UploadOptions{
			Progress: func(progress filetransfer.UploadProgress) {
				if progress.TotalBytes < 0 {
					fmt.Printf("\rUploaded %d bytes", progress.BytesCommitted)
					return
				}
				fmt.Printf("\rUploaded %d of %d bytes", progress.BytesCommitted, progress.TotalBytes)
			},
		},
	)
	fmt.Println()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Uploaded to %s\n", result.Metadata.PathDisplay)
}

func requiredEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(name + " is required")
	}
	return value
}
