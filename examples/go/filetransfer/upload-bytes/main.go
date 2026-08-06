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

	data := []byte("Hello from the Dropbox Go SDK file transfer helper.\n")
	result, err := uploader.Upload(
		ctx,
		filetransfer.BytesUpload(data),
		files.NewCommitInfo(requiredEnv("DROPBOX_UPLOAD_PATH")),
		filetransfer.UploadOptions{},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Uploaded %s (%d bytes)\n", result.Metadata.PathDisplay, result.Metadata.Size)
}

func requiredEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(name + " is required")
	}
	return value
}
