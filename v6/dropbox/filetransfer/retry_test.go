package filetransfer

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
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
			name: "network non-timeout",
			err:  networkError{},
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

	if err := waitForRetry(ctx, 2, 3, nil); err != nil {
		t.Fatalf("waitForRetry() error = %v", err)
	}
}

func TestWaitForRetryCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := waitForRetry(ctx, 0, 3, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf(
			"waitForRetry() error = %v, want context.Canceled",
			err,
		)
	}
}

func TestRetryDelayUsesRateLimitRetryAfter(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "value",
			err: auth.RateLimitAPIError{
				RateLimitError: &auth.RateLimitError{RetryAfter: 7},
			},
		},
		{
			name: "pointer",
			err: &auth.RateLimitAPIError{
				RateLimitError: &auth.RateLimitError{RetryAfter: 7},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !isRetryableTransferError(tt.err) {
				t.Fatal("isRetryableTransferError() = false, want true")
			}
			if got, want := retryDelay(tt.err, 0), 7*time.Second; got != want {
				t.Fatalf("retryDelay() = %v, want %v", got, want)
			}
		})
	}
}

func TestRetryDelayUsesJitteredExponentialBackoff(t *testing.T) {
	for attempt, maximum := range map[int]time.Duration{
		0: retryBaseDelay,
		1: retryBaseDelay * 2,
		2: retryBaseDelay * 4,
	} {
		delay := retryDelay(errors.New("temporary failure"), attempt)
		if delay < maximum/2 || delay > maximum {
			t.Fatalf("retryDelay(attempt %d) = %v, want between %v and %v", attempt, delay, maximum/2, maximum)
		}
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

type networkError struct{}

func (networkError) Error() string   { return "network error" }
func (networkError) Timeout() bool   { return false }
func (networkError) Temporary() bool { return false }
