package filetransfer

import "sync"

// UploadProgress reports upload progress.
type UploadProgress struct {
	// BytesCommitted is the number of unique bytes confirmed by Dropbox.
	// The value is monotonic.
	BytesCommitted int64

	// TotalBytes is the size of the upload, or -1 when the size is unknown.
	TotalBytes int64
}

// UploadProgressFunc receives upload progress updates.
type UploadProgressFunc func(UploadProgress)

type uploadProgressTracker struct {
	mu        sync.Mutex
	committed int64
	total     int64
	progress  UploadProgressFunc
}

func newUploadProgressTracker(
	total int64,
	progress UploadProgressFunc,
) *uploadProgressTracker {
	return &uploadProgressTracker{
		total:    total,
		progress: progress,
	}
}

func (p *uploadProgressTracker) add(n int64) {
	if n <= 0 {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.committed += n
	if p.progress != nil {
		p.progress(UploadProgress{
			BytesCommitted: p.committed,
			TotalBytes:     p.total,
		})
	}
}

func (p *uploadProgressTracker) committedBytes() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.committed
}
