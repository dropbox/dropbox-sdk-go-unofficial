package filetransfer

import (
	"sync"
	"testing"
)

func TestUploadProgressTracker(t *testing.T) {
	var updates []UploadProgress

	tracker := newUploadProgressTracker(
		10,
		func(progress UploadProgress) {
			updates = append(updates, progress)
		},
	)

	tracker.add(4)
	tracker.add(0)
	tracker.add(-1)
	tracker.add(3)

	if got := tracker.committedBytes(); got != 7 {
		t.Fatalf("committedBytes() = %d, want 7", got)
	}

	want := []UploadProgress{
		{
			BytesCommitted: 4,
			TotalBytes:     10,
		},
		{
			BytesCommitted: 7,
			TotalBytes:     10,
		},
	}

	if len(updates) != len(want) {
		t.Fatalf("len(updates) = %d, want %d", len(updates), len(want))
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

func TestUploadProgressTrackerUnknownSize(t *testing.T) {
	var got UploadProgress

	tracker := newUploadProgressTracker(
		-1,
		func(progress UploadProgress) {
			got = progress
		},
	)

	tracker.add(5)

	want := UploadProgress{
		BytesCommitted: 5,
		TotalBytes:     -1,
	}

	if got != want {
		t.Fatalf("progress = %+v, want %+v", got, want)
	}
}

func TestUploadProgressTrackerWithoutCallback(t *testing.T) {
	tracker := newUploadProgressTracker(10, nil)

	tracker.add(4)

	if got := tracker.committedBytes(); got != 4 {
		t.Fatalf("committedBytes() = %d, want 4", got)
	}
}

func TestUploadProgressTrackerConcurrent(t *testing.T) {
	const (
		workers = 10
		perCall = 10
	)

	tracker := newUploadProgressTracker(workers*perCall, nil)

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
