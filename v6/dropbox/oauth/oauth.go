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

// Package oauth provides helpers for Dropbox OAuth 2 flows.
package oauth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"golang.org/x/oauth2"
)

const (
	tokenAccessTypeParam      = "token_access_type"
	includeGrantedScopesParam = "include_granted_scopes"
)

// TokenAccessType controls whether Dropbox returns refreshable credentials.
type TokenAccessType string

const (
	// TokenAccessTypeOffline requests a refresh token.
	TokenAccessTypeOffline TokenAccessType = "offline"
	// TokenAccessTypeOnline requests only a short-lived access token.
	TokenAccessTypeOnline TokenAccessType = "online"
)

// IncludeGrantedScopes controls whether Dropbox reuses previously granted scopes.
type IncludeGrantedScopes string

const (
	// IncludeGrantedScopesUser requests previously granted user scopes.
	IncludeGrantedScopesUser IncludeGrantedScopes = "user"
	// IncludeGrantedScopesTeam requests previously granted team scopes.
	IncludeGrantedScopesTeam IncludeGrantedScopes = "team"
)

// Option configures Dropbox OAuth helpers.
type Option func(*options)

type options struct {
	domain               string
	redirectURL          string
	state                string
	verifier             string
	scopes               []string
	tokenAccessType      TokenAccessType
	includeGrantedScopes IncludeGrantedScopes
	httpClient           *http.Client
}

// PKCEFlow manages authorization URL generation and code exchange for PKCE.
type PKCEFlow struct {
	opts options
	conf *oauth2.Config
}

// NewPKCEFlow creates a Dropbox OAuth 2 PKCE flow.
func NewPKCEFlow(appKey string, opts ...Option) (*PKCEFlow, error) {
	appKey = strings.TrimSpace(appKey)
	if appKey == "" {
		return nil, errors.New("dropbox oauth: app key is required")
	}

	o := options{
		tokenAccessType: TokenAccessTypeOffline,
	}
	applyOptions(&o, opts)
	if o.state == "" {
		o.state = oauth2.GenerateVerifier()
	}
	if o.verifier == "" {
		o.verifier = oauth2.GenerateVerifier()
	}

	return &PKCEFlow{
		opts: o,
		conf: oauthConfig(appKey, o.domain, o.redirectURL, o.scopes),
	}, nil
}

// WithDomain configures the Dropbox API domain.
func WithDomain(domain string) Option {
	return func(opts *options) {
		opts.domain = domain
	}
}

// WithRedirectURL configures the OAuth redirect URL.
func WithRedirectURL(url string) Option {
	return func(opts *options) {
		opts.redirectURL = url
	}
}

// WithState configures the OAuth state value.
func WithState(state string) Option {
	return func(opts *options) {
		opts.state = state
	}
}

// WithVerifier configures the PKCE code verifier.
func WithVerifier(verifier string) Option {
	return func(opts *options) {
		opts.verifier = verifier
	}
}

// WithScopes configures OAuth scopes.
func WithScopes(scopes ...string) Option {
	return func(opts *options) {
		opts.scopes = append([]string(nil), scopes...)
	}
}

// WithTokenAccessType configures the Dropbox token access type.
func WithTokenAccessType(value TokenAccessType) Option {
	return func(opts *options) {
		opts.tokenAccessType = value
	}
}

// WithIncludeGrantedScopes configures Dropbox include_granted_scopes.
func WithIncludeGrantedScopes(value IncludeGrantedScopes) Option {
	return func(opts *options) {
		opts.includeGrantedScopes = value
	}
}

// WithHTTPClient configures the HTTP client used for token exchange or refresh.
func WithHTTPClient(client *http.Client) Option {
	return func(opts *options) {
		opts.httpClient = client
	}
}

// AuthCodeURL returns the authorization URL for this PKCE flow.
func (f *PKCEFlow) AuthCodeURL() string {
	options := []oauth2.AuthCodeOption{oauth2.S256ChallengeOption(f.opts.verifier)}
	if f.opts.tokenAccessType != "" {
		options = append(options, oauth2.SetAuthURLParam(tokenAccessTypeParam, string(f.opts.tokenAccessType)))
	}
	if f.opts.includeGrantedScopes != "" {
		options = append(options, oauth2.SetAuthURLParam(includeGrantedScopesParam, string(f.opts.includeGrantedScopes)))
	}
	return f.conf.AuthCodeURL(f.opts.state, options...)
}

// Exchange exchanges an authorization code for an OAuth token.
func (f *PKCEFlow) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return f.conf.Exchange(withHTTPClient(ctx, f.opts.httpClient), code, oauth2.VerifierOption(f.opts.verifier))
}

// State returns the OAuth state for this flow.
func (f *PKCEFlow) State() string {
	return f.opts.state
}

// Verifier returns the PKCE verifier for this flow.
func (f *PKCEFlow) Verifier() string {
	return f.opts.verifier
}

// Refresh refreshes a Dropbox OAuth token.
func Refresh(ctx context.Context, appKey string, token *oauth2.Token, opts ...Option) (*oauth2.Token, error) {
	appKey = strings.TrimSpace(appKey)
	if appKey == "" {
		return nil, errors.New("dropbox oauth: app key is required")
	}
	if token == nil {
		return nil, errors.New("dropbox oauth: token is required")
	}
	if token.RefreshToken == "" {
		return nil, errors.New("dropbox oauth: refresh token is required")
	}

	o := options{}
	applyOptions(&o, opts)

	expired := *token
	expired.Expiry = time.Now().Add(-time.Second)

	refreshed, err := oauthConfig(appKey, o.domain, "", nil).TokenSource(withHTTPClient(ctx, o.httpClient), &expired).Token()
	if err != nil {
		return nil, err
	}
	if refreshed == nil {
		return nil, errors.New("dropbox oauth: token refresh returned nil token")
	}
	return refreshed, nil
}

// TokenSource returns an OAuth token source backed by Dropbox refresh tokens.
func TokenSource(ctx context.Context, appKey string, token *oauth2.Token, opts ...Option) oauth2.TokenSource {
	appKey = strings.TrimSpace(appKey)
	o := options{}
	applyOptions(&o, opts)
	return oauthConfig(appKey, o.domain, "", nil).TokenSource(withHTTPClient(ctx, o.httpClient), token)
}

func applyOptions(o *options, opts []Option) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}

func oauthConfig(appKey string, domain string, redirectURL string, scopes []string) *oauth2.Config {
	endpoint := dropbox.OAuthEndpoint(domain)
	endpoint.AuthStyle = oauth2.AuthStyleInParams
	return &oauth2.Config{
		ClientID:    appKey,
		Endpoint:    endpoint,
		RedirectURL: redirectURL,
		Scopes:      append([]string(nil), scopes...),
	}
}

func withHTTPClient(ctx context.Context, client *http.Client) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		return ctx
	}
	return context.WithValue(ctx, oauth2.HTTPClient, client)
}
