// Copyright (c) Dropbox, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/filetransfer"
)

func TestUserFileTransferUploadDownloadFile(t *testing.T) {
	ctx := context.Background()
	client := files.NewContext(userConfig(t))
	remotePath := fmt.Sprintf("/sdk-integration-filetransfer-%d.txt", time.Now().UnixNano())
	defer deleteIntegrationPath(t, client, remotePath)

	payload := "hello from the filetransfer integration test\n"
	localUploadPath := filepath.Join(t.TempDir(), "upload.txt")
	if err := os.WriteFile(localUploadPath, []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}

	uploadSource, err := filetransfer.FileUpload(localUploadPath)
	if err != nil {
		t.Fatal(err)
	}

	var uploadUpdates []filetransfer.UploadProgress
	uploadResult, err := filetransfer.NewUploader(client).Upload(
		ctx,
		uploadSource,
		files.NewCommitInfo(remotePath),
		filetransfer.UploadOptions{
			Progress: func(progress filetransfer.UploadProgress) {
				uploadUpdates = append(uploadUpdates, progress)
			},
		},
	)
	if err != nil {
		t.Fatalf("Upload(%q): %v", remotePath, err)
	}
	assertUploadProgress(t, uploadUpdates, int64(len(payload)))
	if uploadResult.Metadata == nil {
		t.Fatal("Upload() metadata is nil")
	}
	if uploadResult.Metadata.PathDisplay != remotePath {
		t.Fatalf("upload metadata path = %q, want %q", uploadResult.Metadata.PathDisplay, remotePath)
	}
	if uploadResult.Metadata.Size != uint64(len(payload)) {
		t.Fatalf("upload metadata size = %d, want %d", uploadResult.Metadata.Size, len(payload))
	}

	localDownloadPath := filepath.Join(t.TempDir(), "download.txt")
	var downloadUpdates []filetransfer.DownloadProgress
	downloadResult, err := filetransfer.NewDownloader(client).Download(
		ctx,
		remotePath,
		filetransfer.File(localDownloadPath),
		filetransfer.DownloadOptions{
			Progress: func(progress filetransfer.DownloadProgress) {
				downloadUpdates = append(downloadUpdates, progress)
			},
		},
	)
	if err != nil {
		t.Fatalf("Download(%q): %v", remotePath, err)
	}
	assertDownloadProgress(t, downloadUpdates, int64(len(payload)))

	data, err := os.ReadFile(localDownloadPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(data); got != payload {
		t.Fatalf("downloaded content = %q, want %q", got, payload)
	}
	if downloadResult.Metadata == nil {
		t.Fatal("Download() metadata is nil")
	}
	if downloadResult.Metadata.PathDisplay != remotePath {
		t.Fatalf("download metadata path = %q, want %q", downloadResult.Metadata.PathDisplay, remotePath)
	}
	if downloadResult.Metadata.Size != uint64(len(payload)) {
		t.Fatalf("download metadata size = %d, want %d", downloadResult.Metadata.Size, len(payload))
	}
}

func TestUserFileTransferDownloadParallelToMemory(t *testing.T) {
	ctx := context.Background()
	client := files.NewContext(userConfig(t))
	remotePath := fmt.Sprintf("/sdk-integration-filetransfer-parallel-%d.txt", time.Now().UnixNano())
	payload := strings.Repeat("parallel filetransfer payload\n", 512)
	uploadIntegrationFile(t, ctx, client, remotePath, payload)
	defer deleteIntegrationPath(t, client, remotePath)

	target := filetransfer.Bytes()
	var (
		progressMu      sync.Mutex
		downloadUpdates []filetransfer.DownloadProgress
	)
	result, err := filetransfer.NewDownloader(client).Download(
		ctx,
		remotePath,
		target,
		filetransfer.DownloadOptions{
			ParallelDownloads: 4,
			Progress: func(progress filetransfer.DownloadProgress) {
				progressMu.Lock()
				downloadUpdates = append(downloadUpdates, progress)
				progressMu.Unlock()
			},
		},
	)
	if err != nil {
		t.Fatalf("Download(%q): %v", remotePath, err)
	}

	if got := string(target.Bytes()); got != payload {
		t.Fatalf("downloaded content = %q, want %q", got, payload)
	}
	if result.Metadata == nil {
		t.Fatal("Download() metadata is nil")
	}
	if result.Metadata.Size != uint64(len(payload)) {
		t.Fatalf("metadata size = %d, want %d", result.Metadata.Size, len(payload))
	}
	progressMu.Lock()
	captured := append([]filetransfer.DownloadProgress(nil), downloadUpdates...)
	progressMu.Unlock()
	assertDownloadProgress(t, captured, int64(len(payload)))
}

func assertUploadProgress(
	t *testing.T,
	updates []filetransfer.UploadProgress,
	total int64,
) {
	t.Helper()
	if total == 0 {
		return
	}
	if len(updates) == 0 {
		t.Fatal("no upload progress updates")
	}
	previous := int64(0)
	for _, update := range updates {
		if update.TotalBytes != total {
			t.Fatalf("upload TotalBytes = %d, want %d", update.TotalBytes, total)
		}
		if update.BytesCommitted <= previous {
			t.Fatalf("upload progress is not increasing: previous=%d current=%d", previous, update.BytesCommitted)
		}
		if update.BytesCommitted > total {
			t.Fatalf("upload progress exceeds total: %d > %d", update.BytesCommitted, total)
		}
		previous = update.BytesCommitted
	}
	if previous != total {
		t.Fatalf("final upload progress = %d, want %d", previous, total)
	}
}

func assertDownloadProgress(
	t *testing.T,
	updates []filetransfer.DownloadProgress,
	total int64,
) {
	t.Helper()
	if total == 0 {
		return
	}
	if len(updates) == 0 {
		t.Fatal("no download progress updates")
	}
	previous := int64(0)
	for _, update := range updates {
		if update.TotalBytes != total {
			t.Fatalf("download TotalBytes = %d, want %d", update.TotalBytes, total)
		}
		if update.BytesCommitted <= previous {
			t.Fatalf("download progress is not increasing: previous=%d current=%d", previous, update.BytesCommitted)
		}
		if update.BytesCommitted > total {
			t.Fatalf("download progress exceeds total: %d > %d", update.BytesCommitted, total)
		}
		previous = update.BytesCommitted
	}
	if previous != total {
		t.Fatalf("final download progress = %d, want %d", previous, total)
	}
}

func uploadIntegrationFile(
	t *testing.T,
	ctx context.Context,
	client files.ContextClient,
	path string,
	payload string,
) {
	t.Helper()

	arg := files.NewUploadArg(path)
	if _, err := client.UploadContext(ctx, arg, strings.NewReader(payload)); err != nil {
		t.Fatalf("UploadContext(%q): %v", path, err)
	}
}

func deleteIntegrationPath(t *testing.T, client files.Client, path string) {
	t.Helper()

	if _, err := client.DeleteV2(files.NewDeleteArg(path)); err != nil {
		t.Errorf("cleanup DeleteV2(%q): %v", path, err)
	}
}
