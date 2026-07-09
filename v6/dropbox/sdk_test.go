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

package dropbox_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/retry"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/users"
	"golang.org/x/oauth2"
)

func generateURL(base string, namespace string, route string) string {
	return fmt.Sprintf("%s/%s/%s", base, namespace, route)
}

func retryTestPolicy(maxRetries int) retry.Policy {
	return retry.Policy{
		MaxRetries:     maxRetries,
		InitialBackoff: time.Nanosecond,
		MaxBackoff:     time.Nanosecond,
		MaxRetryAfter:  time.Second,
	}
}

func retryTestPolicyPtr(maxRetries int) *retry.Policy {
	policy := retryTestPolicy(maxRetries)
	return &policy
}

func retryTestPolicyPtrWith409Tags(maxRetries int, tags ...string) *retry.Policy {
	policy := retryTestPolicy(maxRetries)
	policy.Retryable409Tags = tags
	return &policy
}

type legacyUsersClient struct{}

var _ users.Client = legacyUsersClient{}
var _ func(dropbox.Config) users.Client = users.New
var _ func(dropbox.Config) users.ContextClient = users.NewContext

func (legacyUsersClient) FeaturesGetValues(arg *users.UserFeaturesGetValuesBatchArg) (*users.UserFeaturesGetValuesBatchResult, error) {
	return nil, nil
}

func (legacyUsersClient) GetAccount(arg *users.GetAccountArg) (*users.BasicAccount, error) {
	return nil, nil
}

func (legacyUsersClient) GetAccountBatch(arg *users.GetAccountBatchArg) ([]*users.BasicAccount, error) {
	return nil, nil
}

func (legacyUsersClient) GetCurrentAccount() (*users.FullAccount, error) {
	return nil, nil
}

func (legacyUsersClient) GetSpaceUsage() (*users.SpaceUsage, error) {
	return nil, nil
}

func TestInternalError(t *testing.T) {
	eString := "internal server error"
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, eString, http.StatusInternalServerError)
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogDebug,
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	client := users.NewContext(config)
	v, e := client.GetCurrentAccountContext(context.Background())
	if v != nil || strings.Trim(e.Error(), "\n") != eString {
		t.Errorf("v: %v e: '%s'\n", v, e.Error())
	}
}

func TestRateLimitJSON(t *testing.T) {
	eString := `{"error_summary": "too_many_requests/..", "error": {"reason": {".tag": "too_many_requests"}, "retry_after": 300}}`
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Retry-After", "10")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(eString))
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogDebug,
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	client := users.NewContext(config)
	_, e := client.GetCurrentAccountContext(context.Background())
	re, ok := e.(auth.RateLimitAPIError)
	if !ok {
		t.Errorf("Unexpected error type: %T\n", e)
	}
	if re.RateLimitError.RetryAfter != 300 {
		t.Errorf("Unexpected retry-after value: %d\n", re.RateLimitError.RetryAfter)
	}
	if re.RateLimitError.Reason.Tag != auth.RateLimitReasonTooManyRequests {
		t.Errorf("Unexpected reason: %v\n", re.RateLimitError.Reason)
	}
}

func TestAuthError(t *testing.T) {
	eString := `{"error_summary": "user_suspended/...", "error": {".tag": "user_suspended"}}`
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(eString))
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogDebug,
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	client := users.New(config)
	_, e := client.GetCurrentAccount()
	re, ok := e.(auth.AuthAPIError)
	if !ok {
		t.Errorf("Unexpected error type: %T\n", e)
	}
	fmt.Printf("ERROR is %v\n", re)
	if re.AuthError.Tag != auth.AuthErrorUserSuspended {
		t.Errorf("Unexpected tag: %s\n", re.AuthError.Tag)
	}
}

func TestAccessError(t *testing.T) {
	eString := `{"error_summary": "access_error/...",
	"error": {
		".tag": "paper_access_denied",
	  "paper_access_denied": {".tag": "not_paper_user"}
	}}`
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(eString))
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogDebug,
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	client := users.New(config)
	_, e := client.GetCurrentAccount()
	re, ok := e.(auth.AccessAPIError)
	if !ok {
		t.Errorf("Unexpected error type: %T\n", e)
	}
	if re.AccessError.Tag != auth.AccessErrorPaperAccessDenied {
		t.Errorf("Unexpected tag: %s\n", re.AccessError.Tag)
	}
	if re.AccessError.PaperAccessDenied.Tag != auth.PaperAccessErrorNotPaperUser {
		t.Errorf("Unexpected tag: %s\n", re.AccessError.PaperAccessDenied.Tag)
	}
}

func TestAppError(t *testing.T) {
	eString := `{"error_summary":"","error":{".tag":"app_id_mismatch"}}`
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(eString))
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogDebug,
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	client := auth.New(config)
	_, e := client.TokenFromOauth1(nil)
	re, ok := e.(auth.TokenFromOauth1APIError)
	if !ok {
		t.Errorf("Unexpected error type: %T\n%v\n", e, e)
	}
	if re.EndpointError.Tag != auth.TokenFromOAuth1ErrorAppIdMismatch {
		t.Errorf("Unexpected tag: %s\n", re.EndpointError.Tag)
	}
}

func TestHTTPHeaderSafeJSON(t *testing.T) {
	for _, test := range []struct {
		name string
		in   interface{}
		want string
	}{
		{
			name: "empty string",
			in:   ``,
			want: `""`,
		},
		{
			name: "integer",
			in:   123,
			want: `123`,
		},
		{
			name: "normal string",
			in:   `Normal string!`,
			want: `"Normal string!"`,
		},
		{
			name: "unicode",
			in:   `üñîcødé`,
			want: `"\u00fc\u00f1\u00eec\u00f8d\u00e9"`,
		},
		{
			name: "7f",
			in:   "\x7f",
			want: `"\u007f"`,
		},
		{
			name: "example from the docs",
			in: struct {
				Field string `json:"field"`
			}{
				Field: "some_üñîcødé_and_\x7F",
			},
			want: `{"field":"some_\u00fc\u00f1\u00eec\u00f8d\u00e9_and_\u007f"}`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			b, err := json.Marshal(test.in)
			if err != nil {
				t.Fatal(err)
			}
			got := dropbox.HTTPHeaderSafeJSON(b)
			if got != test.want {
				t.Errorf("Want %q got %q", test.want, got)
			}
		})
	}
}

func TestExecuteBackwardCompatibility(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{}`))
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogOff,
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	ctx := dropbox.NewContext(config)
	resp, respBody, err := ctx.Execute(dropbox.Request{
		Host:      "api",
		Namespace: "users",
		Route:     "get_current_account",
		Auth:      "user",
		Style:     "rpc",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(resp) != "{}" {
		t.Errorf("unexpected response: %s", string(resp))
	}
	if respBody != nil {
		t.Errorf("unexpected response body: %v", respBody)
	}
}

func TestExecuteDoesNotRetryByDefault(t *testing.T) {
	var attempts atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			http.Error(w, "temporary", http.StatusServiceUnavailable)
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogOff,
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	ctx := dropbox.NewContext(config)
	_, _, err := ctx.Execute(dropbox.Request{
		Host:      "api",
		Namespace: "users",
		Route:     "get_current_account",
		Auth:      "user",
		Style:     "rpc",
	}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts.Load() != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts.Load())
	}
}

func TestExecuteUsesConfigRetryPolicy(t *testing.T) {
	var attempts atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if attempts.Add(1) == 1 {
				http.Error(w, "temporary", http.StatusServiceUnavailable)
				return
			}
			_, _ = w.Write([]byte(`{}`))
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogOff,
		RetryPolicy: retryTestPolicyPtr(1),
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	ctx := dropbox.NewContext(config)
	resp, respBody, err := ctx.Execute(dropbox.Request{
		Host:      "api",
		Namespace: "users",
		Route:     "get_current_account",
		Auth:      "user",
		Style:     "rpc",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if attempts.Load() != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts.Load())
	}
	if string(resp) != "{}" {
		t.Errorf("unexpected response: %s", string(resp))
	}
	if respBody != nil {
		t.Errorf("unexpected response body: %v", respBody)
	}
}

func TestExecuteContextDoesNotRetryWithoutPolicy(t *testing.T) {
	var attempts atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			http.Error(w, "temporary", http.StatusServiceUnavailable)
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogOff,
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	callCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	ctx := dropbox.NewContext(config)
	_, _, err := ctx.ExecuteContext(callCtx, dropbox.Request{
		Host:      "api",
		Namespace: "users",
		Route:     "get_current_account",
		Auth:      "user",
		Style:     "rpc",
	}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts.Load() != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts.Load())
	}
}

func TestExecuteContextUsesConfigRetryPolicy(t *testing.T) {
	var attempts atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if attempts.Add(1) == 1 {
				http.Error(w, "temporary", http.StatusServiceUnavailable)
				return
			}
			_, _ = w.Write([]byte(`{}`))
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogOff,
		RetryPolicy: retryTestPolicyPtr(1),
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	callCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	ctx := dropbox.NewContext(config)
	_, _, err := ctx.ExecuteContext(callCtx, dropbox.Request{
		Host:      "api",
		Namespace: "users",
		Route:     "get_current_account",
		Auth:      "user",
		Style:     "rpc",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if attempts.Load() != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts.Load())
	}
}

func TestExecuteDoesNotRetryUnlistedStatus(t *testing.T) {
	var attempts atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			http.Error(w, "temporary", http.StatusInternalServerError)
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogOff,
		RetryPolicy: retryTestPolicyPtr(1),
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	ctx := dropbox.NewContext(config)
	_, _, err := ctx.Execute(dropbox.Request{
		Host:      "api",
		Namespace: "users",
		Route:     "get_current_account",
		Auth:      "user",
		Style:     "rpc",
	}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts.Load() != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts.Load())
	}
}

func TestExecuteDoesNotRetryNetworkError(t *testing.T) {
	var attempts atomic.Int32
	networkErr := errors.New("network unavailable")
	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			attempts.Add(1)
			return nil, networkErr
		}),
	}

	config := dropbox.Config{Client: client, LogLevel: dropbox.LogOff,
		RetryPolicy: retryTestPolicyPtr(1),
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL("http://example.com", namespace, route)
		}}
	ctx := dropbox.NewContext(config)
	_, _, err := ctx.Execute(dropbox.Request{
		Host:      "api",
		Namespace: "users",
		Route:     "get_current_account",
		Auth:      "user",
		Style:     "rpc",
	}, nil)
	if !errors.Is(err, networkErr) {
		t.Fatalf("expected network error, got %v", err)
	}
	if attempts.Load() != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts.Load())
	}
}

func TestNewContextUsesTokenSourceBeforeToken(t *testing.T) {
	var authHeader string
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader = r.Header.Get("Authorization")
			_, _ = w.Write([]byte(`{}`))
		}))
	defer ts.Close()

	config := dropbox.Config{
		Token:       "static-token",
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "source-token"}),
		LogLevel:    dropbox.LogOff,
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		},
	}
	executeTestRequest(t, config)
	if authHeader != "Bearer source-token" {
		t.Fatalf("Authorization = %q, want %q", authHeader, "Bearer source-token")
	}
}

func TestNewContextUsesTokenWhenTokenSourceUnset(t *testing.T) {
	var authHeader string
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader = r.Header.Get("Authorization")
			_, _ = w.Write([]byte(`{}`))
		}))
	defer ts.Close()

	config := dropbox.Config{
		Token:    "static-token",
		LogLevel: dropbox.LogOff,
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		},
	}
	executeTestRequest(t, config)
	if authHeader != "Bearer static-token" {
		t.Fatalf("Authorization = %q, want %q", authHeader, "Bearer static-token")
	}
}

func TestNewContextClientTakesPrecedenceOverTokenSource(t *testing.T) {
	var authHeader string
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader = r.Header.Get("Authorization")
			_, _ = w.Write([]byte(`{}`))
		}))
	defer ts.Close()

	config := dropbox.Config{
		Client:      ts.Client(),
		Token:       "static-token",
		TokenSource: failingTokenSource{t: t},
		LogLevel:    dropbox.LogOff,
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		},
	}
	executeTestRequest(t, config)
	if authHeader != "" {
		t.Fatalf("Authorization = %q, want empty", authHeader)
	}
}

func TestNewContextNoAuthDoesNotUseTokenSource(t *testing.T) {
	var authHeader string
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader = r.Header.Get("Authorization")
			_, _ = w.Write([]byte(`{}`))
		}))
	defer ts.Close()

	config := dropbox.Config{
		TokenSource: failingTokenSource{t: t},
		LogLevel:    dropbox.LogOff,
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		},
	}
	executeTestRequestWithAuth(t, config, "noauth")
	if authHeader != "" {
		t.Fatalf("Authorization = %q, want empty", authHeader)
	}
}

func TestContextCancellation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogOff,
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	client := users.NewContext(config)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetCurrentAccountContext(ctx)
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
	if !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("expected context canceled error, got: %v", err)
	}
}

func TestRetryRPCOn503ThenSuccess(t *testing.T) {
	var attempts atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if attempts.Add(1) == 1 {
				http.Error(w, "temporary", http.StatusServiceUnavailable)
				return
			}
			_, _ = w.Write([]byte(`{}`))
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogOff,
		RetryPolicy: retryTestPolicyPtr(1),
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	ctx := dropbox.NewContext(config)
	resp, respBody, err := ctx.ExecuteContext(context.Background(), dropbox.Request{
		Host:      "api",
		Namespace: "users",
		Route:     "get_current_account",
		Auth:      "user",
		Style:     "rpc",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if attempts.Load() != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts.Load())
	}
	if string(resp) != "{}" {
		t.Errorf("unexpected response: %s", string(resp))
	}
	if respBody != nil {
		t.Errorf("unexpected response body: %v", respBody)
	}
}

func TestRetryStatusWhenErrorBodyReadFails(t *testing.T) {
	var attempts atomic.Int32
	readErr := errors.New("read response body")
	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if attempts.Add(1) == 1 {
				return &http.Response{
					StatusCode: http.StatusServiceUnavailable,
					Header:     http.Header{},
					Body:       failingReadCloser{err: readErr},
					Request:    r,
				}, nil
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{},
				Body:       io.NopCloser(strings.NewReader(`{}`)),
				Request:    r,
			}, nil
		}),
	}

	config := dropbox.Config{Client: client, LogLevel: dropbox.LogOff,
		RetryPolicy: retryTestPolicyPtr(1),
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL("http://example.com", namespace, route)
		}}
	ctx := dropbox.NewContext(config)
	resp, respBody, err := ctx.ExecuteContext(context.Background(), dropbox.Request{
		Host:      "api",
		Namespace: "users",
		Route:     "get_current_account",
		Auth:      "user",
		Style:     "rpc",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if attempts.Load() != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts.Load())
	}
	if string(resp) != "{}" {
		t.Errorf("unexpected response: %s", string(resp))
	}
	if respBody != nil {
		t.Errorf("unexpected response body: %v", respBody)
	}
}

func TestRetryAfterHeader(t *testing.T) {
	var attempts atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if attempts.Add(1) == 1 {
				w.Header().Set("Retry-After", "0")
				http.Error(w, "rate limited", http.StatusTooManyRequests)
				return
			}
			_, _ = w.Write([]byte(`{}`))
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogOff,
		RetryPolicy: retryTestPolicyPtr(1),
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	ctx := dropbox.NewContext(config)
	_, _, err := ctx.ExecuteContext(context.Background(), dropbox.Request{
		Host:      "api",
		Namespace: "users",
		Route:     "get_current_account",
		Auth:      "user",
		Style:     "rpc",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if attempts.Load() != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts.Load())
	}
}

func TestRetry409ConfiguredStoneTag(t *testing.T) {
	var attempts atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if attempts.Add(1) == 1 {
				w.WriteHeader(http.StatusConflict)
				_, _ = w.Write([]byte(`{"error_summary":"too_many_write_operations/..","error":{".tag":"too_many_write_operations"}}`))
				return
			}
			_, _ = w.Write([]byte(`{}`))
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogOff,
		RetryPolicy: retryTestPolicyPtrWith409Tags(1, "too_many_write_operations"),
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	ctx := dropbox.NewContext(config)
	_, _, err := ctx.ExecuteContext(context.Background(), dropbox.Request{
		Host:      "api",
		Namespace: "files",
		Route:     "upload",
		Auth:      "user",
		Style:     "rpc",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if attempts.Load() != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts.Load())
	}
}

func TestRetry409UnconfiguredStoneTagNotRetried(t *testing.T) {
	var attempts atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"error_summary":"too_many_write_operations/..","error":{".tag":"too_many_write_operations"}}`))
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogOff,
		RetryPolicy: retryTestPolicyPtr(1),
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	ctx := dropbox.NewContext(config)
	_, _, err := ctx.ExecuteContext(context.Background(), dropbox.Request{
		Host:      "api",
		Namespace: "files",
		Route:     "upload",
		Auth:      "user",
		Style:     "rpc",
	}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts.Load() != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts.Load())
	}
}

func TestContextCancellationDuringRetryBackoff(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "temporary", http.StatusServiceUnavailable)
		}))
	defer ts.Close()

	policy := retry.Policy{
		MaxRetries:     1,
		InitialBackoff: time.Hour,
		MaxBackoff:     time.Hour,
		MaxRetryAfter:  time.Hour,
	}
	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogOff,
		RetryPolicy: &policy,
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	base, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	ctx := dropbox.NewContext(config)
	_, _, err := ctx.ExecuteContext(base, dropbox.Request{
		Host:      "api",
		Namespace: "users",
		Route:     "get_current_account",
		Auth:      "user",
		Style:     "rpc",
	}, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
}

func TestUploadReadSeekCloserRetriedAndRewound(t *testing.T) {
	var attempts atomic.Int32
	var mu sync.Mutex
	var bodies []string
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("read request body: %v", err)
			}
			mu.Lock()
			bodies = append(bodies, string(body))
			mu.Unlock()

			if attempts.Add(1) == 1 {
				http.Error(w, "temporary", http.StatusServiceUnavailable)
				return
			}
			_, _ = w.Write([]byte(`{}`))
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogOff,
		RetryPolicy: retryTestPolicyPtr(1),
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	body, err := os.CreateTemp(t.TempDir(), "upload")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := body.Close(); err != nil {
			t.Error(err)
		}
	}()
	if _, err := body.WriteString("payload"); err != nil {
		t.Fatal(err)
	}
	if _, err := body.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	ctx := dropbox.NewContext(config)
	_, _, err = ctx.ExecuteContext(context.Background(), dropbox.Request{
		Host:      "content",
		Namespace: "files",
		Route:     "upload",
		Auth:      "user",
		Style:     "upload",
		Arg:       map[string]string{"path": "/x"},
	}, body)
	if err != nil {
		t.Fatal(err)
	}
	if attempts.Load() != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts.Load())
	}
	mu.Lock()
	defer mu.Unlock()
	if len(bodies) != 2 || bodies[0] != "payload" || bodies[1] != "payload" {
		t.Fatalf("unexpected request bodies: %#v", bodies)
	}
}

func TestUploadNonSeekableNotRetried(t *testing.T) {
	var attempts atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			_, _ = io.Copy(io.Discard, r.Body)
			http.Error(w, "temporary", http.StatusServiceUnavailable)
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogOff,
		RetryPolicy: retryTestPolicyPtr(1),
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	body := struct{ io.Reader }{Reader: strings.NewReader("payload")}
	ctx := dropbox.NewContext(config)
	_, _, err := ctx.ExecuteContext(context.Background(), dropbox.Request{
		Host:      "content",
		Namespace: "files",
		Route:     "upload",
		Auth:      "user",
		Style:     "upload",
		Arg:       map[string]string{"path": "/x"},
	}, body)
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts.Load() != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts.Load())
	}
}

func TestDownloadRetryBeforeReturningBody(t *testing.T) {
	var attempts atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if attempts.Add(1) == 1 {
				http.Error(w, "temporary", http.StatusServiceUnavailable)
				return
			}
			w.Header().Set("Dropbox-API-Result", `{"name":"x"}`)
			_, _ = w.Write([]byte("content"))
		}))
	defer ts.Close()

	config := dropbox.Config{Client: ts.Client(), LogLevel: dropbox.LogOff,
		RetryPolicy: retryTestPolicyPtr(1),
		URLGenerator: func(hostType string, namespace string, route string) string {
			return generateURL(ts.URL, namespace, route)
		}}
	ctx := dropbox.NewContext(config)
	resp, respBody, err := ctx.ExecuteContext(context.Background(), dropbox.Request{
		Host:      "content",
		Namespace: "files",
		Route:     "download",
		Auth:      "user",
		Style:     "download",
		Arg:       map[string]string{"path": "/x"},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := respBody.Close(); err != nil {
			t.Error(err)
		}
	}()
	content, err := io.ReadAll(respBody)
	if err != nil {
		t.Fatal(err)
	}
	if attempts.Load() != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts.Load())
	}
	if string(resp) != `{"name":"x"}` {
		t.Fatalf("unexpected response metadata: %s", string(resp))
	}
	if string(content) != "content" {
		t.Fatalf("unexpected content: %s", string(content))
	}
}

type failingTokenSource struct {
	t *testing.T
}

func (s failingTokenSource) Token() (*oauth2.Token, error) {
	s.t.Helper()
	s.t.Fatal("TokenSource should not be used")
	return nil, nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type failingReadCloser struct {
	err error
}

func (r failingReadCloser) Read([]byte) (int, error) {
	return 0, r.err
}

func (r failingReadCloser) Close() error {
	return nil
}

func executeTestRequest(t *testing.T, config dropbox.Config) {
	t.Helper()
	executeTestRequestWithAuth(t, config, "user")
}

func executeTestRequestWithAuth(t *testing.T, config dropbox.Config, auth string) {
	t.Helper()

	ctx := dropbox.NewContext(config)
	resp, respBody, err := ctx.Execute(dropbox.Request{
		Host:      "api",
		Namespace: "users",
		Route:     "get_current_account",
		Auth:      auth,
		Style:     "rpc",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(resp) != "{}" {
		t.Errorf("unexpected response: %s", string(resp))
	}
	if respBody != nil {
		t.Errorf("unexpected response body: %v", respBody)
	}
}
