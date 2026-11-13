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

package oauth_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	dropboxoauth "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/oauth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/users"
	"golang.org/x/oauth2"
)

var _ oauth2.TokenSource = (*dropboxoauth.CachedTokenSource)(nil)

func TestCachedTokenSourceInitialSnapshot(t *testing.T) {
	expiry := time.Date(2030, 1, 2, 3, 4, 5, 123, time.FixedZone("offset", 3600))
	for _, test := range []struct {
		name  string
		token *oauth2.Token
		ok    bool
	}{
		{"known expiry", &oauth2.Token{AccessToken: "access", Expiry: expiry}, true},
		{"unknown expiry", &oauth2.Token{AccessToken: "access"}, true},
		{"expired", &oauth2.Token{AccessToken: "access", Expiry: time.Unix(1, 0)}, true},
		{"nil", nil, false},
		{"empty", &oauth2.Token{}, false},
		{"refresh only", &oauth2.Token{RefreshToken: "secret", Expiry: expiry}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			client := httpClient(func(*http.Request) (*http.Response, error) {
				t.Error("Snapshot made an HTTP request")
				return nil, fmt.Errorf("unexpected HTTP request")
			})
			source := dropboxoauth.NewCachedTokenSource(context.Background(), "app-key", test.token, dropboxoauth.WithHTTPClient(client))
			for i := 0; i < 3; i++ {
				snapshot, ok := source.Snapshot()
				if ok != test.ok {
					t.Fatalf("Snapshot ok = %v, want %v", ok, test.ok)
				}
				want := dropboxoauth.AccessTokenSnapshot{}
				if test.ok {
					want = dropboxoauth.AccessTokenSnapshot{AccessToken: test.token.AccessToken, Expiry: test.token.Expiry}
				}
				if snapshot != want {
					t.Fatalf("Snapshot = %#v, want %#v", snapshot, want)
				}
			}
			if test.token.Valid() {
				token, err := source.Token()
				if err != nil {
					t.Fatal(err)
				}
				if token.AccessToken != test.token.AccessToken || token.Expiry != test.token.Expiry {
					t.Fatal("Token changed valid initial credentials")
				}
			}
		})
	}
}

func TestCachedTokenSourceCopiesInitialToken(t *testing.T) {
	initial := &oauth2.Token{AccessToken: "access", RefreshToken: "secret", Expiry: time.Now().Add(time.Hour)}
	source := dropboxoauth.NewCachedTokenSource(context.Background(), "app-key", initial)
	before, _ := source.Snapshot()
	encoded, err := json.Marshal(before)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "secret") || strings.Contains(string(encoded), "RefreshToken") {
		t.Fatal("snapshot exposes refresh credentials")
	}
	initial.AccessToken = "changed"
	initial.Expiry = time.Unix(1, 0)
	token, err := source.Token()
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken != before.AccessToken || token.Expiry != before.Expiry {
		t.Fatal("constructor retained caller-owned token")
	}
	changed := before
	changed.AccessToken = "changed snapshot"
	changed.Expiry = time.Time{}
	if after, _ := source.Snapshot(); after != before || after.AccessToken == changed.AccessToken || after.Expiry.Equal(changed.Expiry) {
		t.Fatal("modifying snapshot changed cache")
	}
}

func TestCachedTokenSourceRefreshAndImmutability(t *testing.T) {
	calls := 0
	client := httpClient(func(req *http.Request) (*http.Response, error) {
		calls++
		form, err := readForm(req)
		if err != nil {
			return nil, err
		}
		assertQueryValue(t, form, "grant_type", "refresh_token")
		assertQueryValue(t, form, "refresh_token", "secret")
		assertQueryValue(t, form, "client_id", "app-key")
		assertQueryValue(t, form, "client_secret", "app-secret")
		if req.URL.Host != "api.example.com" {
			t.Errorf("refresh host = %q", req.URL.Host)
		}
		// OAuth2's early-expiry window makes this token immediately stale,
		// allowing a second refresh without sleeping or mutating cached tokens.
		if calls == 1 {
			return jsonResponse(http.StatusOK, `{"access_token":"first","expires_in":1}`), nil
		}
		return jsonResponse(http.StatusOK, `{"access_token":"second","expires_in":3600}`), nil
	})
	source := dropboxoauth.NewCachedTokenSource(context.Background(), "app-key", &oauth2.Token{RefreshToken: "secret"},
		dropboxoauth.WithHTTPClient(client), dropboxoauth.WithAppSecret("app-secret"), dropboxoauth.WithDomain(".example.com"))
	firstToken, err := source.Token()
	if err != nil {
		t.Fatal(err)
	}
	first, ok := source.Snapshot()
	if !ok || first.AccessToken != "first" || first.Expiry != firstToken.Expiry || first.Expiry.IsZero() {
		t.Fatalf("unexpected first snapshot: %#v, ok %v", first, ok)
	}
	secondToken, err := source.Token()
	if err != nil {
		t.Fatal(err)
	}
	second, ok := source.Snapshot()
	if !ok || second.AccessToken != "second" || second.Expiry != secondToken.Expiry || !second.Expiry.After(first.Expiry) {
		t.Fatalf("unexpected second snapshot: %#v, ok %v", second, ok)
	}
	if first.AccessToken != "first" || first.Expiry != firstToken.Expiry {
		t.Fatal("subsequent refresh changed previous snapshot")
	}
	if _, err := source.Token(); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("refresh calls = %d, want 2", calls)
	}
}

func TestCachedTokenSourceFailedRefresh(t *testing.T) {
	for _, response := range []struct {
		name   string
		status int
		body   string
	}{
		{"HTTP error", http.StatusBadRequest, `{"error":"invalid_grant"}`},
		{"invalid expiry", http.StatusOK, `{"access_token":"partial","expires_in":"invalid"}`},
		{"missing access token", http.StatusOK, `{"expires_in":3600}`},
	} {
		for _, hasAccess := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/cached=%v", response.name, hasAccess), func(t *testing.T) {
				initial := &oauth2.Token{RefreshToken: "secret", Expiry: time.Unix(1, 0)}
				if hasAccess {
					initial.AccessToken = "old"
				}
				client := httpClient(func(*http.Request) (*http.Response, error) {
					return jsonResponse(response.status, response.body), nil
				})
				source := dropboxoauth.NewCachedTokenSource(context.Background(), "app-key", initial, dropboxoauth.WithHTTPClient(client))
				before, beforeOK := source.Snapshot()
				if _, err := source.Token(); err == nil {
					t.Fatal("Token succeeded, want refresh error")
				}
				if after, ok := source.Snapshot(); after != before || ok != beforeOK {
					t.Fatal("failed refresh changed snapshot")
				}
			})
		}
	}
}

func TestCachedTokenSourceConcurrentRefreshAndSnapshots(t *testing.T) {
	entered := make(chan struct{})
	release := make(chan struct{})
	var releaseOnce sync.Once
	unblock := func() { releaseOnce.Do(func() { close(release) }) }
	defer unblock()
	var calls atomic.Int32
	client := httpClient(func(*http.Request) (*http.Response, error) {
		if calls.Add(1) == 1 {
			close(entered)
		}
		<-release
		return jsonResponse(http.StatusOK, `{"access_token":"new","expires_in":3600}`), nil
	})
	source := dropboxoauth.NewCachedTokenSource(context.Background(), "app-key", &oauth2.Token{
		AccessToken: "old", RefreshToken: "secret", Expiry: time.Unix(1, 0),
	}, dropboxoauth.WithHTTPClient(client))
	before, _ := source.Snapshot()
	var tokens sync.WaitGroup
	for i := 0; i < 16; i++ {
		tokens.Go(func() {
			token, err := source.Token()
			if err != nil {
				t.Error(err)
			} else if token.AccessToken != "new" {
				t.Errorf("token = %q, want new", token.AccessToken)
			}
		})
	}
	<-entered
	// All these reads must finish while the refresh remains blocked on release.
	var readers sync.WaitGroup
	for i := 0; i < 16; i++ {
		readers.Go(func() {
			for j := 0; j < 100; j++ {
				if snapshot, ok := source.Snapshot(); !ok || snapshot != before {
					t.Error("in-flight refresh changed cached pair")
				}
			}
		})
	}
	readers.Wait()
	// Keep readers active across publication and accept only complete pairs.
	stop := make(chan struct{})
	for i := 0; i < 16; i++ {
		readers.Go(func() {
			for {
				select {
				case <-stop:
					return
				default:
				}
				snapshot, ok := source.Snapshot()
				if !ok || (snapshot != before && (snapshot.AccessToken != "new" || snapshot.Expiry.IsZero() || snapshot.Expiry.Equal(before.Expiry))) {
					t.Error("inconsistent token/expiry pair")
					return
				}
			}
		})
	}
	unblock()
	tokens.Wait()
	close(stop)
	readers.Wait()
	if calls.Load() != 1 {
		t.Fatalf("refresh calls = %d, want 1", calls.Load())
	}
	if after, ok := source.Snapshot(); !ok || after.AccessToken != "new" {
		t.Fatal("successful refresh not published")
	}
}

func TestCachedTokenSourceSharedClientsAndReuse(t *testing.T) {
	var refreshes, apiCalls atomic.Int32
	client := httpClient(func(*http.Request) (*http.Response, error) {
		refreshes.Add(1)
		return jsonResponse(http.StatusOK, `{"access_token":"new","expires_in":3600}`), nil
	})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls.Add(1)
		if got := r.Header.Get("Authorization"); got != "Bearer new" {
			t.Errorf("Authorization = %q, want Bearer new", got)
		}
		fmt.Fprint(w, `{"used":0,"allocation":{".tag":"individual","allocated":100}}`)
	}))
	defer server.Close()
	source := dropboxoauth.NewCachedTokenSource(context.Background(), "app-key", &oauth2.Token{RefreshToken: "secret"}, dropboxoauth.WithHTTPClient(client))
	config := dropbox.Config{TokenSource: source, URLGenerator: func(_, _, _ string) string { return server.URL }}
	clients := []users.Client{users.New(config), users.New(config)}
	var wg sync.WaitGroup
	for _, dbx := range clients {
		wg.Go(func() {
			if _, err := dbx.GetSpaceUsage(); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	snapshot, ok := source.Snapshot()
	if !ok {
		t.Fatal("no snapshot after SDK calls")
	}
	reused := dropboxoauth.NewCachedTokenSource(context.Background(), "app-key", &oauth2.Token{
		AccessToken: snapshot.AccessToken, Expiry: snapshot.Expiry, RefreshToken: "secret",
	}, dropboxoauth.WithHTTPClient(client))
	config.TokenSource = reused
	if _, err := users.New(config).GetSpaceUsage(); err != nil {
		t.Fatal(err)
	}
	if got, ok := reused.Snapshot(); !ok || got != snapshot {
		t.Fatal("reused snapshot changed access token or expiration")
	}
	if refreshes.Load() != 1 || apiCalls.Load() != 3 {
		t.Fatalf("refreshes = %d, API calls = %d; want 1, 3", refreshes.Load(), apiCalls.Load())
	}
}

func TestCachedTokenSourceIsolatesOAuth2Caches(t *testing.T) {
	calls := 0
	client := httpClient(func(*http.Request) (*http.Response, error) {
		calls++
		return jsonResponse(http.StatusOK, `{"access_token":"access","expires_in":3600}`), nil
	})
	source := dropboxoauth.NewCachedTokenSource(context.Background(), "app-key",
		&oauth2.Token{RefreshToken: "secret"}, dropboxoauth.WithHTTPClient(client))
	// A client's reuse source writes its own early-expiry metadata into the
	// returned token. That must not change the shared source's cached token.
	outer := oauth2.ReuseTokenSourceWithExpiry(nil, source, 2*time.Hour)
	if _, err := outer.Token(); err != nil {
		t.Fatal(err)
	}
	if _, err := source.Token(); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("refresh calls = %d, want 1; outer cache changed shared expiry metadata", calls)
	}
}

func TestAccessTokenSnapshotTextFormatting(t *testing.T) {
	for _, value := range []struct {
		name       string
		snapshot   dropboxoauth.AccessTokenSnapshot
		expiryText string
	}{
		{"nonzero expiry", dropboxoauth.AccessTokenSnapshot{
			AccessToken: "sl.secret-token",
			Expiry:      time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC),
		}, "2030-01-02T03:04:05Z"},
		{"zero expiry", dropboxoauth.AccessTokenSnapshot{AccessToken: "sl.secret-token"}, "0001-01-01T00:00:00Z"},
		{"empty access token", dropboxoauth.AccessTokenSnapshot{}, "0001-01-01T00:00:00Z"},
		{"fractional expiry with offset", dropboxoauth.AccessTokenSnapshot{
			AccessToken: "sl.secret-token",
			Expiry:      time.Date(2030, 1, 2, 3, 4, 5, 123456789, time.FixedZone("offset", 3600)),
		}, "2030-01-02T03:04:05.123456789+01:00"},
	} {
		t.Run(value.name, func(t *testing.T) {
			snapshot := value.snapshot
			before := snapshot
			want := "AccessTokenSnapshot{AccessToken:<redacted>, Expiry:" + value.expiryText + "}"
			for _, form := range []struct {
				name   string
				render func(any) string
			}{
				{"fmt.Sprint", func(v any) string { return fmt.Sprint(v) }},
				{"fmt %v", func(v any) string { return fmt.Sprintf("%v", v) }},
				{"fmt %+v", func(v any) string { return fmt.Sprintf("%+v", v) }},
				{"fmt %#v", func(v any) string { return fmt.Sprintf("%#v", v) }},
				{"fmt %s", func(v any) string { return fmt.Sprintf("%s", v) }},
				{"log %v", func(v any) string {
					var output strings.Builder
					log.New(&output, "", 0).Printf("%v", v)
					return strings.TrimSuffix(output.String(), "\n")
				}},
				{"log %+v", func(v any) string {
					var output strings.Builder
					log.New(&output, "", 0).Printf("%+v", v)
					return strings.TrimSuffix(output.String(), "\n")
				}},
				{"log %#v", func(v any) string {
					var output strings.Builder
					log.New(&output, "", 0).Printf("%#v", v)
					return strings.TrimSuffix(output.String(), "\n")
				}},
			} {
				t.Run(form.name, func(t *testing.T) {
					for _, input := range []any{snapshot, &snapshot} {
						got := form.render(input)
						if got != want {
							t.Errorf("formatted snapshot = %q, want %q", got, want)
						}
						if strings.Contains(got, "sl.secret-token") {
							t.Error("text formatting exposes access token")
						}
						if snapshot != before {
							t.Fatal("formatting changed snapshot fields")
						}
					}
				})
			}
			data, err := json.Marshal(snapshot)
			if err != nil {
				t.Fatal(err)
			}
			var restored dropboxoauth.AccessTokenSnapshot
			if err := json.Unmarshal(data, &restored); err != nil {
				t.Fatal(err)
			}
			if restored.AccessToken != before.AccessToken || !restored.Expiry.Equal(before.Expiry) {
				t.Fatal("JSON serialization did not retain original credentials")
			}
			if snapshot != before {
				t.Fatal("JSON serialization changed snapshot fields")
			}
		})
	}
}
