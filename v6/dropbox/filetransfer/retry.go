package filetransfer

import (
	"context"
	"errors"
	"io"
	"math"
	"math/rand/v2"
	"net"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
)

const (
	defaultMaxAttempts = 3
	retryBaseDelay     = 200 * time.Millisecond
	retryMaxDelay      = 5 * time.Second
)

func waitForRetry(
	ctx context.Context,
	attempt int,
	maxAttempts int,
	retryErr error,
) error {
	if attempt+1 >= maxAttempts {
		return nil
	}

	delay := retryDelay(retryErr, attempt)
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func retryDelay(err error, attempt int) time.Duration {
	if rateLimitErr := rateLimitTransferError(err); rateLimitErr != nil &&
		rateLimitErr.RateLimitError != nil {
		retryAfter := rateLimitErr.RateLimitError.RetryAfter
		if retryAfter > uint64(math.MaxInt64/int64(time.Second)) {
			return time.Duration(math.MaxInt64)
		}
		return time.Duration(retryAfter) * time.Second
	}

	delay := retryBaseDelay
	for i := 0; i < attempt && delay < retryMaxDelay; i++ {
		delay *= 2
		if delay > retryMaxDelay {
			delay = retryMaxDelay
		}
	}
	// Jitter between 50% and 100% of the exponential delay.
	half := delay / 2
	delay = half + time.Duration(rand.Int64N(int64(half)+1))

	return delay
}

func rateLimitTransferError(err error) *auth.RateLimitAPIError {
	var rateLimitErr auth.RateLimitAPIError
	if errors.As(err, &rateLimitErr) {
		return &rateLimitErr
	}

	var rateLimitErrPtr *auth.RateLimitAPIError
	if errors.As(err, &rateLimitErrPtr) {
		return rateLimitErrPtr
	}

	return nil
}

func isRetryableTransferError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	if errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, io.ErrNoProgress) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	var serverErr auth.ServerError
	if errors.As(err, &serverErr) {
		return true
	}
	if rateLimitTransferError(err) != nil {
		return true
	}
	var internalErr dropbox.SDKInternalError
	if errors.As(err, &internalErr) {
		return internalErr.StatusCode == 408 ||
			internalErr.StatusCode == 429 ||
			(internalErr.StatusCode >= 500 && internalErr.StatusCode <= 599)
	}

	return false
}
