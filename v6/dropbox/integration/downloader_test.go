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
	"testing"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/filedownload"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

func TestUserDownloadFile(t *testing.T) {
	ctx := context.Background()
	client := files.NewContext(userConfig(t))
	remotePath := fmt.Sprintf("/sdk-integration-download-%d.txt", time.Now().UnixNano())
	payload := "hello from the filedownload integration test\n"
	uploadIntegrationFile(t, ctx, client, remotePath, payload)
	defer deleteIntegrationPath(t, client, remotePath)

	localPath := filepath.Join(t.TempDir(), "download.txt")
	result, err := filedownload.New(client).DownloadFile(ctx, remotePath, localPath)
	if err != nil {
		t.Fatalf("DownloadFile(%q): %v", remotePath, err)
	}

	data, err := os.ReadFile(localPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(data); got != payload {
		t.Fatalf("downloaded content = %q, want %q", got, payload)
	}
	if result.Metadata == nil {
		t.Fatal("DownloadFile() metadata is nil")
	}
	if result.Metadata.PathDisplay != remotePath {
		t.Fatalf("metadata path = %q, want %q", result.Metadata.PathDisplay, remotePath)
	}
	if result.Metadata.Size != uint64(len(payload)) {
		t.Fatalf("metadata size = %d, want %d", result.Metadata.Size, len(payload))
	}
}

func TestUserDownloadFileResumesFromPartFile(t *testing.T) {
	ctx := context.Background()
	client := files.NewContext(userConfig(t))
	remotePath := fmt.Sprintf("/sdk-integration-download-resume-%d.txt", time.Now().UnixNano())
	prefix := "already downloaded "
	suffix := "and fetched with a range request\n"
	payload := prefix + suffix
	uploadIntegrationFile(t, ctx, client, remotePath, payload)
	defer deleteIntegrationPath(t, client, remotePath)

	localPath := filepath.Join(t.TempDir(), "download.txt")
	if err := os.WriteFile(localPath+".part", []byte(prefix), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := filedownload.New(client).DownloadFile(ctx, remotePath, localPath)
	if err != nil {
		t.Fatalf("DownloadFile(%q): %v", remotePath, err)
	}

	data, err := os.ReadFile(localPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(data); got != payload {
		t.Fatalf("downloaded content = %q, want %q", got, payload)
	}
	if result.ResumedFrom != int64(len(prefix)) {
		t.Fatalf("ResumedFrom = %d, want %d", result.ResumedFrom, len(prefix))
	}
	if _, err := os.Stat(localPath + ".part"); !os.IsNotExist(err) {
		t.Fatalf("part file still exists: %v", err)
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
