// Copyright (c) Dropbox, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package retry_test

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/retry"
)

func TestPolicyNormalizedAppliesDefaults(t *testing.T) {
	// A policy with unset backoff fields should be normalized to the defaults.
	p := retry.Policy{MaxRetries: -5}.Normalized()

	if p.MaxRetries != 0 {
		t.Fatalf("negative MaxRetries = %d, want clamped to 0", p.MaxRetries)
	}
	if p.InitialBackoff != 500*time.Millisecond {
		t.Fatalf("InitialBackoff = %v, want default 500ms", p.InitialBackoff)
	}
	if p.MaxBackoff != 30*time.Second {
		t.Fatalf("MaxBackoff = %v, want default 30s", p.MaxBackoff)
	}
	if p.MaxRetryAfter != 30*time.Second {
		t.Fatalf("MaxRetryAfter = %v, want default 30s", p.MaxRetryAfter)
	}
}

func TestCanRetry(t *testing.T) {
	p := retry.Policy{MaxRetries: 2}
	for attempt, want := range map[int]bool{0: true, 1: true, 2: false, 3: false} {
		if got := p.CanRetry(attempt); got != want {
			t.Errorf("CanRetry(%d) = %v, want %v", attempt, got, want)
		}
	}
}

func TestBackoffIsExponentialAndCapped(t *testing.T) {
	p := retry.Policy{InitialBackoff: time.Second, MaxBackoff: 5 * time.Second}
	cases := []struct {
		attempt int
		want    time.Duration
	}{
		{0, 1 * time.Second},
		{1, 2 * time.Second},
		{2, 4 * time.Second},
		{3, 5 * time.Second}, // 8s capped to 5s
		{10, 5 * time.Second},
	}
	for _, c := range cases {
		if got := p.Backoff(c.attempt); got != c.want {
			t.Errorf("Backoff(%d) = %v, want %v", c.attempt, got, c.want)
		}
	}
}

func TestBackoffCapsAfterDoublingToMax(t *testing.T) {
	p := retry.Policy{InitialBackoff: 3 * time.Second, MaxBackoff: 6 * time.Second}
	if got := p.Backoff(1); got != 6*time.Second {
		t.Fatalf("Backoff(1) = %v, want capped 6s", got)
	}
}

func TestBackoffSaturatesBeforeDurationOverflow(t *testing.T) {
	p := retry.Policy{
		InitialBackoff: time.Duration(math.MaxInt64/2 + 1),
		MaxBackoff:     time.Duration(math.MaxInt64),
	}
	if got := p.Backoff(1); got != time.Duration(math.MaxInt64) {
		t.Fatalf("Backoff(1) = %v, want saturated MaxBackoff", got)
	}
}

func TestDelayRetryableStatuses(t *testing.T) {
	p := retry.Policy{MaxRetries: 3, InitialBackoff: time.Second, MaxBackoff: time.Second, MaxRetryAfter: time.Minute}
	retryable := []int{
		http.StatusTooManyRequests,
		http.StatusServiceUnavailable,
	}
	for _, code := range retryable {
		resp := &http.Response{StatusCode: code, Header: http.Header{}}
		if _, ok := p.Delay(resp, nil, 0); !ok {
			t.Errorf("status %d should be retryable", code)
		}
	}

	notRetryable := []int{
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusRequestTimeout,
		http.StatusConflict, // 409 needs a configured tag
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusGatewayTimeout,
	}
	for _, code := range notRetryable {
		resp := &http.Response{StatusCode: code, Header: http.Header{}}
		if _, ok := p.Delay(resp, nil, 0); ok {
			t.Errorf("status %d should not be retryable", code)
		}
	}
}

func TestDelayStopsAtMaxRetries(t *testing.T) {
	p := retry.Policy{MaxRetries: 1, InitialBackoff: time.Second, MaxBackoff: time.Second, MaxRetryAfter: time.Minute}
	resp := &http.Response{StatusCode: http.StatusServiceUnavailable, Header: http.Header{}}

	if _, ok := p.Delay(resp, nil, 0); !ok {
		t.Fatal("attempt 0 should be retryable")
	}
	if _, ok := p.Delay(resp, nil, 1); ok {
		t.Fatal("attempt 1 should exhaust MaxRetries")
	}
}

func TestDelayHonorsRetryAfterSeconds(t *testing.T) {
	p := retry.Policy{MaxRetries: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond, MaxRetryAfter: time.Minute}
	resp := &http.Response{StatusCode: http.StatusServiceUnavailable, Header: http.Header{"Retry-After": {"5"}}}

	wait, ok := p.Delay(resp, nil, 0)
	if !ok {
		t.Fatal("expected retry")
	}
	if wait != 5*time.Second {
		t.Fatalf("wait = %v, want 5s from Retry-After header", wait)
	}
}

func TestDelayRetryAfterOverMaxIsNotRetried(t *testing.T) {
	p := retry.Policy{MaxRetries: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond, MaxRetryAfter: 10 * time.Second}
	resp := &http.Response{StatusCode: http.StatusServiceUnavailable, Header: http.Header{"Retry-After": {"60"}}}

	if _, ok := p.Delay(resp, nil, 0); ok {
		t.Fatal("Retry-After beyond MaxRetryAfter should not be retried")
	}
}

func TestDelayRetryAfterHTTPDate(t *testing.T) {
	p := retry.Policy{MaxRetries: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond, MaxRetryAfter: time.Hour}
	when := time.Now().Add(20 * time.Second).UTC().Format(http.TimeFormat)
	resp := &http.Response{StatusCode: http.StatusServiceUnavailable, Header: http.Header{"Retry-After": {when}}}

	wait, ok := p.Delay(resp, nil, 0)
	if !ok {
		t.Fatal("expected retry")
	}
	// Allow slack for the second-granularity HTTP date and test execution time.
	if wait <= 0 || wait > 21*time.Second {
		t.Fatalf("wait = %v, want roughly 20s", wait)
	}
}

func TestDelayRetryAfterPastDate(t *testing.T) {
	// An HTTP-date already in the past means "retry now": wait clamps to 0 but
	// the response is still retryable.
	p := retry.Policy{MaxRetries: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond, MaxRetryAfter: time.Hour}
	past := time.Now().Add(-time.Hour).UTC().Format(http.TimeFormat)
	resp := &http.Response{StatusCode: http.StatusServiceUnavailable, Header: http.Header{"Retry-After": {past}}}

	wait, ok := p.Delay(resp, nil, 0)
	if !ok {
		t.Fatal("past Retry-After date should still be retryable")
	}
	if wait != 0 {
		t.Fatalf("wait = %v, want 0 for a past date", wait)
	}
}

func TestDelayRetryAfterOverflowRejected(t *testing.T) {
	// A Retry-After value large enough to overflow time.Duration must not wrap to
	// a negative duration and slip past the MaxRetryAfter cap (which would cause a
	// tight retry loop). It should saturate and be rejected by the cap instead.
	p := retry.Policy{MaxRetries: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond, MaxRetryAfter: 30 * time.Second}

	// Header path (int64 seconds).
	huge := strconv.FormatInt(math.MaxInt64/int64(time.Second)+1, 10)
	resp := &http.Response{StatusCode: http.StatusServiceUnavailable, Header: http.Header{"Retry-After": {huge}}}
	if wait, ok := p.Delay(resp, nil, 0); ok {
		t.Fatalf("overflowing header Retry-After should be rejected, got ok=true wait=%v", wait)
	}

	tooLargeForInt64 := strconv.FormatUint(uint64(math.MaxInt64)+1, 10)
	resp.Header.Set("Retry-After", tooLargeForInt64)
	if wait, ok := p.Delay(resp, nil, 0); ok {
		t.Fatalf("uint64-sized header Retry-After should be rejected, got ok=true wait=%v", wait)
	}

	// Body path (uint64 seconds).
	body := []byte(`{"error":{"retry_after":18446744073}}`)
	resp2 := &http.Response{StatusCode: http.StatusTooManyRequests, Header: http.Header{}}
	if wait, ok := p.Delay(resp2, body, 0); ok {
		t.Fatalf("overflowing body retry_after should be rejected, got ok=true wait=%v", wait)
	}
}

func TestDelay429RetryAfterBody(t *testing.T) {
	p := retry.Policy{MaxRetries: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond, MaxRetryAfter: time.Minute}
	body := []byte(`{"error_summary":"too_many_requests/..","error":{"reason":{".tag":"too_many_requests"},"retry_after":7}}`)
	resp := &http.Response{StatusCode: http.StatusTooManyRequests, Header: http.Header{}}

	wait, ok := p.Delay(resp, body, 0)
	if !ok {
		t.Fatal("expected retry")
	}
	if wait != 7*time.Second {
		t.Fatalf("wait = %v, want 7s from body retry_after", wait)
	}
}

func TestDelay429UnparseableBodyFallsBackToBackoff(t *testing.T) {
	// A 429 with no Retry-After header and a body that carries no usable
	// retry_after must still be retried, using exponential backoff.
	p := retry.Policy{MaxRetries: 3, InitialBackoff: 3 * time.Second, MaxBackoff: time.Minute, MaxRetryAfter: time.Minute}
	resp := &http.Response{StatusCode: http.StatusTooManyRequests, Header: http.Header{}}

	wait, ok := p.Delay(resp, []byte("not json at all"), 0)
	if !ok {
		t.Fatal("429 should be retryable even with an unparseable body")
	}
	if wait != 3*time.Second {
		t.Fatalf("wait = %v, want backoff fallback of 3s", wait)
	}
}

func TestDelayNegativeRetryAfterClampsToZero(t *testing.T) {
	p := retry.Policy{MaxRetries: 3, InitialBackoff: time.Second, MaxBackoff: time.Second, MaxRetryAfter: time.Minute}
	resp := &http.Response{StatusCode: http.StatusServiceUnavailable, Header: http.Header{"Retry-After": {"-10"}}}

	wait, ok := p.Delay(resp, nil, 0)
	if !ok {
		t.Fatal("expected retry")
	}
	if wait != 0 {
		t.Fatalf("wait = %v, want 0 for a negative Retry-After", wait)
	}
}

func TestBackoffInitialExceedingMaxIsCappedAtFirstAttempt(t *testing.T) {
	// When InitialBackoff already exceeds MaxBackoff, even attempt 0 is capped.
	p := retry.Policy{InitialBackoff: 10 * time.Second, MaxBackoff: 2 * time.Second}
	if got := p.Backoff(0); got != 2*time.Second {
		t.Fatalf("Backoff(0) = %v, want capped 2s", got)
	}
}

func TestDelayFallsBackToBackoff(t *testing.T) {
	p := retry.Policy{MaxRetries: 3, InitialBackoff: 2 * time.Second, MaxBackoff: time.Minute, MaxRetryAfter: time.Minute}
	resp := &http.Response{StatusCode: http.StatusServiceUnavailable, Header: http.Header{}}

	wait, ok := p.Delay(resp, nil, 1)
	if !ok {
		t.Fatal("expected retry")
	}
	if wait != 4*time.Second { // 2s * 2^1
		t.Fatalf("wait = %v, want exponential backoff 4s", wait)
	}
}

func TestDelay409TagMatching(t *testing.T) {
	body := []byte(`{"error_summary":"too_many_write_operations/..","error":{".tag":"too_many_write_operations"}}`)
	resp := &http.Response{StatusCode: http.StatusConflict, Header: http.Header{}}

	// Configured tag: retried.
	p := retry.Policy{MaxRetries: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond, MaxRetryAfter: time.Minute, Retryable409Tags: []string{"too_many_write_operations"}}
	if _, ok := p.Delay(resp, body, 0); !ok {
		t.Fatal("configured 409 tag should be retried")
	}

	// Unconfigured tag: not retried.
	other := retry.Policy{MaxRetries: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond, MaxRetryAfter: time.Minute, Retryable409Tags: []string{"some_other_tag"}}
	if _, ok := other.Delay(resp, body, 0); ok {
		t.Fatal("unconfigured 409 tag should not be retried")
	}

	// No tags configured: not retried.
	none := retry.Policy{MaxRetries: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond, MaxRetryAfter: time.Minute}
	if _, ok := none.Delay(resp, body, 0); ok {
		t.Fatal("409 with no configured tags should not be retried")
	}
}

func TestDelay409TagFromErrorSummary(t *testing.T) {
	// Body carries no ".tag" but a slash-delimited error_summary.
	body := []byte(`{"error_summary":"too_many_write_operations/lock_conflict/..","error":{}}`)
	resp := &http.Response{StatusCode: http.StatusConflict, Header: http.Header{}}
	p := retry.Policy{MaxRetries: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond, MaxRetryAfter: time.Minute, Retryable409Tags: []string{"too_many_write_operations"}}

	if _, ok := p.Delay(resp, body, 0); !ok {
		t.Fatal("409 tag derived from error_summary should be retried")
	}
}

func TestDelay409TagPresentButNoMatch(t *testing.T) {
	// The response carries a valid Stone tag, but it isn't in Retryable409Tags.
	body := []byte(`{"error_summary":"some_other_error/..","error":{".tag":"some_other_error"}}`)
	resp := &http.Response{StatusCode: http.StatusConflict, Header: http.Header{}}
	p := retry.Policy{MaxRetries: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond, MaxRetryAfter: time.Minute, Retryable409Tags: []string{"too_many_write_operations"}}

	if _, ok := p.Delay(resp, body, 0); ok {
		t.Fatal("409 with a non-matching tag should not be retried")
	}
}

func TestDelay409MalformedBodyNotRetried(t *testing.T) {
	p := retry.Policy{MaxRetries: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond, MaxRetryAfter: time.Minute, Retryable409Tags: []string{"too_many_write_operations"}}

	cases := map[string][]byte{
		"invalid json":     []byte("}{ not json"),
		"empty summary":    []byte(`{"error_summary":"","error":{}}`),
		"summary no slash": []byte(`{"error_summary":"noslashvalue","error":{}}`),
	}
	for name, body := range cases {
		resp := &http.Response{StatusCode: http.StatusConflict, Header: http.Header{}}
		if _, ok := p.Delay(resp, body, 0); ok {
			t.Errorf("%s: 409 should not be retried, got ok=true", name)
		}
	}
}

func TestDelayUnparseableRetryAfterFallsBackToBackoff(t *testing.T) {
	// A Retry-After header that is neither an integer nor an HTTP date is ignored,
	// and the response falls back to exponential backoff.
	p := retry.Policy{MaxRetries: 3, InitialBackoff: 2 * time.Second, MaxBackoff: time.Minute, MaxRetryAfter: time.Minute}
	resp := &http.Response{StatusCode: http.StatusServiceUnavailable, Header: http.Header{"Retry-After": {"soon-ish"}}}

	wait, ok := p.Delay(resp, nil, 0)
	if !ok {
		t.Fatal("expected retry with backoff fallback")
	}
	if wait != 2*time.Second {
		t.Fatalf("wait = %v, want backoff fallback 2s", wait)
	}
}

func TestSleepReturnsAfterDuration(t *testing.T) {
	start := time.Now()
	if err := retry.Sleep(context.Background(), 10*time.Millisecond); err != nil {
		t.Fatalf("Sleep returned error: %v", err)
	}
	if elapsed := time.Since(start); elapsed < 10*time.Millisecond {
		t.Fatalf("Sleep returned after %v, want >= 10ms", elapsed)
	}
}

func TestSleepZeroReturnsImmediately(t *testing.T) {
	if err := retry.Sleep(context.Background(), 0); err != nil {
		t.Fatalf("Sleep(0) returned error: %v", err)
	}
	if err := retry.Sleep(context.Background(), -time.Second); err != nil {
		t.Fatalf("Sleep(negative) returned error: %v", err)
	}
}

func TestSleepCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := retry.Sleep(ctx, time.Hour)
	if err == nil {
		t.Fatal("expected context error from cancelled Sleep")
	}
	if err != context.Canceled {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
}

func TestSleepDeadlineDuringWait(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := retry.Sleep(ctx, time.Hour)
	if err != context.DeadlineExceeded {
		t.Fatalf("err = %v, want context.DeadlineExceeded", err)
	}
}
