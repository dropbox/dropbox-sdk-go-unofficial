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
	// TokenAccessTypeLegacy requests a long-lived access token.
	//
	// Deprecated: prefer TokenAccessTypeOffline and refresh tokens.
	TokenAccessTypeLegacy TokenAccessType = "legacy"
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
	appSecret            string
	redirectURL          string
	state                string
	verifier             string
	usePKCE              bool
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

// OAuth2FlowNoRedirect manages authorization URL generation and code exchange
// for OAuth 2 apps that do not use a redirect URI.
type OAuth2FlowNoRedirect struct {
	opts options
	conf *oauth2.Config
}

// OAuth2Flow manages redirect-based authorization URL generation and code
// exchange for OAuth 2 web apps.
type OAuth2Flow struct {
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

	o, err := newPKCEOptions(opts)
	if err != nil {
		return nil, err
	}
	if o.appSecret != "" {
		return nil, errors.New("dropbox oauth: WithAppSecret is not supported by PKCEFlow")
	}
	if o.state == "" {
		o.state, err = generateVerifier()
		if err != nil {
			return nil, err
		}
	}

	return &PKCEFlow{
		opts: o,
		conf: oauthConfig(appKey, "", o.domain, o.redirectURL, o.scopes),
	}, nil
}

// NewOAuth2FlowNoRedirect creates a Dropbox OAuth 2 flow for apps that do not
// use a redirect URI. Pass WithAppSecret for confidential apps or WithPKCE for
// public clients.
func NewOAuth2FlowNoRedirect(appKey string, opts ...Option) (*OAuth2FlowNoRedirect, error) {
	appKey, err := cleanAppKey(appKey)
	if err != nil {
		return nil, err
	}

	o, err := newOAuth2Options(opts)
	if err != nil {
		return nil, err
	}

	return &OAuth2FlowNoRedirect{
		opts: o,
		conf: oauthConfig(appKey, flowAppSecret(o), o.domain, "", o.scopes),
	}, nil
}

// NewOAuth2Flow creates a Dropbox OAuth 2 flow for web redirects. Pass
// WithAppSecret for confidential apps or WithPKCE for public clients.
func NewOAuth2Flow(appKey string, redirectURL string, opts ...Option) (*OAuth2Flow, error) {
	appKey, err := cleanAppKey(appKey)
	if err != nil {
		return nil, err
	}

	redirectURL = strings.TrimSpace(redirectURL)
	o, err := newOAuth2Options(opts)
	if err != nil {
		return nil, err
	}
	if o.state != "" {
		return nil, errors.New("dropbox oauth: WithState is not supported by OAuth2Flow; pass URL state to Start")
	}
	o.redirectURL = strings.TrimSpace(o.redirectURL)
	if o.redirectURL == "" {
		o.redirectURL = redirectURL
	}
	if o.redirectURL == "" {
		return nil, errors.New("dropbox oauth: redirect URL is required")
	}

	return &OAuth2Flow{
		opts: o,
		conf: oauthConfig(appKey, flowAppSecret(o), o.domain, o.redirectURL, o.scopes),
	}, nil
}

// NewWebPKCEFlow creates a Dropbox OAuth 2 PKCE flow for web redirects.
func NewWebPKCEFlow(appKey string, redirectURL string, opts ...Option) (*WebPKCEFlow, error) {
	appKey, err := cleanAppKey(appKey)
	if err != nil {
		return nil, err
	}

	redirectURL = strings.TrimSpace(redirectURL)
	o, err := newPKCEOptions(opts)
	if err != nil {
		return nil, err
	}
	if o.appSecret != "" {
		return nil, errors.New("dropbox oauth: WithAppSecret is not supported by WebPKCEFlow")
	}
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
		conf: oauthConfig(appKey, "", o.domain, o.redirectURL, o.scopes),
	}, nil
}

// WithDomain configures the Dropbox API domain.
func WithDomain(domain string) Option {
	return func(opts *options) {
		opts.domain = domain
	}
}

// WithAppSecret configures the Dropbox API app secret for confidential OAuth apps.
func WithAppSecret(secret string) Option {
	return func(opts *options) {
		opts.appSecret = strings.TrimSpace(secret)
	}
}

// WithPKCE configures general OAuth flows to use PKCE instead of an app secret.
func WithPKCE() Option {
	return func(opts *options) {
		opts.usePKCE = true
	}
}

// WithRedirectURL configures the OAuth redirect URL.
func WithRedirectURL(url string) Option {
	return func(opts *options) {
		opts.redirectURL = url
	}
}

// WithState configures the OAuth state value for no-redirect flows. Web flows
// build state from their CSRF token and the urlState passed to Start.
func WithState(state string) Option {
	return func(opts *options) {
		opts.state = state
	}
}

// WithVerifier configures the PKCE code verifier.
func WithVerifier(verifier string) Option {
	return func(opts *options) {
		opts.verifier = verifier
		opts.usePKCE = true
	}
}

// WithScopes configures OAuth scopes.
func WithScopes(scopes ...string) Option {
	return func(opts *options) {
		opts.scopes = append(make([]string, 0, len(scopes)), scopes...)
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

// WithLocale configures the locale sent to Dropbox during OAuth.
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

// Start returns the authorization URL for this no-redirect OAuth flow.
func (f *OAuth2FlowNoRedirect) Start() string {
	return f.conf.AuthCodeURL(f.opts.state, authCodeOptions(f.opts)...)
}

// Finish exchanges an authorization code for OAuth authorization information.
func (f *OAuth2FlowNoRedirect) Finish(ctx context.Context, code string) (*FlowResult, error) {
	token, err := f.conf.Exchange(withHTTPClient(ctx, f.opts.httpClient), code, exchangeOptions(f.opts)...)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, errors.New("dropbox oauth: token exchange returned nil token")
	}
	return flowResultFromToken(token, ""), nil
}

// Start returns the authorization URL and CSRF token for this web OAuth flow.
func (f *OAuth2Flow) Start(urlState string) (authURL string, csrfToken string, err error) {
	return startWebFlow(f.conf, f.opts, urlState)
}

// Finish validates a redirect query and exchanges its authorization code for
// OAuth authorization information.
func (f *OAuth2Flow) Finish(ctx context.Context, query url.Values, csrfToken string) (*FlowResult, error) {
	return finishWebFlow(ctx, f.conf, f.opts, query, csrfToken)
}

// Start returns the authorization URL and CSRF token for this web PKCE flow.
func (f *WebPKCEFlow) Start(urlState string) (authURL string, csrfToken string, err error) {
	return startWebFlow(f.conf, f.opts, urlState)
}

func startWebFlow(conf *oauth2.Config, opts options, urlState string) (authURL string, csrfToken string, err error) {
	csrfToken, err = generateCSRFToken()
	if err != nil {
		return "", "", err
	}

	state := csrfToken
	if urlState != "" {
		state += stateSeparator + urlState
	}

	return conf.AuthCodeURL(state, authCodeOptions(opts)...), csrfToken, nil
}

// Finish validates a redirect query and exchanges its authorization code for an OAuth token.
func (f *WebPKCEFlow) Finish(ctx context.Context, query url.Values, csrfToken string) (*FlowResult, error) {
	return finishWebFlow(ctx, f.conf, f.opts, query, csrfToken)
}

func finishWebFlow(ctx context.Context, conf *oauth2.Config, opts options, query url.Values, csrfToken string) (*FlowResult, error) {
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

	token, err := conf.Exchange(withHTTPClient(ctx, opts.httpClient), code, exchangeOptions(opts)...)
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

	refreshed, err := oauthConfig(appKey, o.appSecret, o.domain, "", nil).TokenSource(withHTTPClient(ctx, o.httpClient), &expired).Token()
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
	return oauthConfig(appKey, o.appSecret, o.domain, "", nil).TokenSource(withHTTPClient(ctx, o.httpClient), token)
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

func newPKCEOptions(opts []Option) (options, error) {
	o := options{
		usePKCE:         true,
		tokenAccessType: TokenAccessTypeOffline,
	}
	applyOptions(&o, opts)
	if o.verifier == "" {
		verifier, err := generateVerifier()
		if err != nil {
			return options{}, err
		}
		o.verifier = verifier
	}
	if err := validateOAuth2Options(o); err != nil {
		return options{}, err
	}
	return o, nil
}

func newOAuth2Options(opts []Option) (options, error) {
	o := options{}
	applyOptions(&o, opts)
	o.appSecret = strings.TrimSpace(o.appSecret)
	if o.verifier != "" {
		o.usePKCE = true
	}
	if o.appSecret != "" && o.usePKCE {
		return options{}, errors.New("dropbox oauth: WithAppSecret cannot be combined with WithPKCE or WithVerifier")
	}
	if o.usePKCE && o.verifier == "" {
		verifier, err := generateVerifier()
		if err != nil {
			return options{}, err
		}
		o.verifier = verifier
	}
	if o.appSecret == "" && !o.usePKCE {
		return options{}, errors.New("dropbox oauth: either WithAppSecret or WithPKCE is required")
	}
	if err := validateOAuth2Options(o); err != nil {
		return options{}, err
	}
	return o, nil
}

func flowAppSecret(o options) string {
	if o.usePKCE {
		return ""
	}
	return o.appSecret
}

func validateOAuth2Options(o options) error {
	if o.tokenAccessType != "" && !validTokenAccessType(o.tokenAccessType) {
		return errors.New("dropbox oauth: token access type must be offline, online, or legacy")
	}
	if o.scopes != nil && len(o.scopes) == 0 {
		return errors.New("dropbox oauth: scopes must not be empty")
	}
	if o.includeGrantedScopes != "" {
		if !validIncludeGrantedScopes(o.includeGrantedScopes) {
			return errors.New("dropbox oauth: include granted scopes must be user or team")
		}
		if o.scopes == nil {
			return errors.New("dropbox oauth: scopes are required when include granted scopes is set")
		}
	}
	return nil
}

func validTokenAccessType(value TokenAccessType) bool {
	switch value {
	case TokenAccessTypeOffline, TokenAccessTypeOnline, TokenAccessTypeLegacy:
		return true
	default:
		return false
	}
}

func validIncludeGrantedScopes(value IncludeGrantedScopes) bool {
	switch value {
	case IncludeGrantedScopesUser, IncludeGrantedScopesTeam:
		return true
	default:
		return false
	}
}

func authCodeOptions(o options) []oauth2.AuthCodeOption {
	options := []oauth2.AuthCodeOption{}
	if o.usePKCE {
		options = append(options, oauth2.S256ChallengeOption(o.verifier))
	}
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
	options := []oauth2.AuthCodeOption{}
	if o.usePKCE {
		options = append(options, oauth2.VerifierOption(o.verifier))
	}
	if o.locale != "" {
		options = append(options, oauth2.SetAuthURLParam(localeParam, o.locale))
	}
	return options
}

func oauthConfig(appKey string, appSecret string, domain string, redirectURL string, scopes []string) *oauth2.Config {
	endpoint := dropbox.OAuthEndpoint(domain)
	endpoint.AuthStyle = oauth2.AuthStyleInParams
	return &oauth2.Config{
		ClientID:     appKey,
		ClientSecret: appSecret,
		Endpoint:     endpoint,
		RedirectURL:  redirectURL,
		Scopes:       append([]string(nil), scopes...),
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

func generateVerifier() (string, error) {
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
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
	accountID := extraString(token, "account_id")
	teamID := extraString(token, "team_id")
	if accountID == "" {
		accountID = teamID
	}
	return &FlowResult{
		Token:     token,
		AccountID: accountID,
		TeamID:    teamID,
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
