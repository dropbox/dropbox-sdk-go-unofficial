package filetransfer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/contenthash"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

func TestDownloadSequential(t *testing.T) {
	data := []byte("hello reliable download")
	client := &fakeDownloadClient{data: data, metadata: fileMetadata(data)}
	target := Bytes()

	var updates []DownloadProgress
	result, err := NewDownloader(client).Download(
		context.Background(),
		"/input.txt",
		target,
		DownloadOptions{
			Progress: func(progress DownloadProgress) {
				updates = append(updates, progress)
			},
		},
	)
	if err != nil {
		t.Fatalf("Download() error = %v", err)
	}
	if result == nil || result.Metadata != client.metadata {
		t.Fatalf("Download() metadata = %#v, want %#v", result, client.metadata)
	}
	if got := target.Bytes(); !bytes.Equal(got, data) {
		t.Fatalf("downloaded bytes = %q, want %q", got, data)
	}
	assertDownloadProgress(t, updates, int64(len(data)))
}

func TestDownloadSequentialRetriesFromCommittedOffset(t *testing.T) {
	data := []byte("abcdefghij")
	client := &fakeDownloadClient{
		data:         data,
		metadata:     fileMetadata(data),
		failFirstAt:  4,
		failFirstErr: errors.New("connection reset"),
	}
	target := Bytes()

	var updates []DownloadProgress
	_, err := NewDownloader(client).Download(
		context.Background(),
		"/retry.bin",
		target,
		DownloadOptions{
			MaxAttempts: 3,
			Progress: func(progress DownloadProgress) {
				updates = append(updates, progress)
			},
		},
	)
	if err != nil {
		t.Fatalf("Download() error = %v", err)
	}
	if got := target.Bytes(); !bytes.Equal(got, data) {
		t.Fatalf("downloaded bytes = %q, want %q", got, data)
	}

	client.mu.Lock()
	ranges := append([]string(nil), client.ranges...)
	client.mu.Unlock()
	if len(ranges) != 2 {
		t.Fatalf("download requests = %v, want 2 requests", ranges)
	}
	if ranges[0] != "" || ranges[1] != "bytes=4-" {
		t.Fatalf("download ranges = %v, want [\"\" \"bytes=4-\"]", ranges)
	}
	assertDownloadProgress(t, updates, int64(len(data)))
}

func TestDownloadDoesNotRetryPermanentAPIError(t *testing.T) {
	client := &fakeDownloadClient{
		data:     []byte("unused"),
		metadata: fileMetadata([]byte("unused")),
		apiErr:   downloadPathAPIError(),
	}

	_, err := NewDownloader(client).Download(
		context.Background(),
		"/missing",
		Bytes(),
		DownloadOptions{MaxAttempts: 3},
	)
	if err == nil {
		t.Fatal("Download() error = nil, want API error")
	}
	if client.calls != 1 {
		t.Fatalf("download calls = %d, want 1", client.calls)
	}
}

func TestDownloadParallel(t *testing.T) {
	data := bytes.Repeat([]byte("parallel-download-"), 128)
	client := &fakeDownloadClient{data: data, metadata: fileMetadata(data)}
	target := Bytes()

	var (
		progressMu sync.Mutex
		updates    []DownloadProgress
	)
	_, err := NewDownloader(client).Download(
		context.Background(),
		"/parallel.bin",
		target,
		DownloadOptions{
			ParallelDownloads: 4,
			Progress: func(progress DownloadProgress) {
				progressMu.Lock()
				updates = append(updates, progress)
				progressMu.Unlock()
			},
		},
	)
	if err != nil {
		t.Fatalf("Download() error = %v", err)
	}
	if got := target.Bytes(); !bytes.Equal(got, data) {
		t.Fatalf("downloaded bytes differ: got %d bytes, want %d", len(got), len(data))
	}

	progressMu.Lock()
	captured := append([]DownloadProgress(nil), updates...)
	progressMu.Unlock()
	assertDownloadProgress(t, captured, int64(len(data)))

	client.mu.Lock()
	ranges := append([]string(nil), client.ranges...)
	client.mu.Unlock()
	if len(ranges) != 5 {
		t.Fatalf("download requests = %d, want initial request plus 4 ranges; ranges=%v", len(ranges), ranges)
	}
	if ranges[0] != "bytes=0-0" {
		t.Fatalf("initial range = %q, want bytes=0-0", ranges[0])
	}
}

func TestDownloadParallelZeroByteFallsBackToSequential(t *testing.T) {
	client := &fakeDownloadClient{
		data:     nil,
		metadata: fileMetadata(nil),
	}
	target := Bytes()

	_, err := NewDownloader(client).Download(
		context.Background(),
		"/empty.bin",
		target,
		DownloadOptions{ParallelDownloads: 4},
	)
	if err != nil {
		t.Fatalf("Download() error = %v", err)
	}
	if got := target.Bytes(); len(got) != 0 {
		t.Fatalf("downloaded bytes = %q, want empty", got)
	}

	client.mu.Lock()
	ranges := append([]string(nil), client.ranges...)
	client.mu.Unlock()
	if len(ranges) != 2 {
		t.Fatalf("download requests = %v, want initial range plus sequential fallback", ranges)
	}
	if ranges[0] != "bytes=0-0" || ranges[1] != "" {
		t.Fatalf("download ranges = %v, want [bytes=0-0 \"\"]", ranges)
	}
}

func TestDownloadCommitFailureCallsAbort(t *testing.T) {
	data := []byte("data")
	target := &recordingDownloadTarget{commitErr: errors.New("commit failed")}
	client := &fakeDownloadClient{data: data, metadata: fileMetadata(data)}

	_, err := NewDownloader(client).Download(
		context.Background(),
		"/input",
		target,
		DownloadOptions{},
	)
	if err == nil || !strings.Contains(err.Error(), "commit failed") {
		t.Fatalf("Download() error = %v, want commit failure", err)
	}
	if target.commitCalls != 1 {
		t.Fatalf("Commit calls = %d, want 1", target.commitCalls)
	}
	if target.abortCalls != 1 {
		t.Fatalf("Abort calls = %d, want 1", target.abortCalls)
	}
}

func TestDownloadAbortUsesCleanupContext(t *testing.T) {
	tests := []struct {
		name    string
		options DownloadOptions
	}{
		{
			name: "sequential",
		},
		{
			name: "parallel",
			options: DownloadOptions{
				ParallelDownloads: 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			target := &recordingDownloadTarget{
				commitErr:  context.Canceled,
				commitHook: cancel,
			}
			data := []byte("data for canceled cleanup")
			client := &fakeDownloadClient{data: data, metadata: fileMetadata(data)}

			_, err := NewDownloader(client).Download(
				ctx,
				"/input",
				target,
				tt.options,
			)
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("Download() error = %v, want context.Canceled", err)
			}
			if target.abortCalls != 1 {
				t.Fatalf("Abort calls = %d, want 1", target.abortCalls)
			}
			if target.abortCtxErr != nil {
				t.Fatalf("Abort context error = %v, want nil", target.abortCtxErr)
			}
		})
	}
}

func TestDownloadFileHashMismatchCleansTemporaryFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.bin")
	data := []byte("corrupt")
	metadata := fileMetadata(data)
	metadata.ContentHash = strings.Repeat("0", 64)
	client := &fakeDownloadClient{data: data, metadata: metadata}

	_, err := NewDownloader(client).Download(
		context.Background(),
		"/input",
		File(path),
		DownloadOptions{},
	)
	if err == nil || !strings.Contains(err.Error(), "content hash mismatch") {
		t.Fatalf("Download() error = %v, want content hash mismatch", err)
	}
	if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("destination exists after failure: %v", statErr)
	}
	matches, globErr := filepath.Glob(filepath.Join(dir, ".output.bin.*.part"))
	if globErr != nil {
		t.Fatal(globErr)
	}
	if len(matches) != 0 {
		t.Fatalf("temporary files remain after failure: %v", matches)
	}
}

type fakeDownloadClient struct {
	mu sync.Mutex

	data     []byte
	metadata *files.FileMetadata
	ranges   []string

	failFirstAt  int
	failFirstErr error
	apiErr       error
	calls        int
}

func (c *fakeDownloadClient) DownloadContext(
	_ context.Context,
	arg *files.DownloadArg,
) (*files.FileMetadata, io.ReadCloser, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.calls++
	rangeHeader := ""
	if arg != nil && arg.ExtraHeaders != nil {
		rangeHeader = arg.ExtraHeaders["Range"]
	}
	c.ranges = append(c.ranges, rangeHeader)
	if c.apiErr != nil {
		return nil, nil, c.apiErr
	}
	if len(c.data) == 0 && rangeHeader == "bytes=0-0" {
		return nil, nil, downloadUnsatisfiableRangeAPIError()
	}

	start, end, err := parseRange(rangeHeader, int64(len(c.data)))
	if err != nil {
		return nil, nil, err
	}
	body := append([]byte(nil), c.data[start:end]...)
	if c.calls == 1 && c.failFirstErr != nil {
		limit := c.failFirstAt
		if limit < 0 || limit > len(body) {
			limit = len(body)
		}
		return c.metadata, &errorReadCloser{
			data: body[:limit],
			err:  c.failFirstErr,
		}, nil
	}
	return c.metadata, io.NopCloser(bytes.NewReader(body)), nil
}

func parseRange(header string, size int64) (int64, int64, error) {
	if header == "" {
		return 0, size, nil
	}
	if !strings.HasPrefix(header, "bytes=") {
		return 0, 0, fmt.Errorf("invalid range %q", header)
	}
	parts := strings.Split(strings.TrimPrefix(header, "bytes="), "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range %q", header)
	}
	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	end := size
	if parts[1] != "" {
		inclusiveEnd, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}
		end = inclusiveEnd + 1
	}
	if start < 0 || end < start || end > size {
		return 0, 0, fmt.Errorf("range %q outside size %d", header, size)
	}
	return start, end, nil
}

type errorReadCloser struct {
	data []byte
	err  error
	done bool
}

func (r *errorReadCloser) Read(p []byte) (int, error) {
	if len(r.data) > 0 {
		n := copy(p, r.data)
		r.data = r.data[n:]
		if len(r.data) == 0 {
			return n, r.err
		}
		return n, nil
	}
	if !r.done {
		r.done = true
		return 0, r.err
	}
	return 0, io.EOF
}

func (*errorReadCloser) Close() error { return nil }

func fileMetadata(data []byte) *files.FileMetadata {
	hash, err := contenthash.Compute(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return &files.FileMetadata{
		Size:        uint64(len(data)),
		Rev:         "rev-1",
		ContentHash: hash,
	}
}

type recordingDownloadTarget struct {
	mu sync.Mutex

	data        []byte
	commitErr   error
	commitHook  func()
	commitCalls int
	abortCalls  int
	abortCtxErr error
}

func (t *recordingDownloadTarget) Prepare(_ context.Context, info DownloadInfo) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.data = make([]byte, info.Size)
	return nil
}

func (t *recordingDownloadTarget) WriteAt(p []byte, offset int64) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return copy(t.data[offset:], p), nil
}

func (t *recordingDownloadTarget) Commit(context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.commitCalls++
	if t.commitHook != nil {
		t.commitHook()
	}
	return t.commitErr
}

func (t *recordingDownloadTarget) Abort(ctx context.Context, _ error) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.abortCalls++
	t.abortCtxErr = ctx.Err()
	return nil
}

func downloadPathAPIError() error {
	return files.DownloadAPIError{
		APIError: dropbox.APIError{
			ErrorSummary: "path/not_found/.",
		},
		EndpointError: &files.DownloadError{
			Tagged: dropbox.Tagged{
				Tag: files.DownloadErrorPath,
			},
		},
	}
}

func downloadUnsatisfiableRangeAPIError() error {
	return files.DownloadAPIError{
		APIError: dropbox.APIError{
			ErrorSummary: "range/not_satisfiable/.",
		},
		EndpointError: &files.DownloadError{
			Tagged: dropbox.Tagged{
				Tag: files.DownloadErrorUnsupportedFile,
			},
		},
	}
}

func assertDownloadProgress(t *testing.T, updates []DownloadProgress, total int64) {
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
			t.Fatalf("TotalBytes = %d, want %d", update.TotalBytes, total)
		}
		if update.BytesCommitted <= previous {
			t.Fatalf("progress is not strictly increasing: previous=%d current=%d", previous, update.BytesCommitted)
		}
		if update.BytesCommitted > total {
			t.Fatalf("progress exceeds total: %d > %d", update.BytesCommitted, total)
		}
		previous = update.BytesCommitted
	}
	if previous != total {
		t.Fatalf("final progress = %d, want %d", previous, total)
	}
}
