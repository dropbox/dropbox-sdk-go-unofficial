package filetransfer

import (
	"sync"
	"testing"
)

func TestDownloadProgressTracker(t *testing.T) {
	var updates []DownloadProgress

	tracker := newDownloadProgressTracker(
		10,
		func(progress DownloadProgress) {
			updates = append(updates, progress)
		},
	)

	tracker.add(3)
	tracker.add(0)
	tracker.add(-1)
	tracker.add(2)

	if got := tracker.committedBytes(); got != 5 {
		t.Fatalf("committedBytes() = %d, want 5", got)
	}

	if len(updates) != 2 {
		t.Fatalf("len(updates) = %d, want 2", len(updates))
	}

	want := []DownloadProgress{
		{
			BytesCommitted: 3,
			TotalBytes:     10,
		},
		{
			BytesCommitted: 5,
			TotalBytes:     10,
		},
	}

	for i := range want {
		if updates[i] != want[i] {
			t.Fatalf(
				"updates[%d] = %+v, want %+v",
				i,
				updates[i],
				want[i],
			)
		}
	}
}

func TestDownloadProgressTrackerWithoutCallback(t *testing.T) {
	tracker := newDownloadProgressTracker(10, nil)

	tracker.add(4)

	if got := tracker.committedBytes(); got != 4 {
		t.Fatalf("committedBytes() = %d, want 4", got)
	}
}

func TestDownloadProgressTrackerConcurrent(t *testing.T) {
	const (
		workers = 10
		perCall = 10
	)

	tracker := newDownloadProgressTracker(workers*perCall, nil)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tracker.add(perCall)
		}()
	}

	wg.Wait()

	if got := tracker.committedBytes(); got != workers*perCall {
		t.Fatalf(
			"committedBytes() = %d, want %d",
			got,
			workers*perCall,
		)
	}
}
