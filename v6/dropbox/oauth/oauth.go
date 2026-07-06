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
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"golang.org/x/oauth2"
)

const (
	tokenAccessTypeParam      = "token_access_type"
	includeGrantedScopesParam = "include_granted_scopes"
	localeParam               = "locale"
	stateSeparator            = "|"
	csrfTokenBytes            = 16
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
	locale               string
	httpClient           *http.Client
}

// PKCEFlow manages authorization URL generation and code exchange for PKCE.
type PKCEFlow struct {
	opts options
	conf *oauth2.Config
}

// WebPKCEFlow manages redirect-based authorization URL generation and code
// exchange for PKCE web apps.
type WebPKCEFlow struct {
	opts options
	conf *oauth2.Config
}

// FlowResult contains OAuth authorization information returned by a flow.
type FlowResult struct {
	Token     *oauth2.Token
	AccountID string
	TeamID    string
	UserID    string
	URLState  string
	Scopes    []string
}

// BadRequestError is returned when the redirect query is malformed.
type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	if e.Message == "" {
		return "dropbox oauth: bad request"
	}
	return "dropbox oauth: bad request: " + e.Message
}

// BadStateError is returned when the caller has no stored CSRF state.
type BadStateError struct {
	Message string
}

func (e *BadStateError) Error() string {
	if e.Message == "" {
		return "dropbox oauth: bad state"
	}
	return "dropbox oauth: bad state: " + e.Message
}

// CSRFError is returned when the redirect state does not match the stored CSRF token.
type CSRFError struct{}

func (e *CSRFError) Error() string {
	return "dropbox oauth: csrf token mismatch"
}

// NotApprovedError is returned when the user denies authorization.
type NotApprovedError struct {
	Description string
}

func (e *NotApprovedError) Error() string {
	if e.Description == "" {
		return "dropbox oauth: authorization was not approved"
	}
	return "dropbox oauth: authorization was not approved: " + e.Description
}

// ProviderError is returned when Dropbox redirects with an unexpected error.
type ProviderError struct {
	ErrorCode   string
	Description string
}

func (e *ProviderError) Error() string {
	if e.Description == "" {
		return "dropbox oauth: provider error: " + e.ErrorCode
	}
	return "dropbox oauth: provider error: " + e.ErrorCode + ": " + e.Description
}

// NewPKCEFlow creates a Dropbox OAuth 2 PKCE flow.
func NewPKCEFlow(appKey string, opts ...Option) (*PKCEFlow, error) {
	appKey, err := cleanAppKey(appKey)
	if err != nil {
		return nil, err
	}

	o := newFlowOptions(opts)
	if o.state == "" {
		o.state = oauth2.GenerateVerifier()
	}

	return &PKCEFlow{
		opts: o,
		conf: oauthConfig(appKey, o.domain, o.redirectURL, o.scopes),
	}, nil
}

// NewWebPKCEFlow creates a Dropbox OAuth 2 PKCE flow for web redirects.
func NewWebPKCEFlow(appKey string, redirectURL string, opts ...Option) (*WebPKCEFlow, error) {
	appKey, err := cleanAppKey(appKey)
	if err != nil {
		return nil, err
	}

	redirectURL = strings.TrimSpace(redirectURL)
	o := newFlowOptions(opts)
	if o.state != "" {
		return nil, errors.New("dropbox oauth: WithState is not supported by WebPKCEFlow; pass URL state to Start")
	}
	o.redirectURL = strings.TrimSpace(o.redirectURL)
	if o.redirectURL == "" {
		o.redirectURL = redirectURL
	}
	if o.redirectURL == "" {
		return nil, errors.New("dropbox oauth: redirect URL is required")
	}

	return &WebPKCEFlow{
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

// WithState configures the OAuth state value for PKCEFlow. WebPKCEFlow builds
// state from its CSRF token and the urlState passed to Start.
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

// WithLocale configures the locale shown on the authorization page.
func WithLocale(locale string) Option {
	return func(opts *options) {
		opts.locale = locale
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
	return f.conf.AuthCodeURL(f.opts.state, authCodeOptions(f.opts)...)
}

// Exchange exchanges an authorization code for an OAuth token.
func (f *PKCEFlow) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return f.conf.Exchange(withHTTPClient(ctx, f.opts.httpClient), code, exchangeOptions(f.opts)...)
}

// State returns the OAuth state for this flow.
func (f *PKCEFlow) State() string {
	return f.opts.state
}

// Verifier returns the PKCE verifier for this flow.
func (f *PKCEFlow) Verifier() string {
	return f.opts.verifier
}

// Start returns the authorization URL and CSRF token for this web PKCE flow.
func (f *WebPKCEFlow) Start(urlState string) (authURL string, csrfToken string, err error) {
	csrfToken, err = generateCSRFToken()
	if err != nil {
		return "", "", err
	}

	state := csrfToken
	if urlState != "" {
		state += stateSeparator + urlState
	}

	return f.conf.AuthCodeURL(state, authCodeOptions(f.opts)...), csrfToken, nil
}

// Finish validates a redirect query and exchanges its authorization code for an OAuth token.
func (f *WebPKCEFlow) Finish(ctx context.Context, query url.Values, csrfToken string) (*FlowResult, error) {
	state := query.Get("state")
	if state == "" {
		return nil, &BadRequestError{Message: "missing query parameter state"}
	}

	code := query.Get("code")
	errorCode := query.Get("error")
	if code != "" && errorCode != "" {
		return nil, &BadRequestError{Message: "query parameters code and error are both set"}
	}
	if code == "" && errorCode == "" {
		return nil, &BadRequestError{Message: "neither query parameter code nor error is set"}
	}

	if csrfToken == "" {
		return nil, &BadStateError{Message: "missing csrf token"}
	}

	givenCSRF, urlState := splitState(state)
	if !safeEquals(csrfToken, givenCSRF) {
		return nil, &CSRFError{}
	}

	if errorCode != "" {
		description := query.Get("error_description")
		if errorCode == "access_denied" {
			return nil, &NotApprovedError{Description: description}
		}
		return nil, &ProviderError{ErrorCode: errorCode, Description: description}
	}

	token, err := f.conf.Exchange(withHTTPClient(ctx, f.opts.httpClient), code, exchangeOptions(f.opts)...)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, errors.New("dropbox oauth: token exchange returned nil token")
	}

	return flowResultFromToken(token, urlState), nil
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

func cleanAppKey(appKey string) (string, error) {
	appKey = strings.TrimSpace(appKey)
	if appKey == "" {
		return "", errors.New("dropbox oauth: app key is required")
	}
	return appKey, nil
}

func newFlowOptions(opts []Option) options {
	o := options{
		tokenAccessType: TokenAccessTypeOffline,
	}
	applyOptions(&o, opts)
	if o.verifier == "" {
		o.verifier = oauth2.GenerateVerifier()
	}
	return o
}

func authCodeOptions(o options) []oauth2.AuthCodeOption {
	options := []oauth2.AuthCodeOption{oauth2.S256ChallengeOption(o.verifier)}
	if o.tokenAccessType != "" {
		options = append(options, oauth2.SetAuthURLParam(tokenAccessTypeParam, string(o.tokenAccessType)))
	}
	if o.includeGrantedScopes != "" {
		options = append(options, oauth2.SetAuthURLParam(includeGrantedScopesParam, string(o.includeGrantedScopes)))
	}
	if o.locale != "" {
		options = append(options, oauth2.SetAuthURLParam(localeParam, o.locale))
	}
	return options
}

func exchangeOptions(o options) []oauth2.AuthCodeOption {
	return []oauth2.AuthCodeOption{oauth2.VerifierOption(o.verifier)}
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

func generateCSRFToken() (string, error) {
	b := make([]byte, csrfTokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func splitState(state string) (csrfToken string, urlState string) {
	csrfToken, urlState, _ = strings.Cut(state, stateSeparator)
	return csrfToken, urlState
}

func safeEquals(a string, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func flowResultFromToken(token *oauth2.Token, urlState string) *FlowResult {
	return &FlowResult{
		Token:     token,
		AccountID: extraString(token, "account_id"),
		TeamID:    extraString(token, "team_id"),
		UserID:    extraString(token, "uid"),
		URLState:  urlState,
		Scopes:    strings.Fields(extraString(token, "scope")),
	}
}

func extraString(token *oauth2.Token, key string) string {
	value := token.Extra(key)
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case json.Number:
		return v.String()
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return ""
	}
}
