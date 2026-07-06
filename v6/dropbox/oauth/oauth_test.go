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
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	dropboxoauth "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/oauth"
	"golang.org/x/oauth2"
)

const testVerifier = "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"

func TestPKCEFlowAuthCodeURL(t *testing.T) {
	flow, err := dropboxoauth.NewPKCEFlow(
		"app-key",
		dropboxoauth.WithState("test-state"),
		dropboxoauth.WithVerifier(testVerifier),
		dropboxoauth.WithRedirectURL("http://localhost/callback"),
		dropboxoauth.WithScopes("files.metadata.read", "files.content.write"),
		dropboxoauth.WithIncludeGrantedScopes(dropboxoauth.IncludeGrantedScopesUser),
	)
	if err != nil {
		t.Fatal(err)
	}

	authURL, err := url.Parse(flow.AuthCodeURL())
	if err != nil {
		t.Fatal(err)
	}
	query := authURL.Query()
	if authURL.Scheme != "https" || authURL.Host != "www.dropbox.com" || authURL.Path != "/1/oauth2/authorize" {
		t.Fatalf("auth URL = %s", authURL)
	}
	assertQueryValue(t, query, "client_id", "app-key")
	assertQueryValue(t, query, "response_type", "code")
	assertQueryValue(t, query, "state", "test-state")
	assertQueryValue(t, query, "redirect_uri", "http://localhost/callback")
	assertQueryValue(t, query, "scope", "files.metadata.read files.content.write")
	assertQueryValue(t, query, "include_granted_scopes", "user")
	assertQueryValue(t, query, "token_access_type", "offline")
	assertQueryValue(t, query, "code_challenge", "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM")
	assertQueryValue(t, query, "code_challenge_method", "S256")
}

func TestPKCEFlowGeneratesStateAndVerifier(t *testing.T) {
	flow, err := dropboxoauth.NewPKCEFlow("app-key")
	if err != nil {
		t.Fatal(err)
	}

	if flow.State() == "" {
		t.Fatal("expected generated state")
	}
	if flow.Verifier() == "" {
		t.Fatal("expected generated verifier")
	}
	if !strings.Contains(flow.AuthCodeURL(), "token_access_type=offline") {
		t.Fatalf("expected offline token access type in auth URL, got %q", flow.AuthCodeURL())
	}
}

func TestPKCEFlowExchangeSendsVerifierAndAppKey(t *testing.T) {
	var requestURL string
	var form url.Values
	client := httpClient(func(req *http.Request) (*http.Response, error) {
		requestURL = req.URL.String()
		var err error
		form, err = readForm(req)
		if err != nil {
			return nil, err
		}
		return jsonResponse(http.StatusOK, `{"access_token":"access-token","refresh_token":"refresh-token","token_type":"Bearer","expires_in":3600}`), nil
	})

	flow, err := dropboxoauth.NewPKCEFlow(
		"app-key",
		dropboxoauth.WithVerifier(testVerifier),
		dropboxoauth.WithHTTPClient(client),
	)
	if err != nil {
		t.Fatal(err)
	}

	token, err := flow.Exchange(context.Background(), "auth-code")
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken != "access-token" || token.RefreshToken != "refresh-token" || token.TokenType != "Bearer" {
		t.Fatalf("unexpected token: %#v", token)
	}
	if requestURL != "https://api.dropboxapi.com/1/oauth2/token" {
		t.Fatalf("request URL = %q", requestURL)
	}
	assertQueryValue(t, form, "grant_type", "authorization_code")
	assertQueryValue(t, form, "code", "auth-code")
	assertQueryValue(t, form, "code_verifier", testVerifier)
	assertQueryValue(t, form, "client_id", "app-key")
	if got := form.Get("client_secret"); got != "" {
		t.Fatalf("client_secret = %q, want empty", got)
	}
}

func TestRefreshSendsRefreshTokenAndPreservesRefreshToken(t *testing.T) {
	var form url.Values
	client := httpClient(func(req *http.Request) (*http.Response, error) {
		var err error
		form, err = readForm(req)
		if err != nil {
			return nil, err
		}
		return jsonResponse(http.StatusOK, `{"access_token":"new-access","expires_in":3600}`), nil
	})

	token, err := dropboxoauth.Refresh(context.Background(), "app-key", &oauth2.Token{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(-time.Hour),
	}, dropboxoauth.WithHTTPClient(client))
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken != "new-access" {
		t.Fatalf("access token = %q, want new-access", token.AccessToken)
	}
	if token.RefreshToken != "old-refresh" {
		t.Fatalf("refresh token = %q, want old-refresh", token.RefreshToken)
	}
	if token.Type() != "Bearer" {
		t.Fatalf("token type = %q, want Bearer", token.Type())
	}
	assertQueryValue(t, form, "grant_type", "refresh_token")
	assertQueryValue(t, form, "refresh_token", "old-refresh")
	assertQueryValue(t, form, "client_id", "app-key")
	if got := form.Get("client_secret"); got != "" {
		t.Fatalf("client_secret = %q, want empty", got)
	}
}

func TestTokenSourceRefreshesExpiredToken(t *testing.T) {
	calls := 0
	client := httpClient(func(req *http.Request) (*http.Response, error) {
		calls++
		return jsonResponse(http.StatusOK, `{"access_token":"new-access","refresh_token":"new-refresh","token_type":"Bearer","expires_in":3600}`), nil
	})
	source := dropboxoauth.TokenSource(context.Background(), "app-key", &oauth2.Token{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		Expiry:       time.Now().Add(-time.Hour),
	}, dropboxoauth.WithHTTPClient(client))

	token, err := source.Token()
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken != "new-access" || token.RefreshToken != "new-refresh" {
		t.Fatalf("unexpected token: %#v", token)
	}
	cached, err := source.Token()
	if err != nil {
		t.Fatal(err)
	}
	if cached.AccessToken != "new-access" {
		t.Fatalf("cached access token = %q, want new-access", cached.AccessToken)
	}
	if calls != 1 {
		t.Fatalf("refresh calls = %d, want 1", calls)
	}
}

func TestCustomDomainChangesOAuthEndpoints(t *testing.T) {
	flow, err := dropboxoauth.NewPKCEFlow(
		"app-key",
		dropboxoauth.WithDomain(".example.com"),
		dropboxoauth.WithState("state"),
		dropboxoauth.WithVerifier(testVerifier),
	)
	if err != nil {
		t.Fatal(err)
	}

	authURL, err := url.Parse(flow.AuthCodeURL())
	if err != nil {
		t.Fatal(err)
	}
	if authURL.Host != "meta.example.com" {
		t.Fatalf("auth host = %q, want meta.example.com", authURL.Host)
	}

	var tokenHost string
	client := httpClient(func(req *http.Request) (*http.Response, error) {
		tokenHost = req.URL.Host
		return jsonResponse(http.StatusOK, `{"access_token":"new-access","refresh_token":"new-refresh","token_type":"Bearer"}`), nil
	})
	_, err = dropboxoauth.Refresh(context.Background(), "app-key", &oauth2.Token{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		Expiry:       time.Now().Add(-time.Hour),
	}, dropboxoauth.WithDomain(".example.com"), dropboxoauth.WithHTTPClient(client))
	if err != nil {
		t.Fatal(err)
	}
	if tokenHost != "api.example.com" {
		t.Fatalf("token host = %q, want api.example.com", tokenHost)
	}
}

func TestNewPKCEFlowRejectsMissingAppKey(t *testing.T) {
	if _, err := dropboxoauth.NewPKCEFlow(" "); err == nil {
		t.Fatal("expected missing app key error")
	}
}

func TestRefreshRejectsMissingRefreshToken(t *testing.T) {
	_, err := dropboxoauth.Refresh(context.Background(), "app-key", &oauth2.Token{AccessToken: "access-token"})
	if err == nil {
		t.Fatal("expected missing refresh token error")
	}
}

func assertQueryValue(t *testing.T, values url.Values, key string, want string) {
	t.Helper()
	if got := values.Get(key); got != want {
		t.Fatalf("%s = %q, want %q", key, got, want)
	}
}

func readForm(req *http.Request) (url.Values, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	return url.ParseQuery(string(body))
}

func httpClient(fn func(*http.Request) (*http.Response, error)) *http.Client {
	return &http.Client{Transport: roundTripFunc(fn)}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
