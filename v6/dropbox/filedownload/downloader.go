package filedownload

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sync"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/contenthash"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

// Client is the subset of files.Client required by Downloader.
type Client interface {
	DownloadContext(
		ctx context.Context,
		arg *files.DownloadArg,
	) (*files.FileMetadata, io.ReadCloser, error)
}

// Progress describes download progress.
type Progress struct {
	BytesWritten int64
	TotalBytes   int64
	ResumedFrom  int64
}

// ProgressFunc receives download progress updates.
type ProgressFunc func(Progress)

// Option configures a Downloader.
type Option func(*Downloader)

// WithMaxAttempts configures how many times DownloadFile retries.
func WithMaxAttempts(maxAttempts int) Option {
	return func(d *Downloader) {
		if maxAttempts > 0 {
			d.maxAttempts = maxAttempts
		}
	}
}

// WithProgress configures a callback for download progress updates.
func WithProgress(progress ProgressFunc) Option {
	return func(d *Downloader) {
		d.progress = progress
	}
}

// WithParallelDownloads configures how many ranged downloads to run at once
// for fresh downloads. Resumed downloads continue serially from the part file.
func WithParallelDownloads(parallelDownloads int) Option {
	return func(d *Downloader) {
		if parallelDownloads > 1 {
			d.parallelDownloads = parallelDownloads
		}
	}
}

// Downloader downloads Dropbox files to local storage.
type Downloader struct {
	client            Client
	maxAttempts       int
	parallelDownloads int
	progress          ProgressFunc
}

// New returns a Downloader that uses client.
func New(client Client, opts ...Option) *Downloader {
	d := &Downloader{
		client:            client,
		maxAttempts:       3,
		parallelDownloads: 1,
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// Result describes a completed file download.
type Result struct {
	Metadata    *files.FileMetadata
	ResumedFrom int64
}

func (d *Downloader) downloadFileAttempt(
	ctx context.Context,
	remotePath string,
	localPath string,
	expectedRev string,
) (*Result, error) {
	partPath := localPath + ".part"

	offset, err := partFileSize(partPath)
	if err != nil {
		return nil, err
	}

	if offset == 0 && d.parallelDownloads > 1 {
		return d.downloadFileParallel(ctx, remotePath, localPath, partPath, expectedRev)
	}

	f, err := os.OpenFile(
		partPath,
		os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return nil, err
	}

	arg := files.NewDownloadArg(remotePath)
	if offset > 0 {
		if err := files.SetRange(arg, offset); err != nil {
			_ = f.Close()
			return nil, err
		}
	}

	metadata, body, err := d.client.DownloadContext(ctx, arg)
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	if body == nil {
		_ = f.Close()
		return nil, fmt.Errorf("download response body is nil")
	}
	defer body.Close()

	if err := validateRevision(expectedRev, metadata); err != nil {
		_ = f.Close()
		_ = os.Remove(partPath)
		return nil, err
	}

	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		_ = f.Close()
		return nil, err
	}

	progress := d.newProgress(offset, metadataSize(metadata), offset)
	if _, err := io.Copy(f, progress.reader(body)); err != nil {
		_ = f.Close()
		return &Result{
			Metadata:    metadata,
			ResumedFrom: offset,
		}, err
	}

	if err := f.Close(); err != nil {
		return nil, err
	}

	if err := validatePartFile(partPath, metadata); err != nil {
		_ = os.Remove(partPath)
		return nil, err
	}

	if err := os.Rename(partPath, localPath); err != nil {
		return nil, err
	}

	return &Result{
		Metadata:    metadata,
		ResumedFrom: offset,
	}, nil
}

func (d *Downloader) downloadFileParallel(
	ctx context.Context,
	remotePath string,
	localPath string,
	partPath string,
	expectedRev string,
) (*Result, error) {
	arg := files.NewDownloadArg(remotePath)
	if err := files.SetRangeLength(arg, 0, 1); err != nil {
		return nil, err
	}

	metadata, body, err := d.client.DownloadContext(ctx, arg)
	if err != nil {
		return nil, err
	}
	if body == nil {
		return nil, fmt.Errorf("download response body is nil")
	}
	defer body.Close()

	if err := validateRevision(expectedRev, metadata); err != nil {
		_ = os.Remove(partPath)
		return nil, err
	}

	if metadata == nil {
		return nil, fmt.Errorf("download metadata is nil")
	}
	if metadata.Size > math.MaxInt64 {
		return nil, fmt.Errorf("download too large: %d bytes", metadata.Size)
	}

	total := int64(metadata.Size)
	f, err := os.OpenFile(
		partPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0o644,
	)
	if err != nil {
		return nil, err
	}
	if err := f.Truncate(total); err != nil {
		_ = f.Close()
		return nil, err
	}

	progress := d.newProgress(0, total, 0)
	if total > 0 {
		if _, err := io.Copy(&writeAtWriter{f: f}, progress.reader(body)); err != nil {
			_ = f.Close()
			return nil, err
		}
	}
	if err := f.Close(); err != nil {
		return nil, err
	}

	if total > 1 {
		if err := d.downloadParallelRanges(ctx, remotePath, partPath, metadata, progress); err != nil {
			_ = os.Remove(partPath)
			return nil, err
		}
	}

	if err := validatePartFile(partPath, metadata); err != nil {
		_ = os.Remove(partPath)
		return nil, err
	}

	if err := os.Rename(partPath, localPath); err != nil {
		return nil, err
	}

	return &Result{
		Metadata:    metadata,
		ResumedFrom: 0,
	}, nil
}

func (d *Downloader) downloadParallelRanges(
	ctx context.Context,
	remotePath string,
	partPath string,
	metadata *files.FileMetadata,
	progress *progressTracker,
) error {
	f, err := os.OpenFile(partPath, os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	total := int64(metadata.Size)
	ranges := splitRanges(1, total-1, d.parallelDownloads)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error, len(ranges))
	var wg sync.WaitGroup
	for _, r := range ranges {
		r := r
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := d.downloadRange(ctx, remotePath, f, r, metadata.Rev, progress)
			if err != nil {
				cancel()
			}
			errs <- err
		}()
	}

	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			cancel()
			return err
		}
	}
	return nil
}

func (d *Downloader) downloadRange(
	ctx context.Context,
	remotePath string,
	f *os.File,
	r byteRange,
	expectedRev string,
	progress *progressTracker,
) error {
	arg := files.NewDownloadArg(remotePath)
	if err := files.SetRangeLength(arg, r.offset, r.length); err != nil {
		return err
	}

	metadata, body, err := d.client.DownloadContext(ctx, arg)
	if err != nil {
		return err
	}
	if body == nil {
		return fmt.Errorf("download response body is nil")
	}
	defer body.Close()

	if err := validateRevision(expectedRev, metadata); err != nil {
		return err
	}

	_, err = io.Copy(&writeAtWriter{f: f, offset: r.offset}, progress.reader(body))
	return err
}

func (d *Downloader) DownloadFile(
	ctx context.Context,
	remotePath string,
	localPath string,
) (*Result, error) {
	var (
		lastErr     error
		expectedRev string
	)

	for attempt := 0; attempt < d.maxAttempts; attempt++ {
		result, err := d.downloadFileAttempt(
			ctx,
			remotePath,
			localPath,
			expectedRev,
		)

		if result != nil && result.Metadata != nil && expectedRev == "" {
			expectedRev = result.Metadata.Rev
		}

		if err == nil {
			return result, nil
		}

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		var changed *remoteFileChangedError
		if errors.As(err, &changed) {
			return nil, err
		}

		lastErr = err
	}

	return nil, lastErr
}

type progressTracker struct {
	mu          sync.Mutex
	written     int64
	total       int64
	resumedFrom int64
	progress    ProgressFunc
}

func (d *Downloader) newProgress(written int64, total int64, resumedFrom int64) *progressTracker {
	return &progressTracker{
		written:     written,
		total:       total,
		resumedFrom: resumedFrom,
		progress:    d.progress,
	}
}

func (p *progressTracker) reader(r io.Reader) io.Reader {
	return &progressReader{reader: r, progress: p}
}

func (p *progressTracker) add(n int) {
	if n <= 0 || p.progress == nil {
		return
	}

	p.mu.Lock()
	p.written += int64(n)
	progress := Progress{
		BytesWritten: p.written,
		TotalBytes:   p.total,
		ResumedFrom:  p.resumedFrom,
	}
	p.progress(progress)
	p.mu.Unlock()
}

type progressReader struct {
	reader   io.Reader
	progress *progressTracker
}

func (r *progressReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	r.progress.add(n)
	return n, err
}

type writeAtWriter struct {
	f      *os.File
	offset int64
}

func (w *writeAtWriter) Write(p []byte) (int, error) {
	n, err := w.f.WriteAt(p, w.offset)
	w.offset += int64(n)
	return n, err
}

type byteRange struct {
	offset int64
	length int64
}

func splitRanges(offset int64, length int64, parts int) []byteRange {
	if length <= 0 {
		return nil
	}
	if parts <= 1 {
		return []byteRange{{offset: offset, length: length}}
	}
	if int64(parts) > length {
		parts = int(length)
	}

	ranges := make([]byteRange, 0, parts)
	partSize := length / int64(parts)
	remainder := length % int64(parts)
	for i := 0; i < parts; i++ {
		size := partSize
		if int64(i) < remainder {
			size++
		}
		ranges = append(ranges, byteRange{offset: offset, length: size})
		offset += size
	}
	return ranges
}

func validateRevision(expectedRev string, metadata *files.FileMetadata) error {
	if expectedRev == "" || metadata == nil || metadata.Rev == "" {
		return nil
	}
	if metadata.Rev != expectedRev {
		return &remoteFileChangedError{got: metadata.Rev, expected: expectedRev}
	}
	return nil
}

type remoteFileChangedError struct {
	got      string
	expected string
}

func (e *remoteFileChangedError) Error() string {
	return fmt.Sprintf(
		"remote file changed during retry: got rev %q, expected %q",
		e.got,
		e.expected,
	)
}

func validatePartFile(partPath string, metadata *files.FileMetadata) error {
	if metadata == nil {
		return nil
	}

	stat, err := os.Stat(partPath)
	if err != nil {
		return err
	}

	if uint64(stat.Size()) != metadata.Size {
		return fmt.Errorf(
			"incomplete download: got %d bytes, expected %d",
			stat.Size(),
			metadata.Size,
		)
	}

	if metadata.ContentHash == "" {
		return nil
	}

	f, err := os.Open(partPath)
	if err != nil {
		return err
	}
	defer f.Close()

	hash, err := contenthash.Compute(f)
	if err != nil {
		return err
	}
	if hash != metadata.ContentHash {
		return fmt.Errorf(
			"content hash mismatch: got %q, expected %q",
			hash,
			metadata.ContentHash,
		)
	}
	return nil
}

func partFileSize(partPath string) (int64, error) {
	stat, err := os.Stat(partPath)
	if err == nil {
		return stat.Size(), nil
	}
	if os.IsNotExist(err) {
		return 0, nil
	}
	return 0, err
}

func metadataSize(metadata *files.FileMetadata) int64 {
	if metadata == nil {
		return 0
	}
	return int64(metadata.Size)
}
