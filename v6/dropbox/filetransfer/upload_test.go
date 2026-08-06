package filetransfer

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

func TestUploadSequentialKnownSize(t *testing.T) {
	data := []byte("hello upload")
	client := newFakeUploadClient()

	var updates []UploadProgress
	result, err := NewUploader(client).Upload(
		context.Background(),
		BytesUpload(data),
		files.NewCommitInfo("/uploaded.txt"),
		UploadOptions{
			Progress: func(progress UploadProgress) {
				updates = append(updates, progress)
			},
		},
	)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}
	if result == nil || result.Metadata == nil {
		t.Fatal("Upload() returned nil metadata")
	}
	if got := client.uploadedBytes(); !bytes.Equal(got, data) {
		t.Fatalf("uploaded bytes = %q, want %q", got, data)
	}
	assertUploadProgress(t, updates, int64(len(data)), int64(len(data)))
}

func TestUploadSequentialUnknownSize(t *testing.T) {
	data := append(bytes.Repeat([]byte{'x'}, int(uploadChunkSize)), []byte("tail")...)
	source, err := ReaderUpload(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	client := newFakeUploadClient()

	var updates []UploadProgress
	_, err = NewUploader(client).Upload(
		context.Background(),
		source,
		files.NewCommitInfo("/stream.bin"),
		UploadOptions{
			Progress: func(progress UploadProgress) {
				updates = append(updates, progress)
			},
		},
	)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}
	if got := client.uploadedBytes(); !bytes.Equal(got, data) {
		t.Fatalf("uploaded bytes differ: got %d, want %d", len(got), len(data))
	}
	assertUploadProgress(t, updates, int64(len(data)), -1)
}

func TestUploadRetriesBufferedChunk(t *testing.T) {
	data := []byte("retry the exact same bytes")
	client := newFakeUploadClient()
	client.finishFailures = 1

	_, err := NewUploader(client).Upload(
		context.Background(),
		BytesUpload(data),
		files.NewCommitInfo("/retry.bin"),
		UploadOptions{MaxAttempts: 2},
	)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}
	if client.finishCalls != 2 {
		t.Fatalf("finish calls = %d, want 2", client.finishCalls)
	}
	if len(client.finishBodies) != 2 || !bytes.Equal(client.finishBodies[0], client.finishBodies[1]) {
		t.Fatalf("finish retry bodies differ: %#v", client.finishBodies)
	}
}

func TestUploadDoesNotRetryPermanentAppendAPIError(t *testing.T) {
	data := append(bytes.Repeat([]byte("a"), int(uploadChunkSize)), []byte("tail")...)
	client := newFakeUploadClient()
	client.appendErr = uploadAppendClosedError()

	_, err := NewUploader(client).Upload(
		context.Background(),
		BytesUpload(data),
		files.NewCommitInfo("/closed.bin"),
		UploadOptions{MaxAttempts: 3},
	)
	if err == nil {
		t.Fatal("Upload() error = nil, want append API error")
	}
	if client.appendCalls != 1 {
		t.Fatalf("append calls = %d, want 1", client.appendCalls)
	}
}

func TestUploadDoesNotRetryPermanentFinishAPIError(t *testing.T) {
	client := newFakeUploadClient()
	client.finishErr = uploadFinishPathError()

	_, err := NewUploader(client).Upload(
		context.Background(),
		BytesUpload([]byte("data")),
		files.NewCommitInfo("/path-error.bin"),
		UploadOptions{MaxAttempts: 3},
	)
	if err == nil {
		t.Fatal("Upload() error = nil, want finish API error")
	}
	if client.finishCalls != 1 {
		t.Fatalf("finish calls = %d, want 1", client.finishCalls)
	}
}

func TestUploadTreatsCommittedAppendIncorrectOffsetAsSuccess(t *testing.T) {
	data := append(bytes.Repeat([]byte("a"), int(uploadChunkSize)), []byte("tail")...)
	client := newFakeUploadClient()
	client.appendCommittedLostResponse = true

	_, err := NewUploader(client).Upload(
		context.Background(),
		BytesUpload(data),
		files.NewCommitInfo("/retry-append.bin"),
		UploadOptions{MaxAttempts: 2},
	)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}
	if got := client.uploadedBytes(); !bytes.Equal(got, data) {
		t.Fatalf("uploaded bytes differ: got %d, want %d", len(got), len(data))
	}
	if client.appendCalls != 1 {
		t.Fatalf("append calls = %d, want 1", client.appendCalls)
	}
}

func TestUploadRetriesFinishEmptyAfterCommittedFinalBody(t *testing.T) {
	data := []byte("final body committed")
	client := newFakeUploadClient()
	client.finishCommittedLostResponse = true

	_, err := NewUploader(client).Upload(
		context.Background(),
		BytesUpload(data),
		files.NewCommitInfo("/retry-finish.bin"),
		UploadOptions{MaxAttempts: 2},
	)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}
	if got := client.uploadedBytes(); !bytes.Equal(got, data) {
		t.Fatalf("uploaded bytes = %q, want %q", got, data)
	}
	if client.finishCalls != 2 {
		t.Fatalf("finish calls = %d, want 2", client.finishCalls)
	}
	if len(client.finishBodies) != 2 {
		t.Fatalf("finish bodies = %d, want 2", len(client.finishBodies))
	}
	if !bytes.Equal(client.finishBodies[0], data) {
		t.Fatalf("first finish body = %q, want %q", client.finishBodies[0], data)
	}
	if len(client.finishBodies[1]) != 0 {
		t.Fatalf("second finish body = %q, want empty retry body", client.finishBodies[1])
	}
}

func TestUploadParallel(t *testing.T) {
	data := append(bytes.Repeat([]byte("a"), int(2*uploadChunkSize)), []byte("tail")...)
	client := newFakeUploadClient()

	var (
		progressMu sync.Mutex
		updates    []UploadProgress
	)
	_, err := NewUploader(client).Upload(
		context.Background(),
		BytesUpload(data),
		files.NewCommitInfo("/parallel.bin"),
		UploadOptions{
			ParallelUploads: 3,
			Progress: func(progress UploadProgress) {
				progressMu.Lock()
				updates = append(updates, progress)
				progressMu.Unlock()
			},
		},
	)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}
	if got := client.uploadedBytes(); !bytes.Equal(got, data) {
		t.Fatalf("uploaded bytes differ: got %d, want %d", len(got), len(data))
	}
	if !client.concurrentSession {
		t.Fatal("parallel upload did not request a concurrent session")
	}

	progressMu.Lock()
	captured := append([]UploadProgress(nil), updates...)
	progressMu.Unlock()
	assertUploadProgress(t, captured, int64(len(data)), int64(len(data)))
}

func TestUploadParallelZeroByteFallsBackToSequential(t *testing.T) {
	client := newFakeUploadClient()

	_, err := NewUploader(client).Upload(
		context.Background(),
		BytesUpload(nil),
		files.NewCommitInfo("/empty.bin"),
		UploadOptions{ParallelUploads: 3},
	)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}
	if client.concurrentSession {
		t.Fatal("zero-byte upload used a concurrent session")
	}
	if got := client.uploadedBytes(); len(got) != 0 {
		t.Fatalf("uploaded bytes = %q, want empty", got)
	}
}

func TestUploadRejectsParallelOneShotSource(t *testing.T) {
	source, err := ReaderUpload(strings.NewReader("data"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewUploader(newFakeUploadClient()).Upload(
		context.Background(),
		source,
		files.NewCommitInfo("/data"),
		UploadOptions{ParallelUploads: 2},
	)
	if err == nil || !strings.Contains(err.Error(), "ranged upload source") {
		t.Fatalf("Upload() error = %v, want ranged source validation", err)
	}
}

func TestUploadRejectsDeclaredSizeMismatch(t *testing.T) {
	source, err := SizedReaderUpload(strings.NewReader("abc"), 5)
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewUploader(newFakeUploadClient()).Upload(
		context.Background(),
		source,
		files.NewCommitInfo("/data"),
		UploadOptions{},
	)
	if err == nil || !strings.Contains(err.Error(), "read upload content") {
		t.Fatalf("Upload() error = %v, want declared size mismatch", err)
	}
}

func assertUploadProgress(t *testing.T, updates []UploadProgress, committed, total int64) {
	t.Helper()
	if committed == 0 {
		return
	}
	if len(updates) == 0 {
		t.Fatal("no upload progress updates")
	}
	previous := int64(0)
	for _, update := range updates {
		if update.TotalBytes != total {
			t.Fatalf("TotalBytes = %d, want %d", update.TotalBytes, total)
		}
		if update.BytesCommitted <= previous {
			t.Fatalf("progress is not strictly increasing: previous=%d current=%d", previous, update.BytesCommitted)
		}
		if total >= 0 && update.BytesCommitted > total {
			t.Fatalf("progress exceeds total: %d > %d", update.BytesCommitted, total)
		}
		previous = update.BytesCommitted
	}
	if previous != committed {
		t.Fatalf("final progress = %d, want %d", previous, committed)
	}
}

type fakeUploadClient struct {
	mu sync.Mutex

	sessionID string
	chunks    map[int64][]byte

	concurrentSession           bool
	appendErr                   error
	appendCommittedLostResponse bool
	appendCalls                 int
	finishFailures              int
	finishErr                   error
	finishCommittedLostResponse bool
	finishCalls                 int
	finishBodies                [][]byte
}

func newFakeUploadClient() *fakeUploadClient {
	return &fakeUploadClient{
		sessionID: "session-1",
		chunks:    make(map[int64][]byte),
	}
}

func (c *fakeUploadClient) UploadContext(
	context.Context,
	*files.UploadArg,
	io.Reader,
) (*files.FileMetadata, error) {
	return nil, errors.New("unexpected direct upload")
}

func (c *fakeUploadClient) UploadSessionStartContext(
	_ context.Context,
	arg *files.UploadSessionStartArg,
	content io.Reader,
) (*files.UploadSessionStartResult, error) {
	data, err := io.ReadAll(content)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(data) != 0 {
		c.chunks[0] = append([]byte(nil), data...)
	}
	if arg != nil && arg.SessionType != nil && arg.SessionType.Tag == files.UploadSessionTypeConcurrent {
		c.concurrentSession = true
	}
	return &files.UploadSessionStartResult{SessionId: c.sessionID}, nil
}

func (c *fakeUploadClient) UploadSessionAppendV2Context(
	_ context.Context,
	arg *files.UploadSessionAppendArg,
	content io.Reader,
) error {
	data, err := io.ReadAll(content)
	if err != nil {
		return err
	}
	if arg == nil || arg.Cursor == nil {
		return errors.New("append cursor is required")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.appendCalls++
	if c.appendErr != nil {
		return c.appendErr
	}
	c.chunks[int64(arg.Cursor.Offset)] = append([]byte(nil), data...)
	if c.appendCommittedLostResponse && c.appendCalls == 1 {
		return uploadAppendIncorrectOffsetError(
			arg.Cursor.Offset + uint64(len(data)),
		)
	}
	return nil
}

func (c *fakeUploadClient) UploadSessionFinishContext(
	_ context.Context,
	arg *files.UploadSessionFinishArg,
	content io.Reader,
) (*files.FileMetadata, error) {
	data, err := io.ReadAll(content)
	if err != nil {
		return nil, err
	}
	if arg == nil || arg.Cursor == nil || arg.Commit == nil {
		return nil, errors.New("finish arguments are required")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.finishCalls++
	c.finishBodies = append(c.finishBodies, append([]byte(nil), data...))
	if c.finishErr != nil {
		return nil, c.finishErr
	}
	if c.finishCalls <= c.finishFailures {
		return nil, authServerError()
	}
	if len(data) != 0 {
		c.chunks[int64(arg.Cursor.Offset)] = append([]byte(nil), data...)
	}
	if c.finishCommittedLostResponse && c.finishCalls == 1 {
		return nil, uploadFinishIncorrectOffsetError(
			arg.Cursor.Offset + uint64(len(data)),
		)
	}
	assembled := assembleChunks(c.chunks)
	return &files.FileMetadata{Size: uint64(len(assembled)), Rev: "uploaded-rev"}, nil
}

func uploadAppendClosedError() error {
	return files.UploadSessionAppendV2APIError{
		APIError: dropbox.APIError{
			ErrorSummary: "closed/.",
		},
		EndpointError: &files.UploadSessionAppendError{
			Tagged: dropbox.Tagged{
				Tag: files.UploadSessionAppendErrorClosed,
			},
		},
	}
}

func uploadFinishPathError() error {
	return files.UploadSessionFinishAPIError{
		APIError: dropbox.APIError{
			ErrorSummary: "path/conflict/.",
		},
		EndpointError: &files.UploadSessionFinishError{
			Tagged: dropbox.Tagged{
				Tag: files.UploadSessionFinishErrorPath,
			},
		},
	}
}

func authServerError() error {
	return auth.ServerError{
		APIError: dropbox.APIError{
			ErrorSummary: "server error",
		},
		StatusCode: 503,
	}
}

func uploadAppendIncorrectOffsetError(correctOffset uint64) error {
	return files.UploadSessionAppendV2APIError{
		APIError: dropbox.APIError{
			ErrorSummary: "incorrect_offset/.",
		},
		EndpointError: &files.UploadSessionAppendError{
			Tagged: dropbox.Tagged{
				Tag: files.UploadSessionAppendErrorIncorrectOffset,
			},
			IncorrectOffset: files.NewUploadSessionOffsetError(correctOffset),
		},
	}
}

func uploadFinishIncorrectOffsetError(correctOffset uint64) error {
	return files.UploadSessionFinishAPIError{
		APIError: dropbox.APIError{
			ErrorSummary: "lookup_failed/incorrect_offset/.",
		},
		EndpointError: &files.UploadSessionFinishError{
			Tagged: dropbox.Tagged{
				Tag: files.UploadSessionFinishErrorLookupFailed,
			},
			LookupFailed: &files.UploadSessionLookupError{
				Tagged: dropbox.Tagged{
					Tag: files.UploadSessionLookupErrorIncorrectOffset,
				},
				IncorrectOffset: files.NewUploadSessionOffsetError(correctOffset),
			},
		},
	}
}

func (c *fakeUploadClient) uploadedBytes() []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	return assembleChunks(c.chunks)
}

func assembleChunks(chunks map[int64][]byte) []byte {
	offsets := make([]int64, 0, len(chunks))
	for offset := range chunks {
		offsets = append(offsets, offset)
	}
	sort.Slice(offsets, func(i, j int) bool { return offsets[i] < offsets[j] })

	var result []byte
	for _, offset := range offsets {
		if int64(len(result)) < offset {
			result = append(result, make([]byte, offset-int64(len(result)))...)
		}
		chunk := chunks[offset]
		end := offset + int64(len(chunk))
		if int64(len(result)) < end {
			result = append(result, make([]byte, end-int64(len(result)))...)
		}
		copy(result[offset:end], chunk)
	}
	return result
}
