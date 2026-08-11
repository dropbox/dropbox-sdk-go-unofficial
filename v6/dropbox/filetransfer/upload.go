package filetransfer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"sync"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

// -----------------------------------------------------------------------------
// Upload
// -----------------------------------------------------------------------------

// UploadOptions configures an upload operation.
//
// The zero value uses default retry settings.
type UploadOptions struct {
	// MaxAttempts is the maximum number of attempts for retryable failures.
	// A non-positive value uses the default.
	MaxAttempts int

	// ParallelUploads is the maximum number of concurrent upload requests.
	//
	// Parallel uploads require a RangedUploadSource. Values less than 2 disable
	// parallel uploads.
	ParallelUploads int

	// Progress receives upload progress updates.
	Progress UploadProgressFunc
}

// -----------------------------------------------------------------------------
// Uploader
// -----------------------------------------------------------------------------

// UploadResult describes a completed upload.
type UploadResult struct {
	Metadata *files.FileMetadata
}

// UploadClient is the subset of files.Client required by Uploader.
type UploadClient interface {
	UploadSessionStartContext(
		ctx context.Context,
		arg *files.UploadSessionStartArg,
		content io.Reader,
	) (*files.UploadSessionStartResult, error)

	UploadSessionAppendV2Context(
		ctx context.Context,
		arg *files.UploadSessionAppendArg,
		content io.Reader,
	) error

	UploadSessionFinishContext(
		ctx context.Context,
		arg *files.UploadSessionFinishArg,
		content io.Reader,
	) (*files.FileMetadata, error)
}

// Uploader uploads files to Dropbox.
type Uploader struct {
	client UploadClient
}

// NewUploader creates an Uploader.
func NewUploader(client UploadClient) *Uploader {
	return &Uploader{
		client: client,
	}
}

// Upload uploads source to Dropbox using commit.
//
// Upload returns an error before opening source if commit is nil or does not
// contain a destination path.
//
// Upload uses a sequential upload session when the source size is unknown.
// Each chunk is retained until Dropbox confirms it, allowing retryable failures
// to retry the current chunk without reopening the source.
//
// Upload may use concurrent upload requests when ParallelUploads is greater
// than one and source implements RangedUploadSource.
//
// Chunk sizing is managed by the uploader.
//
// Progress reports only bytes confirmed by Dropbox. Retried bytes are not
// counted more than once.
func (u *Uploader) Upload(
	ctx context.Context,
	source UploadSource,
	commit *files.CommitInfo,
	options UploadOptions,
) (*UploadResult, error) {
	if u == nil || u.client == nil {
		return nil, errors.New("upload client is required")
	}
	if source == nil {
		return nil, errors.New("upload source is required")
	}
	if commit == nil {
		return nil, errors.New("upload commit info is required")
	}
	if commit.Path == "" {
		return nil, errors.New("upload destination path is required")
	}
	if options.ParallelUploads > 1 {
		if _, ok := source.(RangedUploadSource); !ok {
			return nil, errors.New(
				"parallel uploads require a ranged upload source",
			)
		}
	}

	maxAttempts := options.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = defaultMaxAttempts
	}

	if options.ParallelUploads > 1 {
		ranged := source.(RangedUploadSource)
		if ranged.Size() == 0 {
			return u.uploadSequential(
				ctx,
				source,
				commit,
				maxAttempts,
				options.Progress,
			)
		}
		return u.uploadParallel(
			ctx,
			ranged,
			commit,
			maxAttempts,
			options.ParallelUploads,
			options.Progress,
		)
	}

	return u.uploadSequential(
		ctx,
		source,
		commit,
		maxAttempts,
		options.Progress,
	)
}

const uploadChunkSize = int64(8 * 1024 * 1024)

func (u *Uploader) startUploadSessionWithRetry(
	ctx context.Context,
	arg *files.UploadSessionStartArg,
	maxAttempts int,
) (*files.UploadSessionStartResult, error) {
	var lastErr error

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		start, err := u.client.UploadSessionStartContext(
			ctx,
			arg,
			bytes.NewReader(nil),
		)
		if err == nil {
			if start == nil || start.SessionId == "" {
				return nil, errors.New("upload session id is empty")
			}
			return start, nil
		}

		if !isRetryableTransferError(err) {
			return nil, err
		}

		lastErr = err
		if err := waitForRetry(ctx, attempt, maxAttempts); err != nil {
			return nil, err
		}
	}

	if lastErr == nil {
		lastErr = errors.New("upload session start failed")
	}
	return nil, lastErr
}

func (u *Uploader) uploadSequential(
	ctx context.Context,
	source UploadSource,
	commit *files.CommitInfo,
	maxAttempts int,
	progress UploadProgressFunc,
) (*UploadResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	total := int64(-1)
	if sized, ok := source.(SizedUploadSource); ok {
		total = sized.Size()
		if total < 0 {
			return nil, errors.New("upload size must not be negative")
		}
	}
	tracker := newUploadProgressTracker(total, progress)

	reader, err := source.Open(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	start, err := u.startUploadSessionWithRetry(
		ctx,
		files.NewUploadSessionStartArg(),
		maxAttempts,
	)
	if err != nil {
		return nil, err
	}

	offset := int64(0)
	for {
		chunk, readErr := readUploadChunk(reader, uploadChunkSize)
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return nil, fmt.Errorf("read upload content: %w", readErr)
		}

		if total >= 0 && offset+int64(len(chunk)) > total {
			return nil, fmt.Errorf(
				"read upload content: got more than declared size %d",
				total,
			)
		}

		if readErr == nil && len(chunk) == 0 {
			return nil, errors.New("read upload content: no progress")
		}

		if errors.Is(readErr, io.EOF) {
			if total >= 0 && offset+int64(len(chunk)) != total {
				return nil, fmt.Errorf(
					"read upload content: got %d bytes, expected %d",
					offset+int64(len(chunk)),
					total,
				)
			}

			metadata, err := u.finishUploadWithRetry(
				ctx,
				start.SessionId,
				offset,
				commit,
				chunk,
				maxAttempts,
			)
			if err != nil {
				return nil, err
			}
			if metadata == nil {
				return nil, errors.New("upload metadata is nil")
			}
			tracker.add(int64(len(chunk)))
			return &UploadResult{Metadata: metadata}, nil
		}

		if err := u.appendUploadWithRetry(
			ctx,
			start.SessionId,
			offset,
			chunk,
			false,
			maxAttempts,
		); err != nil {
			return nil, err
		}
		offset += int64(len(chunk))
		tracker.add(int64(len(chunk)))
	}
}

func (u *Uploader) uploadParallel(
	ctx context.Context,
	source RangedUploadSource,
	commit *files.CommitInfo,
	maxAttempts int,
	parallelUploads int,
	progress UploadProgressFunc,
) (*UploadResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	size := source.Size()
	if size < 0 {
		return nil, errors.New("upload size must not be negative")
	}

	startArg := files.NewUploadSessionStartArg()
	startArg.SessionType = &files.UploadSessionType{
		Tagged: dropbox.Tagged{Tag: files.UploadSessionTypeConcurrent},
	}
	start, err := u.startUploadSessionWithRetry(
		ctx,
		startArg,
		maxAttempts,
	)
	if err != nil {
		return nil, err
	}

	tracker := newUploadProgressTracker(size, progress)
	ranges := splitUploadRanges(size)
	if len(ranges) > 0 {
		finalRange := ranges[len(ranges)-1]
		nonFinalRanges := ranges[:len(ranges)-1]

		if len(nonFinalRanges) > 0 {
			workerCtx, cancel := context.WithCancel(ctx)
			defer cancel()

			jobs := make(chan uploadByteRange)
			errCh := make(chan error, len(nonFinalRanges))

			workers := parallelUploads
			if workers > len(nonFinalRanges) {
				workers = len(nonFinalRanges)
			}
			var wg sync.WaitGroup
			for i := 0; i < workers; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for r := range jobs {
						err := u.uploadRange(
							workerCtx,
							source,
							start.SessionId,
							r,
							maxAttempts,
							tracker,
						)
						if err != nil {
							cancel()
						}
						errCh <- err
					}
				}()
			}

		sendJobs:
			for _, r := range nonFinalRanges {
				select {
				case <-workerCtx.Done():
					break sendJobs
				case jobs <- r:
				}
			}
			close(jobs)

			wg.Wait()
			close(errCh)
			for rangeErr := range errCh {
				if rangeErr != nil {
					return nil, rangeErr
				}
			}
		}

		if err := u.uploadRange(
			ctx,
			source,
			start.SessionId,
			finalRange,
			maxAttempts,
			tracker,
		); err != nil {
			return nil, err
		}
	}

	if tracker.committedBytes() != size {
		return nil, fmt.Errorf(
			"incomplete upload: committed %d of %d bytes",
			tracker.committedBytes(),
			size,
		)
	}

	metadata, err := u.finishUploadWithRetry(
		ctx,
		start.SessionId,
		size,
		commit,
		nil,
		maxAttempts,
	)
	if err != nil {
		return nil, err
	}
	if metadata == nil {
		return nil, errors.New("upload metadata is nil")
	}

	return &UploadResult{Metadata: metadata}, nil
}

func (u *Uploader) uploadRange(
	ctx context.Context,
	source RangedUploadSource,
	sessionID string,
	r uploadByteRange,
	maxAttempts int,
	tracker *uploadProgressTracker,
) error {
	reader, err := source.OpenRange(ctx, r.offset, r.length)
	if err != nil {
		return err
	}
	data, readErr := io.ReadAll(reader)
	closeErr := reader.Close()
	if readErr != nil {
		return fmt.Errorf("read upload range: %w", readErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close upload range: %w", closeErr)
	}
	if int64(len(data)) != r.length {
		return fmt.Errorf(
			"read upload range: got %d bytes, expected %d",
			len(data),
			r.length,
		)
	}

	if err := u.appendUploadWithRetry(
		ctx,
		sessionID,
		r.offset,
		data,
		r.close,
		maxAttempts,
	); err != nil {
		return err
	}
	tracker.add(r.length)

	return nil
}

func (u *Uploader) appendUploadWithRetry(
	ctx context.Context,
	sessionID string,
	offset int64,
	data []byte,
	close bool,
	maxAttempts int,
) error {
	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		arg := files.NewUploadSessionAppendArg(
			files.NewUploadSessionCursor(sessionID, uint64(offset)),
		)
		arg.Close = close
		if err := u.client.UploadSessionAppendV2Context(
			ctx,
			arg,
			bytes.NewReader(data),
		); err != nil {
			if correctOffset, ok := uploadAppendCorrectOffset(err); ok {
				expectedOffset := offset + int64(len(data))
				switch correctOffset {
				case expectedOffset:
					return nil
				case offset:
					lastErr = err
					if err := waitForRetry(ctx, attempt, maxAttempts); err != nil {
						return err
					}
					continue
				default:
					return fmt.Errorf(
						"upload session offset mismatch: got %d, expected %d or %d",
						correctOffset,
						offset,
						expectedOffset,
					)
				}
			}
			if !isRetryableTransferError(err) {
				return err
			}
			lastErr = err
			if err := waitForRetry(ctx, attempt, maxAttempts); err != nil {
				return err
			}
			continue
		}
		return nil
	}
	if lastErr == nil {
		lastErr = errors.New("upload append failed")
	}
	return lastErr
}

func (u *Uploader) finishUploadWithRetry(
	ctx context.Context,
	sessionID string,
	offset int64,
	commit *files.CommitInfo,
	data []byte,
	maxAttempts int,
) (*files.FileMetadata, error) {
	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		arg := files.NewUploadSessionFinishArg(
			files.NewUploadSessionCursor(sessionID, uint64(offset)),
			commit,
		)
		metadata, err := u.client.UploadSessionFinishContext(
			ctx,
			arg,
			bytes.NewReader(data),
		)
		if err != nil {
			if correctOffset, ok := uploadFinishCorrectOffset(err); ok {
				expectedOffset := offset + int64(len(data))
				switch correctOffset {
				case expectedOffset:
					offset = correctOffset
					data = nil
				case offset:
				default:
					return nil, fmt.Errorf(
						"upload session offset mismatch: got %d, expected %d or %d",
						correctOffset,
						offset,
						expectedOffset,
					)
				}
				lastErr = err
				if err := waitForRetry(ctx, attempt, maxAttempts); err != nil {
					return nil, err
				}
				continue
			}
			if !isRetryableTransferError(err) {
				return nil, err
			}
			lastErr = err
			if err := waitForRetry(ctx, attempt, maxAttempts); err != nil {
				return nil, err
			}
			continue
		}
		return metadata, nil
	}
	if lastErr == nil {
		lastErr = errors.New("upload finish failed")
	}
	return nil, lastErr
}

func uploadAppendCorrectOffset(err error) (int64, bool) {
	var valueErr files.UploadSessionAppendV2APIError
	if !errors.As(err, &valueErr) {
		var pointerErr *files.UploadSessionAppendV2APIError
		if !errors.As(err, &pointerErr) || pointerErr == nil {
			return 0, false
		}
		valueErr = *pointerErr
	}
	if valueErr.EndpointError == nil ||
		valueErr.EndpointError.IncorrectOffset == nil {
		return 0, false
	}
	if valueErr.EndpointError.IncorrectOffset.CorrectOffset > math.MaxInt64 {
		return 0, false
	}
	return int64(valueErr.EndpointError.IncorrectOffset.CorrectOffset), true
}

func uploadFinishCorrectOffset(err error) (int64, bool) {
	var valueErr files.UploadSessionFinishAPIError
	if !errors.As(err, &valueErr) {
		var pointerErr *files.UploadSessionFinishAPIError
		if !errors.As(err, &pointerErr) || pointerErr == nil {
			return 0, false
		}
		valueErr = *pointerErr
	}
	if valueErr.EndpointError == nil ||
		valueErr.EndpointError.LookupFailed == nil ||
		valueErr.EndpointError.LookupFailed.IncorrectOffset == nil {
		return 0, false
	}
	if valueErr.EndpointError.LookupFailed.IncorrectOffset.CorrectOffset > math.MaxInt64 {
		return 0, false
	}
	return int64(valueErr.EndpointError.LookupFailed.IncorrectOffset.CorrectOffset), true
}

func readUploadChunk(reader io.Reader, limit int64) ([]byte, error) {
	if limit <= 0 {
		return nil, errors.New("upload chunk size must be positive")
	}

	buffer := make([]byte, int(limit))
	n, err := io.ReadFull(reader, buffer)
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return buffer[:n], io.EOF
		}
		return nil, err
	}

	return buffer, nil
}

type uploadByteRange struct {
	offset int64
	length int64
	close  bool
}

func splitUploadRanges(size int64) []uploadByteRange {
	if size == 0 {
		return nil
	}

	var ranges []uploadByteRange
	for offset := int64(0); offset < size; offset += uploadChunkSize {
		length := uploadChunkSize
		if remaining := size - offset; remaining < length {
			length = remaining
		}
		ranges = append(ranges, uploadByteRange{
			offset: offset,
			length: length,
		})
	}
	ranges[len(ranges)-1].close = true

	return ranges
}
