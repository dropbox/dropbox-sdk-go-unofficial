package filetransfer

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileUpload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "upload.bin")
	data := []byte("hello")

	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	source, err := FileUpload(path)
	if err != nil {
		t.Fatalf("FileUpload() error = %v", err)
	}

	if got := source.Size(); got != int64(len(data)) {
		t.Fatalf("Size() = %d, want %d", got, len(data))
	}
}

func TestFileUploadNotRegularFile(t *testing.T) {
	_, err := FileUpload(t.TempDir())
	if err == nil {
		t.Fatal("FileUpload() error = nil, want error")
	}

	if !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("FileUpload() error = %q, want not a regular file", err)
	}
}

func TestFileSourceOpen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "upload.bin")
	data := []byte("hello")

	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	source, err := FileUpload(path)
	if err != nil {
		t.Fatalf("FileUpload() error = %v", err)
	}

	reader, err := source.Open(context.Background())
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer reader.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if !bytes.Equal(got, data) {
		t.Fatalf("Open() data = %q, want %q", got, data)
	}
}

func TestFileSourceOpenRange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "upload.bin")

	if err := os.WriteFile(path, []byte("hello"), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	source, err := FileUpload(path)
	if err != nil {
		t.Fatalf("FileUpload() error = %v", err)
	}

	reader, err := source.OpenRange(context.Background(), 1, 3)
	if err != nil {
		t.Fatalf("OpenRange() error = %v", err)
	}
	defer reader.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if string(got) != "ell" {
		t.Fatalf("OpenRange() data = %q, want %q", got, "ell")
	}
}

func TestFileSourceOpenRangeOutOfBounds(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "upload.bin")

	if err := os.WriteFile(path, []byte("hello"), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	source, err := FileUpload(path)
	if err != nil {
		t.Fatalf("FileUpload() error = %v", err)
	}

	tests := []struct {
		name   string
		offset int64
		length int64
	}{
		{
			name:   "negative offset",
			offset: -1,
			length: 1,
		},
		{
			name:   "negative length",
			offset: 0,
			length: -1,
		},
		{
			name:   "offset past end",
			offset: 6,
			length: 0,
		},
		{
			name:   "range exceeds end",
			offset: 4,
			length: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, err := source.OpenRange(
				context.Background(),
				tt.offset,
				tt.length,
			)
			if err == nil {
				if reader != nil {
					_ = reader.Close()
				}
				t.Fatal("OpenRange() error = nil, want error")
			}
		})
	}
}

func TestBytesUpload(t *testing.T) {
	source := BytesUpload([]byte("hello"))

	if got := source.Size(); got != 5 {
		t.Fatalf("Size() = %d, want 5", got)
	}
}

func TestBytesSourceOpen(t *testing.T) {
	source := BytesUpload([]byte("hello"))

	reader, err := source.Open(context.Background())
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer reader.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if string(got) != "hello" {
		t.Fatalf("Open() data = %q, want %q", got, "hello")
	}
}

func TestBytesSourceOpenRange(t *testing.T) {
	source := BytesUpload([]byte("hello"))

	reader, err := source.OpenRange(context.Background(), 1, 3)
	if err != nil {
		t.Fatalf("OpenRange() error = %v", err)
	}
	defer reader.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if string(got) != "ell" {
		t.Fatalf("OpenRange() data = %q, want %q", got, "ell")
	}
}

func TestBytesSourceOpenRangeOutOfBounds(t *testing.T) {
	source := BytesUpload([]byte("hello"))

	tests := []struct {
		name   string
		offset int64
		length int64
	}{
		{
			name:   "negative offset",
			offset: -1,
			length: 1,
		},
		{
			name:   "negative length",
			offset: 0,
			length: -1,
		},
		{
			name:   "offset past end",
			offset: 6,
			length: 0,
		},
		{
			name:   "range exceeds end",
			offset: 4,
			length: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, err := source.OpenRange(
				context.Background(),
				tt.offset,
				tt.length,
			)
			if err == nil {
				if reader != nil {
					_ = reader.Close()
				}
				t.Fatal("OpenRange() error = nil, want error")
			}
		})
	}
}

func TestReaderUpload(t *testing.T) {
	source, err := ReaderUpload(bytes.NewReader([]byte("hello")))
	if err != nil {
		t.Fatalf("ReaderUpload() error = %v", err)
	}

	reader, err := source.Open(context.Background())
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer reader.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if string(got) != "hello" {
		t.Fatalf("Open() data = %q, want %q", got, "hello")
	}
}

func TestReaderUploadNil(t *testing.T) {
	_, err := ReaderUpload(nil)
	if err == nil {
		t.Fatal("ReaderUpload() error = nil, want error")
	}
}

func TestReaderSourceOpenOnce(t *testing.T) {
	source, err := ReaderUpload(bytes.NewReader([]byte("hello")))
	if err != nil {
		t.Fatalf("ReaderUpload() error = %v", err)
	}

	reader, err := source.Open(context.Background())
	if err != nil {
		t.Fatalf("first Open() error = %v", err)
	}
	_ = reader.Close()

	if _, err := source.Open(context.Background()); err == nil {
		t.Fatal("second Open() error = nil, want error")
	}
}

func TestReaderSourcePreservesReadCloser(t *testing.T) {
	reader := &trackingReadCloser{
		Reader: bytes.NewReader([]byte("hello")),
	}

	source, err := ReaderUpload(reader)
	if err != nil {
		t.Fatalf("ReaderUpload() error = %v", err)
	}

	opened, err := source.Open(context.Background())
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	if err := opened.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	if !reader.closed {
		t.Fatal("underlying reader was not closed")
	}
}

func TestSizedReaderUpload(t *testing.T) {
	source, err := SizedReaderUpload(
		bytes.NewReader([]byte("hello")),
		5,
	)
	if err != nil {
		t.Fatalf("SizedReaderUpload() error = %v", err)
	}

	if got := source.Size(); got != 5 {
		t.Fatalf("Size() = %d, want 5", got)
	}

	reader, err := source.Open(context.Background())
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer reader.Close()

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if string(got) != "hello" {
		t.Fatalf("Open() data = %q, want %q", got, "hello")
	}
}

func TestSizedReaderUploadNegativeSize(t *testing.T) {
	_, err := SizedReaderUpload(
		bytes.NewReader(nil),
		-1,
	)
	if err == nil {
		t.Fatal("SizedReaderUpload() error = nil, want error")
	}
}

func TestUploadSourceCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	filePath := filepath.Join(t.TempDir(), "upload.bin")
	if err := os.WriteFile(filePath, []byte("hello"), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	fileSource, err := FileUpload(filePath)
	if err != nil {
		t.Fatalf("FileUpload() error = %v", err)
	}

	bytesSource := BytesUpload([]byte("hello"))

	readerSource, err := ReaderUpload(bytes.NewReader([]byte("hello")))
	if err != nil {
		t.Fatalf("ReaderUpload() error = %v", err)
	}

	tests := []struct {
		name string
		open func() (io.ReadCloser, error)
	}{
		{
			name: "file",
			open: func() (io.ReadCloser, error) {
				return fileSource.Open(ctx)
			},
		},
		{
			name: "bytes",
			open: func() (io.ReadCloser, error) {
				return bytesSource.Open(ctx)
			},
		},
		{
			name: "reader",
			open: func() (io.ReadCloser, error) {
				return readerSource.Open(ctx)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, err := tt.open()
			if reader != nil {
				_ = reader.Close()
			}
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("Open() error = %v, want context.Canceled", err)
			}
		})
	}
}

func TestSizedReaderUploadNil(t *testing.T) {
	_, err := SizedReaderUpload(nil, 0)
	if err == nil {
		t.Fatal("SizedReaderUpload() error = nil, want error")
	}
}

type trackingReadCloser struct {
	io.Reader
	closed bool
}

func (r *trackingReadCloser) Close() error {
	r.closed = true
	return nil
}
