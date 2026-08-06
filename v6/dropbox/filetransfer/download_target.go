package filetransfer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/contenthash"
)

// DownloadInfo describes the content being downloaded.
type DownloadInfo struct {
	// Size is the expected size of the download.
	Size int64

	// ContentHash is the expected Dropbox content hash.
	//
	// An empty value means that no content hash is available.
	ContentHash string
}

// DownloadTarget receives downloaded content.
type DownloadTarget interface {
	// Prepare initializes the target for the download.
	Prepare(ctx context.Context, info DownloadInfo) error

	// WriterAt writes data at the specified offset.
	//
	// Calls may be concurrent and may occur out of order when
	// ParallelDownloads is greater than one.
	io.WriterAt

	// Commit validates and finalizes a successful download.
	Commit(ctx context.Context) error

	// Abort discards any incomplete state.
	Abort(ctx context.Context, cause error) error
}

// -----------------------------------------------------------------------------
// Built-in download targets
// -----------------------------------------------------------------------------

// BytesTarget stores downloaded content in memory.
type BytesTarget struct {
	mu        sync.RWMutex
	info      DownloadInfo
	data      []byte
	prepared  bool
	committed bool
}

// Bytes returns a new in-memory download target.
func Bytes() *BytesTarget {
	return &BytesTarget{}
}

// Prepare initializes the target for the download.
func (t *BytesTarget) Prepare(
	ctx context.Context,
	info DownloadInfo,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if info.Size < 0 {
		return errors.New("download size must not be negative")
	}
	if info.Size > int64(maxInt()) {
		return fmt.Errorf(
			"download is too large for an in-memory target: %d bytes",
			info.Size,
		)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.prepared {
		return errors.New("download target is already prepared")
	}

	t.info = info
	t.data = make([]byte, int(info.Size))
	t.prepared = true
	t.committed = false

	return nil
}

// WriteAt writes data at the specified offset.
func (t *BytesTarget) WriteAt(data []byte, offset int64) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.prepared {
		return 0, errors.New("download target is not prepared")
	}
	if offset < 0 || offset > int64(len(t.data)) {
		return 0, fmt.Errorf("invalid write offset: %d", offset)
	}

	start := int(offset)
	if len(data) > len(t.data)-start {
		return 0, io.ErrShortWrite
	}

	return copy(t.data[start:], data), nil
}

// Commit validates and finalizes a successful download.
func (t *BytesTarget) Commit(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.prepared {
		return errors.New("download target is not prepared")
	}

	if t.info.ContentHash != "" {
		actual, err := contenthash.Compute(bytes.NewReader(t.data))
		if err != nil {
			return fmt.Errorf("compute download content hash: %w", err)
		}
		if actual != t.info.ContentHash {
			return fmt.Errorf(
				"download content hash mismatch: got %q, expected %q",
				actual,
				t.info.ContentHash,
			)
		}
	}

	t.prepared = false
	t.committed = true

	return nil
}

// Abort discards any incomplete state.
func (t *BytesTarget) Abort(
	_ context.Context,
	_ error,
) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.info = DownloadInfo{}
	t.data = nil
	t.prepared = false
	t.committed = false

	return nil
}

// Bytes returns the downloaded content.
//
// Bytes returns nil before a successful commit. The returned slice is owned by
// the target and must not be modified.
func (t *BytesTarget) Bytes() []byte {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if !t.committed {
		return nil
	}

	return t.data
}

// fileTarget stores downloaded content in a temporary file.
//
// Commit validates the completed download and renames the temporary file to
// its final destination. Abort removes the temporary file.
type fileTarget struct {
	mu       sync.Mutex
	info     DownloadInfo
	path     string
	tempPath string
	file     *os.File
}

// File returns a target that writes downloaded content to path.
func File(path string) DownloadTarget {
	return &fileTarget{
		path: path,
	}
}

// Prepare initializes the target for the download.
func (t *fileTarget) Prepare(
	ctx context.Context,
	info DownloadInfo,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if info.Size < 0 {
		return errors.New("download size must not be negative")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.file != nil || t.tempPath != "" {
		return errors.New("download target is already prepared")
	}

	dir := filepath.Dir(t.path)
	base := filepath.Base(t.path)

	file, err := os.CreateTemp(dir, "."+base+".*.part")
	if err != nil {
		return err
	}

	tempPath := file.Name()

	if err := file.Truncate(info.Size); err != nil {
		_ = file.Close()
		_ = os.Remove(tempPath)
		return err
	}

	t.info = info
	t.file = file
	t.tempPath = tempPath

	return nil
}

// WriteAt writes data at the specified offset.
func (t *fileTarget) WriteAt(data []byte, offset int64) (int, error) {
	t.mu.Lock()
	file := t.file
	t.mu.Unlock()

	if file == nil {
		return 0, errors.New("download target is not prepared")
	}

	return file.WriteAt(data, offset)
}

// Commit validates and finalizes a successful download.
func (t *fileTarget) Commit(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.file == nil || t.tempPath == "" {
		return errors.New("download target is not prepared")
	}

	stat, err := t.file.Stat()
	if err != nil {
		return fmt.Errorf("stat downloaded content: %w", err)
	}
	if stat.Size() != t.info.Size {
		return fmt.Errorf(
			"download size mismatch: got %d bytes, expected %d",
			stat.Size(),
			t.info.Size,
		)
	}

	if t.info.ContentHash != "" {
		if _, err := t.file.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("seek downloaded content: %w", err)
		}

		actual, err := contenthash.Compute(t.file)
		if err != nil {
			return fmt.Errorf("compute download content hash: %w", err)
		}
		if actual != t.info.ContentHash {
			return fmt.Errorf(
				"download content hash mismatch: got %q, expected %q",
				actual,
				t.info.ContentHash,
			)
		}
	}

	if err := t.file.Close(); err != nil {
		return err
	}
	t.file = nil

	if err := os.Rename(t.tempPath, t.path); err != nil {
		return err
	}

	t.info = DownloadInfo{}
	t.tempPath = ""

	return nil
}

// Abort discards any incomplete state.
func (t *fileTarget) Abort(
	_ context.Context,
	_ error,
) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	var closeErr error
	if t.file != nil {
		closeErr = t.file.Close()
		t.file = nil
	}

	var removeErr error
	if t.tempPath != "" {
		removeErr = os.Remove(t.tempPath)
		if errors.Is(removeErr, os.ErrNotExist) {
			removeErr = nil
		}
	}

	t.info = DownloadInfo{}
	t.tempPath = ""

	return errors.Join(closeErr, removeErr)
}

func maxInt() int {
	return int(^uint(0) >> 1)
}
