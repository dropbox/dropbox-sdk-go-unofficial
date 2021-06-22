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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/users"
)

func generateURL(base string, namespace string, route string) string {
	return fmt.Sprintf("%s/%s/%s", base, namespace, route)
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
	client := users.New(config)
	v, e := client.GetCurrentAccount()
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
	client := users.New(config)
	_, e := client.GetCurrentAccount()
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
