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
	"errors"
	"fmt"
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

func TestPKCEFlowAuthCodeURLOpenIDScopes(t *testing.T) {
	flow, err := dropboxoauth.NewPKCEFlow(
		"app-key",
		dropboxoauth.WithState("test-state"),
		dropboxoauth.WithVerifier(testVerifier),
		dropboxoauth.WithScopes(dropboxoauth.ScopeOpenID, dropboxoauth.ScopeEmail, dropboxoauth.ScopeProfile),
	)
	if err != nil {
		t.Fatal(err)
	}

	authURL, err := url.Parse(flow.AuthCodeURL())
	if err != nil {
		t.Fatal(err)
	}
	assertQueryValue(t, authURL.Query(), "scope", "openid email profile")
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

func TestGenerateVerifier(t *testing.T) {
	verifier, err := dropboxoauth.GenerateVerifier()
	if err != nil {
		t.Fatal(err)
	}
	if verifier == "" {
		t.Fatal("expected verifier")
	}
	if strings.ContainsAny(verifier, "+/=") {
		t.Fatalf("verifier is not raw URL-safe base64: %q", verifier)
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

func TestOAuth2FlowNoRedirectAppSecretStartAndFinish(t *testing.T) {
	var requestURL string
	var form url.Values
	client := httpClient(func(req *http.Request) (*http.Response, error) {
		requestURL = req.URL.String()
		var err error
		form, err = readForm(req)
		if err != nil {
			return nil, err
		}
		return jsonResponse(http.StatusOK, `{"access_token":"access-token","refresh_token":"refresh-token","token_type":"Bearer","expires_in":3600,"account_id":"dbid:account","uid":"12345","scope":"files.metadata.read files.content.write"}`), nil
	})

	flow, err := dropboxoauth.NewOAuth2FlowNoRedirect(
		"app-key",
		dropboxoauth.WithAppSecret("app-secret"),
		dropboxoauth.WithState("test-state"),
		dropboxoauth.WithScopes("files.metadata.read", "files.content.write"),
		dropboxoauth.WithIncludeGrantedScopes(dropboxoauth.IncludeGrantedScopesUser),
		dropboxoauth.WithTokenAccessType(dropboxoauth.TokenAccessTypeLegacy),
		dropboxoauth.WithLocale("en_US"),
		dropboxoauth.WithHTTPClient(client),
	)
	if err != nil {
		t.Fatal(err)
	}

	authURL, err := url.Parse(flow.Start())
	if err != nil {
		t.Fatal(err)
	}
	query := authURL.Query()
	assertQueryValue(t, query, "client_id", "app-key")
	assertQueryValue(t, query, "response_type", "code")
	assertQueryValue(t, query, "state", "test-state")
	assertQueryValue(t, query, "scope", "files.metadata.read files.content.write")
	assertQueryValue(t, query, "include_granted_scopes", "user")
	assertQueryValue(t, query, "token_access_type", "legacy")
	assertQueryValue(t, query, "locale", "en_US")
	assertNoQueryValue(t, query, "redirect_uri")
	assertNoQueryValue(t, query, "code_challenge")
	assertNoQueryValue(t, query, "code_challenge_method")

	result, err := flow.Finish(context.Background(), "auth-code")
	if err != nil {
		t.Fatal(err)
	}
	if result.Token.AccessToken != "access-token" || result.Token.RefreshToken != "refresh-token" || result.Token.TokenType != "Bearer" {
		t.Fatalf("unexpected token: %#v", result.Token)
	}
	if result.Token.Expiry.IsZero() {
		t.Fatal("expected token expiry")
	}
	if result.AccountID != "dbid:account" || result.UserID != "12345" {
		t.Fatalf("unexpected account info: %#v", result)
	}
	assertScopes(t, result.Scopes, "files.metadata.read", "files.content.write")

	if requestURL != "https://api.dropboxapi.com/1/oauth2/token" {
		t.Fatalf("request URL = %q", requestURL)
	}
	assertQueryValue(t, form, "grant_type", "authorization_code")
	assertQueryValue(t, form, "code", "auth-code")
	assertQueryValue(t, form, "client_id", "app-key")
	assertQueryValue(t, form, "client_secret", "app-secret")
	assertQueryValue(t, form, "locale", "en_US")
	assertNoQueryValue(t, form, "redirect_uri")
	assertNoQueryValue(t, form, "code_verifier")
}

func TestOAuth2FlowNoRedirectPKCEStartAndFinish(t *testing.T) {
	var form url.Values
	client := httpClient(func(req *http.Request) (*http.Response, error) {
		var err error
		form, err = readForm(req)
		if err != nil {
			return nil, err
		}
		return jsonResponse(http.StatusOK, `{"access_token":"access-token","token_type":"Bearer"}`), nil
	})

	flow, err := dropboxoauth.NewOAuth2FlowNoRedirect(
		"app-key",
		dropboxoauth.WithPKCE(),
		dropboxoauth.WithVerifier(testVerifier),
		dropboxoauth.WithHTTPClient(client),
	)
	if err != nil {
		t.Fatal(err)
	}

	authURL, err := url.Parse(flow.Start())
	if err != nil {
		t.Fatal(err)
	}
	query := authURL.Query()
	assertQueryValue(t, query, "code_challenge", "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM")
	assertQueryValue(t, query, "code_challenge_method", "S256")

	if _, err := flow.Finish(context.Background(), "auth-code"); err != nil {
		t.Fatal(err)
	}
	assertQueryValue(t, form, "grant_type", "authorization_code")
	assertQueryValue(t, form, "code", "auth-code")
	assertQueryValue(t, form, "client_id", "app-key")
	assertQueryValue(t, form, "code_verifier", testVerifier)
	assertNoQueryValue(t, form, "client_secret")
}

func TestOAuth2FlowNoRedirectAppSecretOmitsTokenAccessTypeByDefault(t *testing.T) {
	flow, err := dropboxoauth.NewOAuth2FlowNoRedirect(
		"app-key",
		dropboxoauth.WithAppSecret("app-secret"),
	)
	if err != nil {
		t.Fatal(err)
	}

	authURL, err := url.Parse(flow.Start())
	if err != nil {
		t.Fatal(err)
	}
	assertNoQueryValue(t, authURL.Query(), "token_access_type")
}

func TestOAuth2FlowAppSecretStartAndFinish(t *testing.T) {
	var form url.Values
	client := httpClient(func(req *http.Request) (*http.Response, error) {
		var err error
		form, err = readForm(req)
		if err != nil {
			return nil, err
		}
		return jsonResponse(http.StatusOK, `{"access_token":"access-token","refresh_token":"refresh-token","token_type":"Bearer","expires_in":3600,"team_id":"dbtid:team","uid":12345}`), nil
	})

	flow, err := dropboxoauth.NewOAuth2Flow(
		"app-key",
		"http://localhost/callback",
		dropboxoauth.WithAppSecret("app-secret"),
		dropboxoauth.WithLocale("en_US"),
		dropboxoauth.WithHTTPClient(client),
	)
	if err != nil {
		t.Fatal(err)
	}

	rawAuthURL, csrfToken, err := flow.Start("url-state")
	if err != nil {
		t.Fatal(err)
	}
	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		t.Fatal(err)
	}
	query := authURL.Query()
	assertQueryValue(t, query, "redirect_uri", "http://localhost/callback")
	assertQueryValue(t, query, "state", csrfToken+"|url-state")
	assertQueryValue(t, query, "locale", "en_US")
	assertNoQueryValue(t, query, "code_challenge")
	assertNoQueryValue(t, query, "code_challenge_method")

	result, err := flow.Finish(context.Background(), url.Values{
		"state": {csrfToken + "|url-state"},
		"code":  {"auth-code"},
	}, csrfToken)
	if err != nil {
		t.Fatal(err)
	}
	if result.Token.AccessToken != "access-token" || result.Token.RefreshToken != "refresh-token" {
		t.Fatalf("unexpected token: %#v", result.Token)
	}
	if result.AccountID != "dbtid:team" || result.TeamID != "dbtid:team" || result.UserID != "12345" {
		t.Fatalf("unexpected account info: %#v", result)
	}
	if result.URLState != "url-state" {
		t.Fatalf("url state = %q, want url-state", result.URLState)
	}

	assertQueryValue(t, form, "grant_type", "authorization_code")
	assertQueryValue(t, form, "code", "auth-code")
	assertQueryValue(t, form, "client_id", "app-key")
	assertQueryValue(t, form, "client_secret", "app-secret")
	assertQueryValue(t, form, "redirect_uri", "http://localhost/callback")
	assertQueryValue(t, form, "locale", "en_US")
	assertNoQueryValue(t, form, "code_verifier")
}

func TestWebPKCEFlowStart(t *testing.T) {
	flow, err := dropboxoauth.NewWebPKCEFlow(
		"app-key",
		"http://localhost/callback",
		dropboxoauth.WithVerifier(testVerifier),
		dropboxoauth.WithScopes("files.metadata.read", "files.content.write"),
		dropboxoauth.WithIncludeGrantedScopes(dropboxoauth.IncludeGrantedScopesUser),
		dropboxoauth.WithLocale("en_US"),
	)
	if err != nil {
		t.Fatal(err)
	}

	rawAuthURL, csrfToken, err := flow.Start("url-state")
	if err != nil {
		t.Fatal(err)
	}
	if csrfToken == "" {
		t.Fatal("expected csrf token")
	}

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		t.Fatal(err)
	}
	query := authURL.Query()
	if authURL.Scheme != "https" || authURL.Host != "www.dropbox.com" || authURL.Path != "/1/oauth2/authorize" {
		t.Fatalf("auth URL = %s", authURL)
	}
	assertQueryValue(t, query, "client_id", "app-key")
	assertQueryValue(t, query, "response_type", "code")
	assertQueryValue(t, query, "redirect_uri", "http://localhost/callback")
	assertQueryValue(t, query, "state", csrfToken+"|url-state")
	assertQueryValue(t, query, "scope", "files.metadata.read files.content.write")
	assertQueryValue(t, query, "include_granted_scopes", "user")
	assertQueryValue(t, query, "token_access_type", "offline")
	assertQueryValue(t, query, "code_challenge", "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM")
	assertQueryValue(t, query, "code_challenge_method", "S256")
	assertQueryValue(t, query, "locale", "en_US")
}

func TestWebPKCEFlowStartGeneratesVerifierAndCSRFOnlyState(t *testing.T) {
	flow, err := dropboxoauth.NewWebPKCEFlow("app-key", "http://localhost/callback")
	if err != nil {
		t.Fatal(err)
	}

	rawAuthURL, csrfToken, err := flow.Start("")
	if err != nil {
		t.Fatal(err)
	}

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		t.Fatal(err)
	}
	query := authURL.Query()
	assertQueryValue(t, query, "state", csrfToken)
	if got := query.Get("code_challenge"); got == "" || got == "47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU" {
		t.Fatalf("code_challenge = %q, want generated verifier challenge", got)
	}
	assertQueryValue(t, query, "code_challenge_method", "S256")
}

func TestWebPKCEFlowFinishExchangesCodeAndReturnsResult(t *testing.T) {
	var requestURL string
	var form url.Values
	client := httpClient(func(req *http.Request) (*http.Response, error) {
		requestURL = req.URL.String()
		var err error
		form, err = readForm(req)
		if err != nil {
			return nil, err
		}
		return jsonResponse(http.StatusOK, `{"access_token":"access-token","refresh_token":"refresh-token","token_type":"Bearer","expires_in":3600,"account_id":"dbid:account","team_id":"dbtid:team","uid":12345,"scope":"files.metadata.read files.content.write"}`), nil
	})

	flow, err := dropboxoauth.NewWebPKCEFlow(
		"app-key",
		"http://localhost/callback",
		dropboxoauth.WithVerifier(testVerifier),
		dropboxoauth.WithHTTPClient(client),
		dropboxoauth.WithLocale("en_US"),
	)
	if err != nil {
		t.Fatal(err)
	}

	result, err := flow.Finish(context.Background(), url.Values{
		"state": {"stored-csrf|url-state"},
		"code":  {"auth-code"},
	}, "stored-csrf")
	if err != nil {
		t.Fatal(err)
	}
	if result.Token.AccessToken != "access-token" || result.Token.RefreshToken != "refresh-token" || result.Token.TokenType != "Bearer" {
		t.Fatalf("unexpected token: %#v", result.Token)
	}
	if result.Token.Expiry.IsZero() {
		t.Fatal("expected token expiry")
	}
	if result.AccountID != "dbid:account" || result.TeamID != "dbtid:team" || result.UserID != "12345" {
		t.Fatalf("unexpected account info: %#v", result)
	}
	if result.URLState != "url-state" {
		t.Fatalf("url state = %q, want url-state", result.URLState)
	}
	assertScopes(t, result.Scopes, "files.metadata.read", "files.content.write")

	if requestURL != "https://api.dropboxapi.com/1/oauth2/token" {
		t.Fatalf("request URL = %q", requestURL)
	}
	assertQueryValue(t, form, "grant_type", "authorization_code")
	assertQueryValue(t, form, "code", "auth-code")
	assertQueryValue(t, form, "code_verifier", testVerifier)
	assertQueryValue(t, form, "client_id", "app-key")
	assertQueryValue(t, form, "redirect_uri", "http://localhost/callback")
	assertQueryValue(t, form, "locale", "en_US")
	if got := form.Get("client_secret"); got != "" {
		t.Fatalf("client_secret = %q, want empty", got)
	}
}

func TestWebPKCEFlowStartUsesCustomDomainAndRedirectFallback(t *testing.T) {
	flow, err := dropboxoauth.NewWebPKCEFlow(
		"app-key",
		"http://localhost/callback",
		dropboxoauth.WithDomain(".example.com"),
		dropboxoauth.WithRedirectURL(""),
	)
	if err != nil {
		t.Fatal(err)
	}

	rawAuthURL, _, err := flow.Start("")
	if err != nil {
		t.Fatal(err)
	}

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		t.Fatal(err)
	}
	if authURL.Host != "meta.example.com" {
		t.Fatalf("auth host = %q, want meta.example.com", authURL.Host)
	}
	assertQueryValue(t, authURL.Query(), "redirect_uri", "http://localhost/callback")
}

func TestWebPKCEFlowFinishRejectsRedirectErrors(t *testing.T) {
	tests := []struct {
		name      string
		query     url.Values
		csrfToken string
		wantErr   error
	}{
		{
			name: "missing state",
			query: url.Values{
				"code": {"auth-code"},
			},
			csrfToken: "stored-csrf",
			wantErr:   &dropboxoauth.BadRequestError{},
		},
		{
			name: "missing code and error",
			query: url.Values{
				"state": {"stored-csrf"},
			},
			csrfToken: "stored-csrf",
			wantErr:   &dropboxoauth.BadRequestError{},
		},
		{
			name: "code and error",
			query: url.Values{
				"state": {"stored-csrf"},
				"code":  {"auth-code"},
				"error": {"access_denied"},
			},
			csrfToken: "stored-csrf",
			wantErr:   &dropboxoauth.BadRequestError{},
		},
		{
			name: "missing csrf token",
			query: url.Values{
				"state": {"stored-csrf"},
				"code":  {"auth-code"},
			},
			wantErr: &dropboxoauth.BadStateError{},
		},
		{
			name: "csrf mismatch",
			query: url.Values{
				"state": {"other-csrf"},
				"code":  {"auth-code"},
			},
			csrfToken: "stored-csrf",
			wantErr:   &dropboxoauth.CSRFError{},
		},
		{
			name: "access denied",
			query: url.Values{
				"state":             {"stored-csrf"},
				"error":             {"access_denied"},
				"error_description": {"denied"},
			},
			csrfToken: "stored-csrf",
			wantErr:   &dropboxoauth.NotApprovedError{},
		},
		{
			name: "provider error",
			query: url.Values{
				"state":             {"stored-csrf"},
				"error":             {"server_error"},
				"error_description": {"bad"},
			},
			csrfToken: "stored-csrf",
			wantErr:   &dropboxoauth.ProviderError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flow, err := dropboxoauth.NewWebPKCEFlow("app-key", "http://localhost/callback")
			if err != nil {
				t.Fatal(err)
			}

			_, err = flow.Finish(context.Background(), tt.query, tt.csrfToken)
			if err == nil {
				t.Fatal("expected error")
			}
			if _, ok := tt.wantErr.(*dropboxoauth.CSRFError); ok && strings.Contains(fmt.Sprintf("%+v", err), tt.csrfToken) {
				t.Fatalf("csrf error leaked stored token: %#v", err)
			}
			assertErrorAs(t, err, tt.wantErr)
		})
	}
}

func TestNewWebPKCEFlowRejectsMissingInputs(t *testing.T) {
	if _, err := dropboxoauth.NewWebPKCEFlow(" ", "http://localhost/callback"); err == nil {
		t.Fatal("expected missing app key error")
	}
	if _, err := dropboxoauth.NewWebPKCEFlow("app-key", " "); err == nil {
		t.Fatal("expected missing redirect URL error")
	}
	if _, err := dropboxoauth.NewWebPKCEFlow("app-key", "http://localhost/callback", dropboxoauth.WithState("state")); err == nil {
		t.Fatal("expected unsupported state error")
	}
	if _, err := dropboxoauth.NewWebPKCEFlow("app-key", "http://localhost/callback", dropboxoauth.WithAppSecret("secret")); err == nil {
		t.Fatal("expected unsupported app secret error")
	}
}

func TestNewPKCEFlowRejectsInvalidOptions(t *testing.T) {
	tests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "invalid token access type",
			fn: func() error {
				_, err := dropboxoauth.NewPKCEFlow(
					"app-key",
					dropboxoauth.WithTokenAccessType(dropboxoauth.TokenAccessType("invalid")),
				)
				return err
			},
		},
		{
			name: "empty scopes",
			fn: func() error {
				_, err := dropboxoauth.NewPKCEFlow("app-key", dropboxoauth.WithScopes())
				return err
			},
		},
		{
			name: "include granted scopes without scopes",
			fn: func() error {
				_, err := dropboxoauth.NewPKCEFlow(
					"app-key",
					dropboxoauth.WithIncludeGrantedScopes(dropboxoauth.IncludeGrantedScopesUser),
				)
				return err
			},
		},
		{
			name: "invalid include granted scopes",
			fn: func() error {
				_, err := dropboxoauth.NewPKCEFlow(
					"app-key",
					dropboxoauth.WithScopes("files.metadata.read"),
					dropboxoauth.WithIncludeGrantedScopes(dropboxoauth.IncludeGrantedScopes("invalid")),
				)
				return err
			},
		},
		{
			name: "web invalid option",
			fn: func() error {
				_, err := dropboxoauth.NewWebPKCEFlow(
					"app-key",
					"http://localhost/callback",
					dropboxoauth.WithTokenAccessType(dropboxoauth.TokenAccessType("invalid")),
				)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestNewOAuth2FlowRejectsInvalidOptions(t *testing.T) {
	tests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "missing secret or pkce",
			fn: func() error {
				_, err := dropboxoauth.NewOAuth2FlowNoRedirect("app-key")
				return err
			},
		},
		{
			name: "invalid token access type",
			fn: func() error {
				_, err := dropboxoauth.NewOAuth2FlowNoRedirect(
					"app-key",
					dropboxoauth.WithAppSecret("secret"),
					dropboxoauth.WithTokenAccessType(dropboxoauth.TokenAccessType("invalid")),
				)
				return err
			},
		},
		{
			name: "empty scopes",
			fn: func() error {
				_, err := dropboxoauth.NewOAuth2FlowNoRedirect(
					"app-key",
					dropboxoauth.WithAppSecret("secret"),
					dropboxoauth.WithScopes(),
				)
				return err
			},
		},
		{
			name: "include granted scopes without scopes",
			fn: func() error {
				_, err := dropboxoauth.NewOAuth2FlowNoRedirect(
					"app-key",
					dropboxoauth.WithAppSecret("secret"),
					dropboxoauth.WithIncludeGrantedScopes(dropboxoauth.IncludeGrantedScopesUser),
				)
				return err
			},
		},
		{
			name: "invalid include granted scopes",
			fn: func() error {
				_, err := dropboxoauth.NewOAuth2FlowNoRedirect(
					"app-key",
					dropboxoauth.WithAppSecret("secret"),
					dropboxoauth.WithScopes("files.metadata.read"),
					dropboxoauth.WithIncludeGrantedScopes(dropboxoauth.IncludeGrantedScopes("invalid")),
				)
				return err
			},
		},
		{
			name: "web missing redirect url",
			fn: func() error {
				_, err := dropboxoauth.NewOAuth2Flow("app-key", " ", dropboxoauth.WithAppSecret("secret"))
				return err
			},
		},
		{
			name: "web state option",
			fn: func() error {
				_, err := dropboxoauth.NewOAuth2Flow(
					"app-key",
					"http://localhost/callback",
					dropboxoauth.WithAppSecret("secret"),
					dropboxoauth.WithState("state"),
				)
				return err
			},
		},
		{
			name: "pkce constructor app secret",
			fn: func() error {
				_, err := dropboxoauth.NewPKCEFlow("app-key", dropboxoauth.WithAppSecret("secret"))
				return err
			},
		},
		{
			name: "app secret with pkce",
			fn: func() error {
				_, err := dropboxoauth.NewOAuth2FlowNoRedirect(
					"app-key",
					dropboxoauth.WithAppSecret("secret"),
					dropboxoauth.WithPKCE(),
				)
				return err
			},
		},
		{
			name: "app secret with verifier",
			fn: func() error {
				_, err := dropboxoauth.NewOAuth2FlowNoRedirect(
					"app-key",
					dropboxoauth.WithAppSecret("secret"),
					dropboxoauth.WithVerifier(testVerifier),
				)
				return err
			},
		},
		{
			name: "web app secret with pkce",
			fn: func() error {
				_, err := dropboxoauth.NewOAuth2Flow(
					"app-key",
					"http://localhost/callback",
					dropboxoauth.WithAppSecret("secret"),
					dropboxoauth.WithPKCE(),
				)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); err == nil {
				t.Fatal("expected error")
			}
		})
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

func TestRefreshSendsAppSecret(t *testing.T) {
	var form url.Values
	client := httpClient(func(req *http.Request) (*http.Response, error) {
		var err error
		form, err = readForm(req)
		if err != nil {
			return nil, err
		}
		return jsonResponse(http.StatusOK, `{"access_token":"new-access","expires_in":3600}`), nil
	})

	_, err := dropboxoauth.Refresh(context.Background(), "app-key", &oauth2.Token{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		Expiry:       time.Now().Add(-time.Hour),
	}, dropboxoauth.WithAppSecret("app-secret"), dropboxoauth.WithHTTPClient(client))
	if err != nil {
		t.Fatal(err)
	}
	assertQueryValue(t, form, "grant_type", "refresh_token")
	assertQueryValue(t, form, "refresh_token", "old-refresh")
	assertQueryValue(t, form, "client_id", "app-key")
	assertQueryValue(t, form, "client_secret", "app-secret")
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

func assertNoQueryValue(t *testing.T, values url.Values, key string) {
	t.Helper()
	if got, ok := values[key]; ok {
		t.Fatalf("%s = %#v, want absent", key, got)
	}
}

func assertScopes(t *testing.T, got []string, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("scopes = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("scopes = %#v, want %#v", got, want)
		}
	}
}

func assertErrorAs(t *testing.T, err error, want error) {
	t.Helper()
	switch want.(type) {
	case *dropboxoauth.BadRequestError:
		var target *dropboxoauth.BadRequestError
		if !errors.As(err, &target) {
			t.Fatalf("error = %T %v, want %T", err, err, want)
		}
	case *dropboxoauth.BadStateError:
		var target *dropboxoauth.BadStateError
		if !errors.As(err, &target) {
			t.Fatalf("error = %T %v, want %T", err, err, want)
		}
	case *dropboxoauth.CSRFError:
		var target *dropboxoauth.CSRFError
		if !errors.As(err, &target) {
			t.Fatalf("error = %T %v, want %T", err, err, want)
		}
	case *dropboxoauth.NotApprovedError:
		var target *dropboxoauth.NotApprovedError
		if !errors.As(err, &target) {
			t.Fatalf("error = %T %v, want %T", err, err, want)
		}
	case *dropboxoauth.ProviderError:
		var target *dropboxoauth.ProviderError
		if !errors.As(err, &target) {
			t.Fatalf("error = %T %v, want %T", err, err, want)
		}
	default:
		t.Fatalf("unhandled wanted error type %T", want)
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
