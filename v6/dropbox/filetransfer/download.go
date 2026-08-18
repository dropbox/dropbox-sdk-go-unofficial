package filetransfer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
	"sync"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

// DownloadOptions configures a download operation.
//
// The zero value uses default retry and concurrency settings.
type DownloadOptions struct {
	// MaxAttempts is the maximum number of attempts for retryable failures.
	// A non-positive value uses the default.
	MaxAttempts int

	// ParallelDownloads is the maximum number of concurrent ranged downloads.
	// Values less than 2 disable parallel downloads.
	ParallelDownloads int

	// Progress receives download progress updates.
	Progress DownloadProgressFunc
}

// DownloadResult describes a completed download.
type DownloadResult struct {
	Metadata *files.FileMetadata
}

// DownloadClient is the subset of files.Client required by Downloader.
type DownloadClient interface {
	DownloadContext(
		ctx context.Context,
		arg *files.DownloadArg,
	) (
		metadata *files.FileMetadata,
		content io.ReadCloser,
		err error,
	)
}

// Downloader downloads files from Dropbox.
type Downloader struct {
	client DownloadClient
}

// NewDownloader creates a Downloader.
func NewDownloader(client DownloadClient) *Downloader {
	return &Downloader{
		client: client,
	}
}

// Download downloads remotePath into target.
//
// Download may retry transient failures according to options. Retries continue
// from the first uncommitted byte within the current operation.
//
// Before calling Commit, Download verifies that every expected byte range has
// been written exactly once for progress-accounting purposes.
//
// After Prepare succeeds, Download calls Commit at most once. If the download
// or Commit fails, Download calls Abort before returning. Built-in targets
// validate the downloaded size and content hash when available.
//
// Progress updates are monotonic and do not count bytes more than once.
func (d *Downloader) Download(
	ctx context.Context,
	remotePath string,
	target DownloadTarget,
	options DownloadOptions,
) (*DownloadResult, error) {
	if d == nil || d.client == nil {
		return nil, errors.New("download client is required")
	}
	if remotePath == "" {
		return nil, errors.New("download path is required")
	}
	if target == nil {
		return nil, errors.New("download target is required")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	maxAttempts := options.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = defaultMaxAttempts
	}

	if options.ParallelDownloads > 1 {
		return d.downloadWithParallelFallback(
			ctx,
			remotePath,
			target,
			maxAttempts,
			options.ParallelDownloads,
			options.Progress,
		)
	}

	return d.downloadSequential(
		ctx,
		remotePath,
		target,
		maxAttempts,
		options.Progress,
	)
}

func (d *Downloader) downloadWithParallelFallback(
	ctx context.Context,
	remotePath string,
	target DownloadTarget,
	maxAttempts int,
	parallelDownloads int,
	progress DownloadProgressFunc,
) (*DownloadResult, error) {
	metadata, info, firstWritten, err := d.prepareParallelDownload(
		ctx,
		remotePath,
		target,
		maxAttempts,
	)
	if isUnsatisfiableInitialRange(err) {
		return d.downloadSequential(
			ctx,
			remotePath,
			target,
			maxAttempts,
			progress,
		)
	}
	if err != nil {
		return nil, err
	}
	return d.downloadPreparedParallel(
		ctx,
		remotePath,
		target,
		metadata,
		info,
		firstWritten,
		maxAttempts,
		parallelDownloads,
		progress,
	)
}

func (d *Downloader) downloadSequential(
	ctx context.Context,
	remotePath string,
	target DownloadTarget,
	maxAttempts int,
	progress DownloadProgressFunc,
) (*DownloadResult, error) {
	var (
		metadata  *files.FileMetadata
		info      DownloadInfo
		prepared  bool
		committed int64
		tracker   *downloadProgressTracker
		lastErr   error
	)

	cleanupCtx := context.WithoutCancel(ctx)
	fail := func(err error) (*DownloadResult, error) {
		if prepared {
			if abortErr := target.Abort(cleanupCtx, err); abortErr != nil {
				err = errors.Join(err, fmt.Errorf("abort download target: %w", abortErr))
			}
		}
		return nil, err
	}

	for attempt := range maxAttempts {
		if err := ctx.Err(); err != nil {
			return fail(err)
		}

		arg := files.NewDownloadArg(remotePath)
		if committed > 0 {
			if err := files.SetRange(arg, committed); err != nil {
				return fail(err)
			}
		}

		responseMetadata, body, err := d.client.DownloadContext(ctx, arg)
		if err != nil {
			if !isRetryableTransferError(err) {
				return fail(err)
			}
			lastErr = err
			if err := waitForRetry(ctx, attempt, maxAttempts, err); err != nil {
				return fail(err)
			}
			continue
		}
		if body == nil {
			lastErr = errors.New("download response body is nil")
			if err := waitForRetry(ctx, attempt, maxAttempts, lastErr); err != nil {
				return fail(err)
			}
			continue
		}

		if !prepared {
			metadata, info, err = downloadMetadata(responseMetadata)
			if err != nil {
				_ = body.Close()
				return nil, err
			}
			if err := target.Prepare(ctx, info); err != nil {
				_ = body.Close()
				return nil, err
			}
			tracker = newDownloadProgressTracker(info.Size, progress)
			prepared = true
		} else if err := validateDownloadMetadata(metadata, responseMetadata); err != nil {
			_ = body.Close()
			return fail(err)
		}

		remaining := info.Size - committed
		if remaining < 0 {
			_ = body.Close()
			return fail(fmt.Errorf(
				"download exceeded expected size: got at least %d bytes, expected %d",
				committed,
				info.Size,
			))
		}

		written, copyErr, retryable := copyDownloadRange(
			body,
			target,
			committed,
			remaining,
			tracker,
		)
		closeErr := body.Close()
		committed += written

		if copyErr == nil && closeErr != nil && committed < info.Size {
			copyErr = closeErr
			retryable = true
		}
		if copyErr != nil {
			lastErr = copyErr
			if !retryable {
				return fail(copyErr)
			}
			if err := waitForRetry(ctx, attempt, maxAttempts, copyErr); err != nil {
				return fail(err)
			}
			continue
		}
		if committed != info.Size {
			lastErr = fmt.Errorf(
				"incomplete download: got %d bytes, expected %d",
				committed,
				info.Size,
			)
			if err := waitForRetry(ctx, attempt, maxAttempts, lastErr); err != nil {
				return fail(err)
			}
			continue
		}

		if err := target.Commit(ctx); err != nil {
			return fail(err)
		}
		return &DownloadResult{Metadata: metadata}, nil
	}

	if err := ctx.Err(); err != nil {
		return fail(err)
	}
	if lastErr == nil {
		lastErr = errors.New("download failed")
	}
	return fail(lastErr)
}

func (d *Downloader) downloadPreparedParallel(
	ctx context.Context,
	remotePath string,
	target DownloadTarget,
	metadata *files.FileMetadata,
	info DownloadInfo,
	firstWritten int64,
	maxAttempts int,
	parallelDownloads int,
	progress DownloadProgressFunc,
) (*DownloadResult, error) {
	prepared := true
	cleanupCtx := context.WithoutCancel(ctx)
	fail := func(err error) (*DownloadResult, error) {
		if prepared {
			if abortErr := target.Abort(cleanupCtx, err); abortErr != nil {
				err = errors.Join(err, fmt.Errorf("abort download target: %w", abortErr))
			}
		}
		return nil, err
	}

	tracker := newDownloadProgressTracker(info.Size, progress)
	tracker.add(firstWritten)

	if info.Size > firstWritten {
		ranges := splitDownloadRanges(
			firstWritten,
			info.Size-firstWritten,
			parallelDownloads,
		)
		workerCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		errCh := make(chan error, len(ranges))
		var wg sync.WaitGroup
		for _, r := range ranges {
			r := r
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := d.downloadRange(
					workerCtx,
					remotePath,
					target,
					r,
					metadata,
					maxAttempts,
					tracker,
				)
				if err != nil {
					cancel()
				}
				errCh <- err
			}()
		}

		wg.Wait()
		close(errCh)
		for rangeErr := range errCh {
			if rangeErr != nil {
				return fail(rangeErr)
			}
		}
	}

	if tracker.committedBytes() != info.Size {
		return fail(fmt.Errorf(
			"incomplete download: committed %d of %d bytes",
			tracker.committedBytes(),
			info.Size,
		))
	}
	if err := target.Commit(ctx); err != nil {
		return fail(err)
	}
	prepared = false

	return &DownloadResult{Metadata: metadata}, nil
}

func (d *Downloader) prepareParallelDownload(
	ctx context.Context,
	remotePath string,
	target DownloadTarget,
	maxAttempts int,
) (*files.FileMetadata, DownloadInfo, int64, error) {
	var lastErr error

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, DownloadInfo{}, 0, err
		}

		arg := files.NewDownloadArg(remotePath)
		if err := files.SetRangeLength(arg, 0, 1); err != nil {
			return nil, DownloadInfo{}, 0, err
		}

		metadata, body, err := d.client.DownloadContext(ctx, arg)
		if err != nil {
			if !isRetryableTransferError(err) {
				return nil, DownloadInfo{}, 0, err
			}
			lastErr = err
			if err := waitForRetry(ctx, attempt, maxAttempts, err); err != nil {
				return nil, DownloadInfo{}, 0, err
			}
			continue
		}
		if body == nil {
			lastErr = errors.New("download response body is nil")
			if err := waitForRetry(ctx, attempt, maxAttempts, lastErr); err != nil {
				return nil, DownloadInfo{}, 0, err
			}
			continue
		}

		stableMetadata, info, err := downloadMetadata(metadata)
		if err != nil {
			_ = body.Close()
			return nil, DownloadInfo{}, 0, err
		}
		if err := target.Prepare(ctx, info); err != nil {
			_ = body.Close()
			return nil, DownloadInfo{}, 0, err
		}

		expected := int64(0)
		if info.Size > 0 {
			expected = 1
		}
		tracker := newDownloadProgressTracker(info.Size, nil)
		written, copyErr, retryable := copyDownloadRange(
			body,
			target,
			0,
			expected,
			tracker,
		)
		_ = body.Close()
		if copyErr == nil && written == expected {
			return stableMetadata, info, written, nil
		}

		if abortErr := target.Abort(context.WithoutCancel(ctx), copyErr); abortErr != nil {
			copyErr = errors.Join(copyErr, abortErr)
		}
		if copyErr == nil {
			copyErr = fmt.Errorf(
				"incomplete initial range: got %d bytes, expected %d",
				written,
				expected,
			)
		}
		if !retryable {
			return nil, DownloadInfo{}, 0, copyErr
		}
		lastErr = copyErr
		if err := waitForRetry(ctx, attempt, maxAttempts, copyErr); err != nil {
			return nil, DownloadInfo{}, 0, err
		}
	}

	if err := ctx.Err(); err != nil {
		return nil, DownloadInfo{}, 0, err
	}
	if lastErr == nil {
		lastErr = errors.New("download failed")
	}
	return nil, DownloadInfo{}, 0, lastErr
}

func (d *Downloader) downloadRange(
	ctx context.Context,
	remotePath string,
	target DownloadTarget,
	r downloadByteRange,
	metadata *files.FileMetadata,
	maxAttempts int,
	tracker *downloadProgressTracker,
) error {
	var (
		committed int64
		lastErr   error
	)

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		remaining := r.length - committed
		if remaining == 0 {
			return nil
		}

		arg := files.NewDownloadArg(remotePath)
		if err := files.SetRangeLength(arg, r.offset+committed, remaining); err != nil {
			return err
		}

		responseMetadata, body, err := d.client.DownloadContext(ctx, arg)
		if err != nil {
			if !isRetryableTransferError(err) {
				return err
			}
			lastErr = err
			if err := waitForRetry(ctx, attempt, maxAttempts, err); err != nil {
				return err
			}
			continue
		}
		if body == nil {
			lastErr = errors.New("download response body is nil")
			if err := waitForRetry(ctx, attempt, maxAttempts, lastErr); err != nil {
				return err
			}
			continue
		}
		if err := validateDownloadMetadata(metadata, responseMetadata); err != nil {
			_ = body.Close()
			return err
		}

		written, copyErr, retryable := copyDownloadRange(
			body,
			target,
			r.offset+committed,
			remaining,
			tracker,
		)
		_ = body.Close()
		committed += written
		if copyErr == nil && committed == r.length {
			return nil
		}
		if copyErr == nil {
			copyErr = fmt.Errorf(
				"incomplete range at offset %d: got %d bytes, expected %d",
				r.offset,
				committed,
				r.length,
			)
		}
		if !retryable {
			return copyErr
		}
		lastErr = copyErr
		if err := waitForRetry(ctx, attempt, maxAttempts, copyErr); err != nil {
			return err
		}
	}

	if err := ctx.Err(); err != nil {
		return err
	}
	if lastErr == nil {
		lastErr = errors.New("download range failed")
	}
	return lastErr
}

type downloadByteRange struct {
	offset int64
	length int64
}

func splitDownloadRanges(
	offset int64,
	length int64,
	parts int,
) []downloadByteRange {
	if length <= 0 {
		return nil
	}
	if parts <= 1 {
		return []downloadByteRange{{offset: offset, length: length}}
	}
	if int64(parts) > length {
		parts = int(length)
	}

	ranges := make([]downloadByteRange, 0, parts)
	partSize := length / int64(parts)
	remainder := length % int64(parts)
	for i := 0; i < parts; i++ {
		size := partSize
		if int64(i) < remainder {
			size++
		}
		ranges = append(ranges, downloadByteRange{
			offset: offset,
			length: size,
		})
		offset += size
	}

	return ranges
}

func downloadMetadata(
	metadata *files.FileMetadata,
) (*files.FileMetadata, DownloadInfo, error) {
	if metadata == nil {
		return nil, DownloadInfo{}, errors.New("download metadata is nil")
	}
	if metadata.Size > math.MaxInt64 {
		return nil, DownloadInfo{}, fmt.Errorf(
			"download is too large: %d bytes",
			metadata.Size,
		)
	}

	return metadata, DownloadInfo{
		Size:        int64(metadata.Size),
		ContentHash: metadata.ContentHash,
	}, nil
}

func validateDownloadMetadata(
	expected *files.FileMetadata,
	actual *files.FileMetadata,
) error {
	if actual == nil {
		return errors.New("download metadata is nil")
	}
	if expected == nil {
		return nil
	}
	if expected.Rev != "" && actual.Rev != "" && actual.Rev != expected.Rev {
		return fmt.Errorf(
			"remote file changed during download: got rev %q, expected %q",
			actual.Rev,
			expected.Rev,
		)
	}
	if actual.Size != expected.Size {
		return fmt.Errorf(
			"remote file size changed during download: got %d, expected %d",
			actual.Size,
			expected.Size,
		)
	}
	if expected.ContentHash != "" &&
		actual.ContentHash != "" &&
		actual.ContentHash != expected.ContentHash {
		return fmt.Errorf(
			"remote file content hash changed during download: got %q, expected %q",
			actual.ContentHash,
			expected.ContentHash,
		)
	}

	return nil
}

func isUnsatisfiableInitialRange(err error) bool {
	var downloadErr files.DownloadAPIError
	if !errors.As(err, &downloadErr) {
		var pointerErr *files.DownloadAPIError
		if !errors.As(err, &pointerErr) || pointerErr == nil {
			return false
		}
		downloadErr = *pointerErr
	}
	return strings.HasPrefix(downloadErr.ErrorSummary, "range/not_satisfiable")
}

func copyDownloadRange(
	reader io.Reader,
	target DownloadTarget,
	offset int64,
	length int64,
	progress *downloadProgressTracker,
) (written int64, err error, retryable bool) {
	if length < 0 {
		return 0, errors.New("download range length must not be negative"), false
	}

	buffer := make([]byte, 32*1024)
	for written < length {
		remaining := length - written
		readBuffer := buffer
		if int64(len(readBuffer)) > remaining {
			readBuffer = readBuffer[:remaining]
		}

		n, readErr := reader.Read(readBuffer)
		if n > 0 {
			writeOffset := offset + written
			writtenNow := 0
			for writtenNow < n {
				m, writeErr := target.WriteAt(
					readBuffer[writtenNow:n],
					writeOffset+int64(writtenNow),
				)
				if m > 0 {
					writtenNow += m
					written += int64(m)
					progress.add(int64(m))
				}
				if writeErr != nil {
					return written, writeErr, false
				}
				if m == 0 {
					return written, io.ErrShortWrite, false
				}
			}
		}

		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				if written == length {
					break
				}
				return written, io.ErrUnexpectedEOF, true
			}
			return written, readErr, true
		}
		if n == 0 {
			return written, io.ErrNoProgress, true
		}
	}

	var extra [1]byte
	n, readErr := reader.Read(extra[:])
	if n > 0 {
		return written, errors.New("download response exceeded requested range"), false
	}
	if readErr != nil && !errors.Is(readErr, io.EOF) {
		return written, nil, false
	}

	return written, nil, false
}
