package filetransfer

import "sync"

// DownloadProgress reports download progress.
type DownloadProgress struct {
	// BytesCommitted is the number of unique bytes successfully written to
	// the target. The value is monotonic.
	BytesCommitted int64

	// TotalBytes is the expected size of the download.
	TotalBytes int64
}

// DownloadProgressFunc receives download progress updates.
type DownloadProgressFunc func(DownloadProgress)

type downloadProgressTracker struct {
	mu        sync.Mutex
	committed int64
	total     int64
	progress  DownloadProgressFunc
}

func newDownloadProgressTracker(
	total int64,
	progress DownloadProgressFunc,
) *downloadProgressTracker {
	return &downloadProgressTracker{
		total:    total,
		progress: progress,
	}
}

func (p *downloadProgressTracker) add(n int64) {
	if n <= 0 {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.committed += n
	if p.progress != nil {
		p.progress(DownloadProgress{
			BytesCommitted: p.committed,
			TotalBytes:     p.total,
		})
	}
}

func (p *downloadProgressTracker) committedBytes() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.committed
}
