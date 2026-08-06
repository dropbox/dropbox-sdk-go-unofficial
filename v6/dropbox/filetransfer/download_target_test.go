package filetransfer

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/contenthash"
)

func TestBytesTargetPrepare(t *testing.T) {
	target := Bytes()

	if err := target.Prepare(context.Background(), DownloadInfo{
		Size: 5,
	}); err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}

	if got := target.Bytes(); got != nil {
		t.Fatalf("Bytes() = %q, want nil before commit", got)
	}

	n, err := target.WriteAt([]byte("hello"), 0)
	if err != nil {
		t.Fatalf("WriteAt() error = %v", err)
	}
	if n != 5 {
		t.Fatalf("WriteAt() n = %d, want 5", n)
	}
}

func TestBytesTargetWriteAt(t *testing.T) {
	target := Bytes()

	if err := target.Prepare(context.Background(), DownloadInfo{
		Size: 5,
	}); err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}

	if _, err := target.WriteAt([]byte("he"), 0); err != nil {
		t.Fatalf("WriteAt() error = %v", err)
	}
	if _, err := target.WriteAt([]byte("llo"), 2); err != nil {
		t.Fatalf("WriteAt() error = %v", err)
	}

	if err := target.Commit(context.Background()); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}

	if got := string(target.Bytes()); got != "hello" {
		t.Fatalf("Bytes() = %q, want %q", got, "hello")
	}
}

func TestBytesTargetWriteAtOutOfRange(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		offset int64
	}{
		{
			name:   "negative offset",
			data:   []byte("x"),
			offset: -1,
		},
		{
			name:   "offset past end",
			data:   []byte("x"),
			offset: 6,
		},
		{
			name:   "write exceeds end",
			data:   []byte("hello"),
			offset: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := Bytes()

			if err := target.Prepare(context.Background(), DownloadInfo{
				Size: 5,
			}); err != nil {
				t.Fatalf("Prepare() error = %v", err)
			}

			if _, err := target.WriteAt(tt.data, tt.offset); err == nil {
				t.Fatal("WriteAt() error = nil, want error")
			}
		})
	}
}

func TestBytesTargetCommit(t *testing.T) {
	data := []byte("hello")

	hash, err := contenthash.Compute(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("contenthash.Compute() error = %v", err)
	}

	target := Bytes()

	if err := target.Prepare(context.Background(), DownloadInfo{
		Size:        int64(len(data)),
		ContentHash: hash,
	}); err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}

	if _, err := target.WriteAt(data, 0); err != nil {
		t.Fatalf("WriteAt() error = %v", err)
	}

	if err := target.Commit(context.Background()); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}

	if got := target.Bytes(); !bytes.Equal(got, data) {
		t.Fatalf("Bytes() = %q, want %q", got, data)
	}
}

func TestBytesTargetCommitHashMismatch(t *testing.T) {
	target := Bytes()

	if err := target.Prepare(context.Background(), DownloadInfo{
		Size:        5,
		ContentHash: strings.Repeat("0", 64),
	}); err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}

	if _, err := target.WriteAt([]byte("hello"), 0); err != nil {
		t.Fatalf("WriteAt() error = %v", err)
	}

	err := target.Commit(context.Background())
	if err == nil {
		t.Fatal("Commit() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "content hash mismatch") {
		t.Fatalf("Commit() error = %q, want content hash mismatch", err)
	}

	if got := target.Bytes(); got != nil {
		t.Fatalf("Bytes() = %q, want nil", got)
	}
}

func TestBytesTargetAbort(t *testing.T) {
	target := Bytes()

	if err := target.Prepare(context.Background(), DownloadInfo{
		Size: 5,
	}); err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}

	if _, err := target.WriteAt([]byte("hello"), 0); err != nil {
		t.Fatalf("WriteAt() error = %v", err)
	}

	if err := target.Abort(context.Background(), errors.New("failed")); err != nil {
		t.Fatalf("Abort() error = %v", err)
	}

	if got := target.Bytes(); got != nil {
		t.Fatalf("Bytes() = %q, want nil", got)
	}

	if _, err := target.WriteAt([]byte("x"), 0); err == nil {
		t.Fatal("WriteAt() after Abort() error = nil, want error")
	}

	if err := target.Commit(context.Background()); err == nil {
		t.Fatal("Commit() after Abort() error = nil, want error")
	}
}

func TestBytesTargetBytesBeforeCommit(t *testing.T) {
	target := Bytes()

	if err := target.Prepare(context.Background(), DownloadInfo{
		Size: 5,
	}); err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}

	if _, err := target.WriteAt([]byte("hello"), 0); err != nil {
		t.Fatalf("WriteAt() error = %v", err)
	}

	if got := target.Bytes(); got != nil {
		t.Fatalf("Bytes() = %q, want nil", got)
	}
}

func TestBytesTargetDoublePrepare(t *testing.T) {
	target := Bytes()

	if err := target.Prepare(context.Background(), DownloadInfo{
		Size: 5,
	}); err != nil {
		t.Fatalf("first Prepare() error = %v", err)
	}

	err := target.Prepare(context.Background(), DownloadInfo{
		Size: 10,
	})
	if err == nil {
		t.Fatal("second Prepare() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "already prepared") {
		t.Fatalf("second Prepare() error = %q, want already prepared", err)
	}
}

func TestFileTargetPrepare(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "download.bin")
	target := File(path)

	if err := target.Prepare(context.Background(), DownloadInfo{
		Size: 5,
	}); err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}
	t.Cleanup(func() {
		_ = target.Abort(context.Background(), nil)
	})

	if _, err := target.WriteAt([]byte("hello"), 0); err != nil {
		t.Fatalf("WriteAt() error = %v", err)
	}
}

func TestFileTargetWriteAt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "download.bin")
	target := File(path)

	if err := target.Prepare(context.Background(), DownloadInfo{
		Size: 5,
	}); err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}

	if _, err := target.WriteAt([]byte("he"), 0); err != nil {
		t.Fatalf("WriteAt() error = %v", err)
	}
	if _, err := target.WriteAt([]byte("llo"), 2); err != nil {
		t.Fatalf("WriteAt() error = %v", err)
	}

	if err := target.Commit(context.Background()); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(got) != "hello" {
		t.Fatalf("file contents = %q, want %q", got, "hello")
	}
}

func TestFileTargetCommit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "download.bin")
	data := []byte("hello")

	hash, err := contenthash.Compute(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("contenthash.Compute() error = %v", err)
	}

	target := File(path)

	if err := target.Prepare(context.Background(), DownloadInfo{
		Size:        int64(len(data)),
		ContentHash: hash,
	}); err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}

	if _, err := target.WriteAt(data, 0); err != nil {
		t.Fatalf("WriteAt() error = %v", err)
	}

	if err := target.Commit(context.Background()); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Fatalf("file contents = %q, want %q", got, data)
	}
}

func TestFileTargetCommitHashMismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "download.bin")
	target := File(path)

	if err := target.Prepare(context.Background(), DownloadInfo{
		Size:        5,
		ContentHash: strings.Repeat("0", 64),
	}); err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}

	if _, err := target.WriteAt([]byte("hello"), 0); err != nil {
		t.Fatalf("WriteAt() error = %v", err)
	}

	err := target.Commit(context.Background())
	if err == nil {
		t.Fatal("Commit() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "content hash mismatch") {
		t.Fatalf("Commit() error = %q, want content hash mismatch", err)
	}

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("destination exists after failed commit, Stat() error = %v", err)
	}

	if err := target.Abort(context.Background(), err); err != nil {
		t.Fatalf("Abort() error = %v", err)
	}
}

func TestFileTargetAbort(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "download.bin")
	target := File(path)

	if err := target.Prepare(context.Background(), DownloadInfo{
		Size: 5,
	}); err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}

	if _, err := target.WriteAt([]byte("hello"), 0); err != nil {
		t.Fatalf("WriteAt() error = %v", err)
	}

	if err := target.Abort(
		context.Background(),
		errors.New("download failed"),
	); err != nil {
		t.Fatalf("Abort() error = %v", err)
	}

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("destination exists after Abort(), Stat() error = %v", err)
	}

	if _, err := target.WriteAt([]byte("x"), 0); err == nil {
		t.Fatal("WriteAt() after Abort() error = nil, want error")
	}

	if err := target.Commit(context.Background()); err == nil {
		t.Fatal("Commit() after Abort() error = nil, want error")
	}
}

func TestFileTargetDoublePrepare(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "download.bin")
	target := File(path)

	if err := target.Prepare(context.Background(), DownloadInfo{
		Size: 5,
	}); err != nil {
		t.Fatalf("first Prepare() error = %v", err)
	}

	err := target.Prepare(context.Background(), DownloadInfo{
		Size: 10,
	})
	if err == nil {
		t.Fatal("second Prepare() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "already prepared") {
		t.Fatalf("second Prepare() error = %q, want already prepared", err)
	}

	if err := target.Abort(context.Background(), err); err != nil {
		t.Fatalf("Abort() error = %v", err)
	}
}

func TestFileTargetWriteBeforePrepare(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "download.bin")
	target := File(path)

	if _, err := target.WriteAt([]byte("hello"), 0); err == nil {
		t.Fatal("WriteAt() error = nil, want error")
	}
}

func TestFileTargetCommitSizeMismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "download.bin")

	target := File(path).(*fileTarget)

	if err := target.Prepare(context.Background(), DownloadInfo{
		Size: 5,
	}); err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}

	// Simulate an incomplete/corrupted target state.
	if err := target.file.Truncate(3); err != nil {
		t.Fatalf("Truncate() error = %v", err)
	}

	err := target.Commit(context.Background())
	if err == nil {
		t.Fatal("Commit() error = nil, want size mismatch")
	}
	if !strings.Contains(err.Error(), "download size mismatch") {
		t.Fatalf("Commit() error = %q, want size mismatch", err)
	}

	if err := target.Abort(context.Background(), err); err != nil {
		t.Fatalf("Abort() error = %v", err)
	}
}

func TestFileTargetCommitWithContentHash(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "download.bin")
	data := []byte("hello")

	hash, err := contenthash.Compute(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("contenthash.Compute() error = %v", err)
	}

	target := File(path)

	if err := target.Prepare(context.Background(), DownloadInfo{
		Size:        int64(len(data)),
		ContentHash: hash,
	}); err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}

	if _, err := target.WriteAt(data, 0); err != nil {
		t.Fatalf("WriteAt() error = %v", err)
	}

	if err := target.Commit(context.Background()); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Fatalf("file contents = %q, want %q", got, data)
	}
}
