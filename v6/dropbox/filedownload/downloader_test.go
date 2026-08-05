package filedownload

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/contenthash"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type fakeClient struct {
	download func(
		context.Context,
		*files.DownloadArg,
	) (*files.FileMetadata, io.ReadCloser, error)
}

func (f *fakeClient) DownloadContext(
	ctx context.Context,
	arg *files.DownloadArg,
) (*files.FileMetadata, io.ReadCloser, error) {
	return f.download(ctx, arg)
}

type failingReader struct {
	data []byte
	err  error
}

func (r *failingReader) Read(p []byte) (int, error) {
	if len(r.data) == 0 {
		return 0, r.err
	}

	n := copy(p, r.data)
	r.data = r.data[n:]
	return n, nil
}

func testMetadata(size uint64) *files.FileMetadata {
	return &files.FileMetadata{
		Metadata: *files.NewMetadata("file.bin"),
		Rev:      "rev1",
		Size:     size,
	}
}

func testMetadataWithRev(size uint64, rev string) *files.FileMetadata {
	metadata := testMetadata(size)
	metadata.Rev = rev
	return metadata
}

func testMetadataWithHash(t *testing.T, payload string) *files.FileMetadata {
	t.Helper()

	hash, err := contenthash.Compute(strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	metadata := testMetadata(uint64(len(payload)))
	metadata.ContentHash = hash
	return metadata
}

func TestDownloadFileResumesFromPartFile(t *testing.T) {
	dir := t.TempDir()
	localPath := filepath.Join(dir, "file.bin")
	partPath := localPath + ".part"

	if err := os.WriteFile(partPath, []byte("hello "), 0o644); err != nil {
		t.Fatal(err)
	}

	client := &fakeClient{
		download: func(
			ctx context.Context,
			arg *files.DownloadArg,
		) (*files.FileMetadata, io.ReadCloser, error) {
			if got := arg.ExtraHeaders["Range"]; got != "bytes=6-" {
				t.Fatalf("Range header = %q, want %q", got, "bytes=6-")
			}

			return testMetadata(11), io.NopCloser(strings.NewReader("world")), nil
		},
	}

	d := New(client)

	result, err := d.DownloadFile(
		context.Background(),
		"/file.bin",
		localPath,
	)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(localPath)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(data), "hello world"; got != want {
		t.Fatalf("file contents = %q, want %q", got, want)
	}

	if result.ResumedFrom != 6 {
		t.Fatalf("ResumedFrom = %d, want 6", result.ResumedFrom)
	}

	if _, err := os.Stat(partPath); !os.IsNotExist(err) {
		t.Fatalf("part file still exists: %v", err)
	}
}

func TestDownloadFileFresh(t *testing.T) {
	dir := t.TempDir()
	localPath := filepath.Join(dir, "file.bin")

	client := &fakeClient{
		download: func(
			ctx context.Context,
			arg *files.DownloadArg,
		) (*files.FileMetadata, io.ReadCloser, error) {
			if got := arg.ExtraHeaders["Range"]; got != "" {
				t.Fatalf("Range header = %q, want empty", got)
			}

			return testMetadata(5),
				io.NopCloser(strings.NewReader("hello")),
				nil
		},
	}

	d := New(client)

	result, err := d.DownloadFile(
		context.Background(),
		"/file.bin",
		localPath,
	)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(localPath)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(data), "hello"; got != want {
		t.Fatalf("file contents = %q, want %q", got, want)
	}

	if result.ResumedFrom != 0 {
		t.Fatalf("ResumedFrom = %d, want 0", result.ResumedFrom)
	}

	if _, err := os.Stat(localPath + ".part"); !os.IsNotExist(err) {
		t.Fatalf("part file still exists: %v", err)
	}
}

func TestDownloadFilePreservesPartFileOnReadError(t *testing.T) {
	dir := t.TempDir()
	localPath := filepath.Join(dir, "file.bin")
	partPath := localPath + ".part"

	client := &fakeClient{
		download: func(
			ctx context.Context,
			arg *files.DownloadArg,
		) (*files.FileMetadata, io.ReadCloser, error) {
			return testMetadata(7),
				io.NopCloser(&failingReader{
					data: []byte("partial"),
					err:  errors.New("connection reset"),
				}),
				nil
		},
	}

	d := New(client)

	_, err := d.downloadFileAttempt(
		context.Background(),
		"/file.bin",
		localPath,
		"",
	)
	if err == nil {
		t.Fatal("DownloadFile() expected an error")
	}

	data, readErr := os.ReadFile(partPath)
	if readErr != nil {
		t.Fatal(readErr)
	}

	if got, want := string(data), "partial"; got != want {
		t.Fatalf("part file contents = %q, want %q", got, want)
	}

	if _, statErr := os.Stat(localPath); !os.IsNotExist(statErr) {
		t.Fatalf("final file unexpectedly exists: %v", statErr)
	}
}

func TestDownloadFileRetriesAndResumes(t *testing.T) {
	dir := t.TempDir()
	localPath := filepath.Join(dir, "file.bin")

	calls := 0
	client := &fakeClient{
		download: func(
			ctx context.Context,
			arg *files.DownloadArg,
		) (*files.FileMetadata, io.ReadCloser, error) {
			calls++

			switch calls {
			case 1:
				if got := arg.ExtraHeaders["Range"]; got != "" {
					t.Fatalf("first Range header = %q, want empty", got)
				}

				return testMetadata(11),
					io.NopCloser(&failingReader{
						data: []byte("hello "),
						err:  errors.New("connection reset"),
					}),
					nil

			case 2:
				if got := arg.ExtraHeaders["Range"]; got != "bytes=6-" {
					t.Fatalf("second Range header = %q, want %q", got, "bytes=6-")
				}

				return testMetadata(11),
					io.NopCloser(strings.NewReader("world")),
					nil

			default:
				t.Fatalf("unexpected download call %d", calls)
				return nil, nil, nil
			}
		},
	}

	d := New(client)

	result, err := d.DownloadFile(
		context.Background(),
		"/file.bin",
		localPath,
	)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(localPath)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(data), "hello world"; got != want {
		t.Fatalf("file contents = %q, want %q", got, want)
	}

	if calls != 2 {
		t.Fatalf("download calls = %d, want 2", calls)
	}

	if result.ResumedFrom != 6 {
		t.Fatalf("ResumedFrom = %d, want 6", result.ResumedFrom)
	}
}

func TestDownloadFileStopsOnContextCancellation(t *testing.T) {
	dir := t.TempDir()
	localPath := filepath.Join(dir, "file.bin")

	ctx, cancel := context.WithCancel(context.Background())
	calls := 0

	client := &fakeClient{
		download: func(
			ctx context.Context,
			arg *files.DownloadArg,
		) (*files.FileMetadata, io.ReadCloser, error) {
			calls++
			cancel()

			return nil, nil, ctx.Err()
		},
	}

	d := New(client)

	_, err := d.DownloadFile(ctx, "/file.bin", localPath)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("DownloadFile() error = %v, want context.Canceled", err)
	}

	if calls != 1 {
		t.Fatalf("download calls = %d, want 1", calls)
	}
}

func TestDownloadFileStopsAfterMaxAttempts(t *testing.T) {
	dir := t.TempDir()
	localPath := filepath.Join(dir, "file.bin")

	calls := 0
	wantErr := errors.New("connection reset")

	client := &fakeClient{
		download: func(
			ctx context.Context,
			arg *files.DownloadArg,
		) (*files.FileMetadata, io.ReadCloser, error) {
			calls++

			return testMetadata(0),
				io.NopCloser(&failingReader{
					err: wantErr,
				}),
				nil
		},
	}

	d := New(client)

	_, err := d.DownloadFile(
		context.Background(),
		"/file.bin",
		localPath,
	)
	if err == nil {
		t.Fatal("DownloadFile() expected an error")
	}

	if calls != d.maxAttempts {
		t.Fatalf(
			"download calls = %d, want %d",
			calls,
			d.maxAttempts,
		)
	}
}

func TestDownloadFileStopsWhenRemoteFileChangesDuringRetry(t *testing.T) {
	dir := t.TempDir()
	localPath := filepath.Join(dir, "file.bin")
	partPath := localPath + ".part"

	calls := 0
	client := &fakeClient{
		download: func(
			ctx context.Context,
			arg *files.DownloadArg,
		) (*files.FileMetadata, io.ReadCloser, error) {
			calls++

			switch calls {
			case 1:
				return testMetadataWithRev(11, "rev1"),
					io.NopCloser(&failingReader{
						data: []byte("hello "),
						err:  errors.New("connection reset"),
					}),
					nil
			case 2:
				return testMetadataWithRev(11, "rev2"),
					io.NopCloser(strings.NewReader("world")),
					nil
			default:
				t.Fatalf("unexpected download call %d", calls)
				return nil, nil, nil
			}
		},
	}

	d := New(client)

	_, err := d.DownloadFile(
		context.Background(),
		"/file.bin",
		localPath,
	)
	if err == nil || !strings.Contains(err.Error(), "remote file changed") {
		t.Fatalf("DownloadFile() error = %v, want remote changed error", err)
	}

	if _, statErr := os.Stat(partPath); !os.IsNotExist(statErr) {
		t.Fatalf("part file still exists: %v", statErr)
	}
}

func TestDownloadFileValidatesTotalSize(t *testing.T) {
	dir := t.TempDir()
	localPath := filepath.Join(dir, "file.bin")
	partPath := localPath + ".part"

	client := &fakeClient{
		download: func(
			ctx context.Context,
			arg *files.DownloadArg,
		) (*files.FileMetadata, io.ReadCloser, error) {
			return testMetadata(10),
				io.NopCloser(strings.NewReader("short")),
				nil
		},
	}

	d := New(client, WithMaxAttempts(1))

	_, err := d.DownloadFile(
		context.Background(),
		"/file.bin",
		localPath,
	)
	if err == nil || !strings.Contains(err.Error(), "incomplete download") {
		t.Fatalf("DownloadFile() error = %v, want incomplete download error", err)
	}
	if _, statErr := os.Stat(partPath); !os.IsNotExist(statErr) {
		t.Fatalf("part file still exists: %v", statErr)
	}
}

func TestDownloadFileValidatesContentHash(t *testing.T) {
	dir := t.TempDir()
	localPath := filepath.Join(dir, "file.bin")
	partPath := localPath + ".part"

	client := &fakeClient{
		download: func(
			ctx context.Context,
			arg *files.DownloadArg,
		) (*files.FileMetadata, io.ReadCloser, error) {
			metadata := testMetadata(5)
			metadata.ContentHash = "not-the-right-hash"
			return metadata, io.NopCloser(strings.NewReader("hello")), nil
		},
	}

	d := New(client, WithMaxAttempts(1))

	_, err := d.DownloadFile(
		context.Background(),
		"/file.bin",
		localPath,
	)
	if err == nil || !strings.Contains(err.Error(), "content hash mismatch") {
		t.Fatalf("DownloadFile() error = %v, want hash mismatch error", err)
	}
	if _, statErr := os.Stat(partPath); !os.IsNotExist(statErr) {
		t.Fatalf("part file still exists: %v", statErr)
	}
}

func TestDownloadFileAcceptsContentHash(t *testing.T) {
	dir := t.TempDir()
	localPath := filepath.Join(dir, "file.bin")

	client := &fakeClient{
		download: func(
			ctx context.Context,
			arg *files.DownloadArg,
		) (*files.FileMetadata, io.ReadCloser, error) {
			return testMetadataWithHash(t, "hello"),
				io.NopCloser(strings.NewReader("hello")),
				nil
		},
	}

	d := New(client)

	if _, err := d.DownloadFile(context.Background(), "/file.bin", localPath); err != nil {
		t.Fatal(err)
	}
}

func TestDownloadFileReportsProgress(t *testing.T) {
	dir := t.TempDir()
	localPath := filepath.Join(dir, "file.bin")
	var updates []Progress

	client := &fakeClient{
		download: func(
			ctx context.Context,
			arg *files.DownloadArg,
		) (*files.FileMetadata, io.ReadCloser, error) {
			return testMetadata(11),
				io.NopCloser(strings.NewReader("hello world")),
				nil
		},
	}

	d := New(client, WithProgress(func(progress Progress) {
		updates = append(updates, progress)
	}))

	if _, err := d.DownloadFile(context.Background(), "/file.bin", localPath); err != nil {
		t.Fatal(err)
	}

	if len(updates) == 0 {
		t.Fatal("progress callback was not called")
	}

	last := updates[len(updates)-1]
	if last.BytesWritten != 11 || last.TotalBytes != 11 || last.ResumedFrom != 0 {
		t.Fatalf("last progress = %+v, want 11/11 from 0", last)
	}
}

func TestDownloadFileParallelDownloadsRanges(t *testing.T) {
	dir := t.TempDir()
	localPath := filepath.Join(dir, "file.bin")
	payload := "hello world"
	var mu sync.Mutex
	var ranges []string

	client := &fakeClient{
		download: func(
			ctx context.Context,
			arg *files.DownloadArg,
		) (*files.FileMetadata, io.ReadCloser, error) {
			header := arg.ExtraHeaders["Range"]
			mu.Lock()
			ranges = append(ranges, header)
			mu.Unlock()

			start, end := parseRangeHeader(t, header)
			return testMetadataWithHash(t, payload),
				io.NopCloser(strings.NewReader(payload[start : end+1])),
				nil
		},
	}

	d := New(client, WithParallelDownloads(3))

	if _, err := d.DownloadFile(context.Background(), "/file.bin", localPath); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(localPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(data); got != payload {
		t.Fatalf("file contents = %q, want %q", got, payload)
	}

	mu.Lock()
	defer mu.Unlock()
	sort.Strings(ranges)
	wantRanges := []string{
		"bytes=0-0",
		"bytes=1-4",
		"bytes=5-7",
		"bytes=8-10",
	}
	sort.Strings(wantRanges)
	if strings.Join(ranges, ",") != strings.Join(wantRanges, ",") {
		t.Fatalf("ranges = %v, want %v", ranges, wantRanges)
	}
}

func TestDownloadFileParallelRemovesPartFileOnRangeError(t *testing.T) {
	dir := t.TempDir()
	localPath := filepath.Join(dir, "file.bin")
	partPath := localPath + ".part"
	payload := "hello world"

	client := &fakeClient{
		download: func(
			ctx context.Context,
			arg *files.DownloadArg,
		) (*files.FileMetadata, io.ReadCloser, error) {
			header := arg.ExtraHeaders["Range"]
			if header == "bytes=0-0" {
				return testMetadataWithHash(t, payload),
					io.NopCloser(strings.NewReader(payload[:1])),
					nil
			}
			return nil, nil, errors.New("range failed")
		},
	}

	d := New(client, WithMaxAttempts(1), WithParallelDownloads(3))

	_, err := d.DownloadFile(context.Background(), "/file.bin", localPath)
	if err == nil {
		t.Fatal("DownloadFile() expected an error")
	}

	if _, statErr := os.Stat(partPath); !os.IsNotExist(statErr) {
		t.Fatalf("part file still exists: %v", statErr)
	}
}

func parseRangeHeader(t *testing.T, header string) (int, int) {
	t.Helper()

	const prefix = "bytes="
	if !strings.HasPrefix(header, prefix) {
		t.Fatalf("Range header = %q, want bytes range", header)
	}

	parts := strings.Split(strings.TrimPrefix(header, prefix), "-")
	if len(parts) != 2 {
		t.Fatalf("Range header = %q, want start-end", header)
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		t.Fatal(err)
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		t.Fatal(err)
	}
	return start, end
}
