package filetransfer

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
)

func TestIsRetryableTransferError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil",
			err:  nil,
			want: false,
		},
		{
			name: "context canceled",
			err:  context.Canceled,
			want: false,
		},
		{
			name: "context deadline exceeded",
			err:  context.DeadlineExceeded,
			want: false,
		},
		{
			name: "unexpected EOF",
			err:  io.ErrUnexpectedEOF,
			want: true,
		},
		{
			name: "no progress",
			err:  io.ErrNoProgress,
			want: true,
		},
		{
			name: "network timeout",
			err:  timeoutError{},
			want: true,
		},
		{
			name: "HTTP 408",
			err: dropbox.SDKInternalError{
				StatusCode: 408,
			},
			want: true,
		},
		{
			name: "HTTP 429",
			err: dropbox.SDKInternalError{
				StatusCode: 429,
			},
			want: true,
		},
		{
			name: "HTTP 500",
			err: dropbox.SDKInternalError{
				StatusCode: 500,
			},
			want: true,
		},
		{
			name: "HTTP 599",
			err: dropbox.SDKInternalError{
				StatusCode: 599,
			},
			want: true,
		},
		{
			name: "HTTP 400",
			err: dropbox.SDKInternalError{
				StatusCode: 400,
			},
			want: false,
		},
		{
			name: "generic error",
			err:  errors.New("failed"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRetryableTransferError(tt.err); got != tt.want {
				t.Fatalf(
					"isRetryableTransferError() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestWaitForRetryLastAttempt(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := waitForRetry(ctx, 2, 3); err != nil {
		t.Fatalf("waitForRetry() error = %v", err)
	}
}

func TestWaitForRetryCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := waitForRetry(ctx, 0, 3)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf(
			"waitForRetry() error = %v, want context.Canceled",
			err,
		)
	}
}

type timeoutError struct{}

func (timeoutError) Error() string {
	return "timeout"
}

func (timeoutError) Timeout() bool {
	return true
}

func (timeoutError) Temporary() bool {
	return false
}
