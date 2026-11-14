package dropbox_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type metadataCarrier interface {
	ResponseMetadata() *dropbox.ResponseMetadata
}

type metadataBody struct {
	io.Reader
	closed            bool
	readErr, closeErr error
	onClose           func()
}

func (b *metadataBody) Read(p []byte) (int, error) {
	if b.readErr != nil {
		return 0, b.readErr
	}
	return b.Reader.Read(p)
}
func (b *metadataBody) Close() error {
	b.closed = true
	if b.onClose != nil {
		b.onClose()
	}
	return b.closeErr
}

func metadataConfig(status int, headers http.Header, body io.ReadCloser) dropbox.Config {
	return dropbox.Config{Client: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: status, Header: headers, Body: body, Request: r}, nil
	})}, URLGenerator: func(_, _, _ string) string { return "http://example.com" }}
}
func requireMetadata(t *testing.T, err error, status int, id string) *dropbox.ResponseMetadata {
	t.Helper()
	var carrier metadataCarrier
	if !errors.As(err, &carrier) {
		t.Fatalf("%T lacks metadata: %v", err, err)
	}
	m := carrier.ResponseMetadata()
	if m == nil || m.StatusCode != status || m.RequestID != id {
		t.Fatalf("metadata = %+v, want %d / %q", m, status, id)
	}
	return m
}
func TestResponseMetadataPublicErrors(t *testing.T) {
	for _, tc := range []struct {
		status int
		body   string
		target error
	}{
		{400, "bad input", auth.BadRequest{}},
		{401, `{"error_summary":"summary","error":{".tag":"user_suspended"}}`, auth.AuthAPIError{}},
		{403, `{"error_summary":"summary","error":{".tag":"invalid_account_type"}}`, auth.AccessAPIError{}},
		{409, `{"error_summary":"summary","error":{".tag":"path","path":{".tag":"not_found"}}}`, files.GetMetadataAPIError{}},
		{429, `{"error_summary":"summary","error":{"reason":{".tag":"too_many_requests"},"retry_after":300}}`, auth.RateLimitAPIError{}},
		{500, "server error", auth.ServerError{}}, {503, "server error", auth.ServerError{}},
		{418, "unexpected", dropbox.SDKInternalError{}},
	} {
		t.Run(fmt.Sprint(tc.status), func(t *testing.T) {
			headers := http.Header{"x-dRoPbOx-ReQuEsT-iD": {"request-secret"}, "x-MuLtI": {"one", "two"}, "X-Multi": {"three"}, "Retry-After": {"10"}, "X-Server-Response-Time": {"42"}}
			body := &metadataBody{Reader: strings.NewReader(tc.body)}
			// Mutate the original map and slices on Close, before the error is parsed.
			body.onClose = func() { headers["x-MuLtI"][0] = "changed"; headers["x-dRoPbOx-ReQuEsT-iD"][0] = "changed" }
			_, err := files.New(metadataConfig(tc.status, headers, body)).GetMetadata(files.NewGetMetadataArg("/missing"))
			if reflect.TypeOf(err) != reflect.TypeOf(tc.target) {
				t.Fatalf("got %T, want %T", err, tc.target)
			}
			target := reflect.New(reflect.TypeOf(tc.target))
			if !errors.As(err, target.Interface()) {
				t.Fatal("concrete errors.As failed")
			}
			m := requireMetadata(t, fmt.Errorf("wrapped: %w", err), tc.status, "request-secret")
			if !body.closed {
				t.Fatal("body not closed")
			}
			vals := m.Headers().Values("x-multi")
			if len(vals) != 3 || !containsString(vals, "one") || !containsString(vals, "two") || !containsString(vals, "three") {
				t.Fatalf("headers: %v", vals)
			}
			h := m.Headers()
			h.Values("x-multi")[0] = "modified"
			h.Set("X-Dropbox-Request-Id", "modified")
			if m.Headers().Get("X-Dropbox-Request-Id") != "request-secret" || containsString(m.Headers().Values("X-Multi"), "modified") {
				t.Fatal("Headers aliases snapshot")
			}
			m.StatusCode = 0
			requireMetadata(t, err, tc.status, "request-secret")
			encoded, e := json.Marshal(err)
			if e != nil {
				t.Fatal(e)
			}
			if strings.Contains(string(encoded), "request-secret") || strings.Contains(string(encoded), "X-Multi") || strings.Contains(err.Error(), "request-secret") {
				t.Fatalf("metadata leaked: %s / %s", encoded, err)
			}
			switch e := err.(type) {
			case auth.ServerError:
				if e.StatusCode != tc.status || e.Error() != tc.body {
					t.Fatalf("server error: %+v", e)
				}
			case auth.BadRequest:
				if e.Error() != tc.body {
					t.Fatal(e)
				}
			case auth.AuthAPIError:
				if e.AuthError.Tag != "user_suspended" {
					t.Fatal(e)
				}
			case auth.AccessAPIError:
				if e.AccessError.Tag != "invalid_account_type" {
					t.Fatal(e)
				}
			case files.GetMetadataAPIError:
				if e.EndpointError.Path.Tag != "not_found" || e.Error() != "summary" {
					t.Fatal(e)
				}
			case auth.RateLimitAPIError:
				if e.RateLimitError.RetryAfter != 300 {
					t.Fatal(e)
				}
			}
			// The same promoted getter also works on pointers to the existing types.
			target.Elem().Set(reflect.ValueOf(err))
			requireMetadata(t, target.Interface().(error), tc.status, "request-secret")
		})
	}
}
func containsString(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func TestResponseMetadataDecodeFailures(t *testing.T) {
	for _, status := range []int{401, 403, 409, 429} {
		for _, raw := range []string{"", "{", `{"error_summary":42}`} {
			t.Run(fmt.Sprintf("%d/%s", status, raw), func(t *testing.T) {
				body := &metadataBody{Reader: strings.NewReader(raw)}
				_, err := files.New(metadataConfig(status, nil, body)).GetMetadata(files.NewGetMetadataArg("/missing"))
				m := requireMetadata(t, err, status, "")
				if len(m.Headers()) != 0 || !body.closed {
					t.Fatal("missing headers or close")
				}
				var syntax *json.SyntaxError
				var mismatch *json.UnmarshalTypeError
				if !errors.As(err, &syntax) && !errors.As(err, &mismatch) {
					t.Fatalf("lost JSON cause: %T %v", err, err)
				}
				if errors.Unwrap(err) == nil {
					t.Fatal("missing Unwrap")
				}
				var structured files.GetMetadataAPIError
				if errors.As(err, &structured) {
					t.Fatal("decoding failure became structured error")
				}
			})
		}
	}
}
func TestResponseMetadataBodyFailures(t *testing.T) {
	cause := errors.New("body failure")
	for _, status := range []int{200, 503} {
		for _, closeFailure := range []bool{false, true} {
			t.Run(fmt.Sprintf("%d/close=%v", status, closeFailure), func(t *testing.T) {
				b := &metadataBody{Reader: strings.NewReader("response")}
				if closeFailure {
					b.closeErr = cause
				} else {
					b.readErr = cause
				}
				c := dropbox.NewContext(metadataConfig(status, http.Header{"X-Dropbox-Request-Id": {"body-id"}}, b))
				_, _, err := c.Execute(dropbox.Request{Style: "rpc"}, nil)
				requireMetadata(t, err, status, "body-id")
				if !errors.Is(err, cause) || err.Error() != cause.Error() || !b.closed {
					t.Fatalf("lost cause or close: %v", err)
				}
			})
		}
	}
}
func TestResponseMetadataAbsentWithoutResponse(t *testing.T) {
	for _, err := range []error{dropbox.APIError{}, &dropbox.APIError{}, dropbox.SDKInternalError{StatusCode: 503}, auth.ServerError{StatusCode: 503}, auth.BadRequest{}, files.GetMetadataAPIError{}} {
		var carrier metadataCarrier
		if !errors.As(err, &carrier) || carrier.ResponseMetadata() != nil {
			t.Fatalf("invented metadata: %T", err)
		}
	}
	cause := errors.New("network failure")
	cfg := metadataConfig(0, nil, nil)
	cfg.Client.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, cause })
	c := dropbox.NewContext(cfg)
	_, _, err := c.ExecuteContext(context.Background(), dropbox.Request{}, nil)
	var carrier metadataCarrier
	if !errors.Is(err, cause) || errors.As(err, &carrier) {
		t.Fatalf("network error: %v", err)
	}
}
func TestResponseMetadataRetryExhaustion(t *testing.T) {
	var attempts int
	cfg := metadataConfig(503, nil, nil)
	cfg.RetryPolicy = retryTestPolicyPtr(2)
	cfg.Client.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		attempts++
		return &http.Response{StatusCode: 503, Header: http.Header{"X-Dropbox-Request-Id": {fmt.Sprint(attempts)}}, Body: io.NopCloser(strings.NewReader("server error")), Request: r}, nil
	})
	_, err := files.New(cfg).GetMetadata(files.NewGetMetadataArg("/missing"))
	if attempts != 3 {
		t.Fatalf("attempts: %d", attempts)
	}
	requireMetadata(t, err, 503, "3")
}
