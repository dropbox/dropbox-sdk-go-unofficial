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

package files_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/contenthash"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type capturedUploadRequest struct {
	path string
	arg  string
	body []byte
}

type nonSeekableReader struct {
	reader *strings.Reader
}

func (r *nonSeekableReader) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

type brokenSeekReader struct {
	reader *strings.Reader
}

func (r *brokenSeekReader) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

func (r *brokenSeekReader) Seek(offset int64, whence int) (int64, error) {
	return 0, errors.New("seek unavailable")
}

type readErrorSeeker struct {
	err error
}

func (r *readErrorSeeker) Read(p []byte) (int, error) {
	return 0, r.err
}

func (r *readErrorSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func TestAutoContentHashUploadRoutes(t *testing.T) {
	const payload = "payload for automatic content hash"
	expectedHash := mustContentHash(t, payload)

	tests := []struct {
		name     string
		wantPath string
		call     func(files.Client, io.Reader) error
	}{
		{
			name:     "alpha upload",
			wantPath: "/files/alpha/upload",
			call: func(client files.Client, content io.Reader) error {
				_, err := client.AlphaUpload(files.NewUploadArg("/alpha.txt"), content)
				return err
			},
		},
		{
			name:     "upload",
			wantPath: "/files/upload",
			call: func(client files.Client, content io.Reader) error {
				_, err := client.Upload(files.NewUploadArg("/upload.txt"), content)
				return err
			},
		},
		{
			name:     "upload session start",
			wantPath: "/files/upload_session/start",
			call: func(client files.Client, content io.Reader) error {
				_, err := client.UploadSessionStart(files.NewUploadSessionStartArg(), content)
				return err
			},
		},
		{
			name:     "upload session append v2",
			wantPath: "/files/upload_session/append_v2",
			call: func(client files.Client, content io.Reader) error {
				err := client.UploadSessionAppendV2(files.NewUploadSessionAppendArg(files.NewUploadSessionCursor("sid", 0)), content)
				return err
			},
		},
		{
			name:     "upload session append batch",
			wantPath: "/files/upload_session/append_batch",
			call: func(client files.Client, content io.Reader) error {
				entry := files.NewUploadSessionAppendBatchArgEntry(files.NewUploadSessionCursor("sid", 0), uint64(len(payload)))
				_, err := client.UploadSessionAppendBatch(files.NewUploadSessionAppendBatchArg([]*files.UploadSessionAppendBatchArgEntry{entry}), content)
				return err
			},
		},
		{
			name:     "upload session finish",
			wantPath: "/files/upload_session/finish",
			call: func(client files.Client, content io.Reader) error {
				arg := files.NewUploadSessionFinishArg(files.NewUploadSessionCursor("sid", 0), files.NewCommitInfo("/finish.txt"))
				_, err := client.UploadSessionFinish(arg, content)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, requests := newUploadTestClient(t)

			if err := tt.call(client, strings.NewReader(payload)); err != nil {
				t.Fatalf("upload route failed: %v", err)
			}

			req := nextUploadRequest(t, requests)
			if req.path != tt.wantPath {
				t.Fatalf("route path = %q, want %q", req.path, tt.wantPath)
			}
			if string(req.body) != payload {
				t.Fatalf("body = %q, want %q", string(req.body), payload)
			}

			arg := decodeDropboxAPIArg(t, req.arg)
			if got := arg["content_hash"]; got != expectedHash {
				t.Fatalf("content_hash = %v, want %q", got, expectedHash)
			}
		})
	}
}

func TestAutoContentHashManualValueWins(t *testing.T) {
	const payload = "payload with manual hash"
	const manualHash = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	client, requests := newUploadTestClient(t)
	arg := files.NewUploadArg("/manual.txt")
	arg.ContentHash = manualHash

	if _, err := client.Upload(arg, strings.NewReader(payload)); err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	req := nextUploadRequest(t, requests)
	gotArg := decodeDropboxAPIArg(t, req.arg)
	if got := gotArg["content_hash"]; got != manualHash {
		t.Fatalf("content_hash = %v, want manual hash %q", got, manualHash)
	}
}

func TestAutoContentHashDoesNotMutateCallerArg(t *testing.T) {
	const payload = "payload for arg mutation check"
	client, requests := newUploadTestClient(t)
	arg := files.NewUploadArg("/mutation.txt")

	if _, err := client.Upload(arg, strings.NewReader(payload)); err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	req := nextUploadRequest(t, requests)
	gotArg := decodeDropboxAPIArg(t, req.arg)
	if got := gotArg["content_hash"]; got != mustContentHash(t, payload) {
		t.Fatalf("content_hash = %v, want computed hash", got)
	}
	if arg.ContentHash != "" {
		t.Fatalf("caller arg ContentHash = %q, want empty", arg.ContentHash)
	}
}

func TestWithoutAutoContentHashOmitsGeneratedHash(t *testing.T) {
	const payload = "payload with opt out"
	client, requests := newUploadTestClient(t)

	if _, err := client.Upload(files.NewUploadArg("/opt-out.txt"), files.WithoutAutoContentHash(strings.NewReader(payload))); err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	req := nextUploadRequest(t, requests)
	if string(req.body) != payload {
		t.Fatalf("body = %q, want %q", string(req.body), payload)
	}
	gotArg := decodeDropboxAPIArg(t, req.arg)
	if _, ok := gotArg["content_hash"]; ok {
		t.Fatalf("content_hash was sent for opt-out reader: %v", gotArg["content_hash"])
	}
}

func TestWithoutAutoContentHashStillSendsManualHash(t *testing.T) {
	const payload = "payload with opt out and manual hash"
	const manualHash = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	client, requests := newUploadTestClient(t)
	arg := files.NewUploadArg("/manual-opt-out.txt")
	arg.ContentHash = manualHash

	if _, err := client.Upload(arg, files.WithoutAutoContentHash(strings.NewReader(payload))); err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	req := nextUploadRequest(t, requests)
	gotArg := decodeDropboxAPIArg(t, req.arg)
	if got := gotArg["content_hash"]; got != manualHash {
		t.Fatalf("content_hash = %v, want manual hash %q", got, manualHash)
	}
}

func TestWithoutAutoContentHashPreservesReaderInterfaces(t *testing.T) {
	content := files.WithoutAutoContentHash(strings.NewReader("payload"))
	if _, ok := content.(io.ReadSeeker); !ok {
		t.Fatal("WithoutAutoContentHash stripped io.ReadSeeker")
	}
	if _, ok := content.(io.WriterTo); !ok {
		t.Fatal("WithoutAutoContentHash stripped io.WriterTo")
	}
}

func TestAutoContentHashSkipsNonSeekableReader(t *testing.T) {
	const payload = "payload from non-seekable reader"
	client, requests := newUploadTestClient(t)
	content := &nonSeekableReader{reader: strings.NewReader(payload)}

	if _, err := client.Upload(files.NewUploadArg("/non-seekable.txt"), content); err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	req := nextUploadRequest(t, requests)
	if string(req.body) != payload {
		t.Fatalf("body = %q, want %q", string(req.body), payload)
	}
	gotArg := decodeDropboxAPIArg(t, req.arg)
	if _, ok := gotArg["content_hash"]; ok {
		t.Fatalf("content_hash was sent for non-seekable reader: %v", gotArg["content_hash"])
	}
}

func TestAutoContentHashSkipsReaderWhenCurrentOffsetCannotBeRead(t *testing.T) {
	const payload = "payload from reader with broken seek"
	client, requests := newUploadTestClient(t)
	content := &brokenSeekReader{reader: strings.NewReader(payload)}

	if _, err := client.Upload(files.NewUploadArg("/broken-seek.txt"), content); err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	req := nextUploadRequest(t, requests)
	if string(req.body) != payload {
		t.Fatalf("body = %q, want %q", string(req.body), payload)
	}
	gotArg := decodeDropboxAPIArg(t, req.arg)
	if _, ok := gotArg["content_hash"]; ok {
		t.Fatalf("content_hash was sent for reader with broken seek: %v", gotArg["content_hash"])
	}
}

func TestAutoContentHashWrapsReadError(t *testing.T) {
	readErr := errors.New("read failed")
	client, requests := newUploadTestClient(t)

	_, err := client.Upload(files.NewUploadArg("/read-error.txt"), &readErrorSeeker{err: readErr})
	if err == nil {
		t.Fatal("upload succeeded, want auto content_hash error")
	}
	if !errors.Is(err, readErr) {
		t.Fatalf("error = %v, want to wrap %v", err, readErr)
	}
	if !strings.Contains(err.Error(), "auto content_hash") {
		t.Fatalf("error = %q, want auto content_hash context", err.Error())
	}

	select {
	case req := <-requests:
		t.Fatalf("unexpected request after auto content_hash failure: %s", req.path)
	default:
	}
}

func TestAutoContentHashUsesAndPreservesCurrentReaderOffset(t *testing.T) {
	const prefix = "already-read:"
	const payload = "payload from current offset"
	client, requests := newUploadTestClient(t)
	content := strings.NewReader(prefix + payload)
	if _, err := content.Seek(int64(len(prefix)), io.SeekStart); err != nil {
		t.Fatalf("seek failed: %v", err)
	}

	if _, err := client.Upload(files.NewUploadArg("/offset.txt"), content); err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	req := nextUploadRequest(t, requests)
	if string(req.body) != payload {
		t.Fatalf("body = %q, want %q", string(req.body), payload)
	}
	gotArg := decodeDropboxAPIArg(t, req.arg)
	if got := gotArg["content_hash"]; got != mustContentHash(t, payload) {
		t.Fatalf("content_hash = %v, want hash of %q", got, payload)
	}
}

func newUploadTestClient(t *testing.T) (files.Client, <-chan capturedUploadRequest) {
	t.Helper()

	requests := make(chan capturedUploadRequest, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		requests <- capturedUploadRequest{
			path: r.URL.Path,
			arg:  r.Header.Get("Dropbox-API-Arg"),
			body: body,
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, err = w.Write([]byte(uploadTestResponse(r.URL.Path)))
		if err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	t.Cleanup(ts.Close)

	config := dropbox.Config{
		Client: ts.Client(),
		URLGenerator: func(hostType string, namespace string, route string) string {
			return fmt.Sprintf("%s/%s/%s", ts.URL, namespace, route)
		},
	}
	return files.New(config), requests
}

func uploadTestResponse(path string) string {
	switch path {
	case "/files/upload_session/start":
		return `{"session_id":"sid"}`
	case "/files/upload_session/append_v2":
		return `{}`
	case "/files/upload_session/append_batch":
		return `{"entries":[{".tag":"success"}]}`
	default:
		return `{"name":"file.txt","id":"id:test","client_modified":"2020-01-02T03:04:05Z","server_modified":"2020-01-02T03:04:06Z","rev":"rev","size":1,"is_downloadable":true}`
	}
}

func nextUploadRequest(t *testing.T, requests <-chan capturedUploadRequest) capturedUploadRequest {
	t.Helper()

	select {
	case req := <-requests:
		return req
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for upload request")
		return capturedUploadRequest{}
	}
}

func decodeDropboxAPIArg(t *testing.T, header string) map[string]interface{} {
	t.Helper()

	var arg map[string]interface{}
	if err := json.Unmarshal([]byte(header), &arg); err != nil {
		t.Fatalf("unmarshal Dropbox-API-Arg %q: %v", header, err)
	}
	return arg
}

func mustContentHash(t *testing.T, payload string) string {
	t.Helper()

	hash, err := contenthash.Compute(strings.NewReader(payload))
	if err != nil {
		t.Fatalf("compute content hash: %v", err)
	}
	return hash
}
