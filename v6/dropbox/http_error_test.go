package dropbox_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

func TestHTTPError(t *testing.T) {
	for _, tc := range []struct {
		status        int
		body, message string
		target        any
	}{
		{400, "bad input", "bad input", &auth.BadRequest{}},
		{401, `{"error_summary":"unauthorized","error":{".tag":"user_suspended"}}`, "unauthorized", &auth.AuthAPIError{}},
		{403, `{"error_summary":"forbidden","error":{".tag":"invalid_account_type"}}`, "forbidden", &auth.AccessAPIError{}},
		{409, `{"error_summary":"missing","error":{".tag":"path","path":{".tag":"not_found"}}}`, "missing", &files.GetMetadataAPIError{}},
		{429, `{"error_summary":"limited","error":{"reason":{".tag":"too_many_requests"},"retry_after":300}}`, "limited", &auth.RateLimitAPIError{}},
		{500, "server error", "server error", &auth.ServerError{}},
		{503, "unavailable", "unavailable", &auth.ServerError{}},
		{418, "unexpected", "Unexpected error: unexpected (code: 418)", &dropbox.SDKInternalError{}},
	} {
		t.Run(fmt.Sprint(tc.status), func(t *testing.T) {
			cfg := httpErrorConfig(func(r *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: tc.status, Header: http.Header{"X-Dropbox-Request-Id": {"request-id"}, "Retry-After": {"10"}}, Body: io.NopCloser(strings.NewReader(tc.body))}, nil
			})
			_, err := files.New(cfg).GetMetadata(files.NewGetMetadataArg("/missing"))
			checkHTTPError(t, err, tc.status, "request-id")
			if !errors.As(err, tc.target) || err.Error() != tc.message {
				t.Fatalf("got %T: %v", err, err)
			}
			var server auth.ServerError
			if errors.As(err, &server) && server.StatusCode != tc.status {
				t.Fatalf("server status: %d", server.StatusCode)
			}
			wrapped, _ := json.Marshal(err)
			underlying, _ := json.Marshal(errors.Unwrap(err))
			if string(wrapped) != string(underlying) || strings.Contains(string(wrapped), "request-id") {
				t.Fatalf("JSON changed: %s", wrapped)
			}
		})
	}
}

func TestHTTPErrorMalformedJSON(t *testing.T) {
	for _, status := range []int{401, 403, 409, 429} {
		for _, body := range []string{"", "{"} {
			t.Run(fmt.Sprintf("%d/%q", status, body), func(t *testing.T) {
				cfg := httpErrorConfig(func(*http.Request) (*http.Response, error) {
					return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body))}, nil
				})
				_, err := files.New(cfg).GetMetadata(files.NewGetMetadataArg("/missing"))
				checkHTTPError(t, err, status, "")
				var syntax *json.SyntaxError
				if !errors.As(err, &syntax) || err.Error() != syntax.Error() {
					t.Fatalf("lost decoding error: %v", err)
				}
			})
		}
	}
}

type errorBody struct {
	io.Reader
	readErr, closeErr error
	closed            bool
}

func (b *errorBody) Read(p []byte) (int, error) {
	if b.readErr != nil {
		return 0, b.readErr
	}
	return b.Reader.Read(p)
}
func (b *errorBody) Close() error { b.closed = true; return b.closeErr }

func TestHTTPErrorBodyFailures(t *testing.T) {
	cause := errors.New("body failure")
	for _, status := range []int{200, 503} {
		for _, closeFailure := range []bool{false, true} {
			t.Run(fmt.Sprintf("%d/close=%v", status, closeFailure), func(t *testing.T) {
				body := &errorBody{Reader: strings.NewReader("response")}
				if closeFailure {
					body.closeErr = cause
				} else {
					body.readErr = cause
				}
				cfg := httpErrorConfig(func(*http.Request) (*http.Response, error) {
					return &http.Response{StatusCode: status, Header: http.Header{"X-Dropbox-Request-Id": {"body-id"}}, Body: body}, nil
				})
				_, err := files.New(cfg).GetMetadata(files.NewGetMetadataArg("/missing"))
				checkHTTPError(t, err, status, "body-id")
				if !errors.Is(err, cause) || err.Error() != cause.Error() || !body.closed {
					t.Fatalf("lost cause or close: %v", err)
				}
			})
		}
	}
}

func TestHTTPErrorTransportFailures(t *testing.T) {
	cause := errors.New("transport failure")
	for _, withResponse := range []bool{false, true} {
		t.Run(fmt.Sprint(withResponse), func(t *testing.T) {
			cfg := httpErrorConfig(func(*http.Request) (*http.Response, error) {
				if !withResponse {
					return nil, cause
				}
				return &http.Response{StatusCode: 302, Header: http.Header{"Location": {"/redirect"}, "X-Dropbox-Request-Id": {"redirect-id"}}, Body: http.NoBody}, nil
			})
			cfg.Client.CheckRedirect = func(*http.Request, []*http.Request) error { return cause }
			c := dropbox.NewContext(cfg)
			_, _, err := c.Execute(dropbox.Request{Style: "rpc"}, nil)
			var httpErr *dropbox.HTTPError
			var endpoint files.GetMetadataAPIError
			if !errors.Is(err, cause) || auth.ParseError(err, &endpoint) != err {
				t.Fatalf("lost transport error: %v", err)
			}
			if withResponse {
				checkHTTPError(t, err, 302, "redirect-id")
			} else if errors.As(err, &httpErr) {
				t.Fatal("wrapped network error without a response")
			}
		})
	}
}

func TestHTTPErrorRetryExhaustion(t *testing.T) {
	attempts := 0
	cfg := httpErrorConfig(func(*http.Request) (*http.Response, error) {
		attempts++
		status := 429
		if attempts == 3 {
			status = 503
		}
		return &http.Response{StatusCode: status, Header: http.Header{"X-Dropbox-Request-Id": {fmt.Sprint(attempts)}}, Body: io.NopCloser(strings.NewReader("unavailable"))}, nil
	})
	cfg.RetryPolicy = retryTestPolicyPtr(2)
	_, err := files.New(cfg).GetMetadata(files.NewGetMetadataArg("/missing"))
	checkHTTPError(t, err, 503, "3")
	if attempts != 3 {
		t.Fatalf("attempts: %d", attempts)
	}
}

func httpErrorConfig(transport roundTripFunc) dropbox.Config {
	return dropbox.Config{Client: &http.Client{Transport: transport}, URLGenerator: func(_, _, _ string) string { return "http://example.com" }}
}
func checkHTTPError(t *testing.T, err error, status int, id string) {
	t.Helper()
	var httpErr *dropbox.HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != status || httpErr.RequestID != id {
		t.Fatalf("HTTPError = %+v, want %d/%q", httpErr, status, id)
	}
}
