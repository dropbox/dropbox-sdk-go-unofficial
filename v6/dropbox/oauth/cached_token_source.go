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

package oauth

import (
	"context"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

// AccessTokenSnapshot is a copy of a cached access token and its expiration.
// It contains no refresh credentials. Access tokens are sensitive and must not
// be logged or disclosed.
type AccessTokenSnapshot struct {
	// AccessToken is the cached OAuth2 access token.
	AccessToken string
	// Expiry is the token's expiration time, or zero if unknown. A zero expiry
	// does not guarantee an unlimited lifetime; OAuth2 reuses such tokens
	// without automatically refreshing them.
	Expiry time.Time
}

// CachedTokenSource is a refreshable OAuth2 token source with local snapshots.
// It is safe for concurrent use and may be shared by multiple Dropbox clients.
// Create one with NewCachedTokenSource; it must not be copied after first use.
type CachedTokenSource struct {
	source oauth2.TokenSource
	// Serialize Token calls through publication so an older result cannot
	// overwrite a newer snapshot. Snapshot never acquires this mutex.
	tokenMu sync.Mutex
	mu      sync.RWMutex
	cached  AccessTokenSnapshot
}

// NewCachedTokenSource returns a Dropbox token source using the same options
// and refresh behavior as TokenSource. It copies token and immediately caches
// its access token and expiry, even if expired. A nil token or an empty access
// token leaves the snapshot empty until Token successfully obtains a token.
// Construction does not make HTTP requests.
func NewCachedTokenSource(ctx context.Context, appKey string, token *oauth2.Token, opts ...Option) *CachedTokenSource {
	if token != nil {
		copied := *token
		token = &copied
	}
	s := &CachedTokenSource{source: TokenSource(ctx, appKey, token, opts...)}
	if token != nil && token.AccessToken != "" {
		s.cached = AccessTokenSnapshot{AccessToken: token.AccessToken, Expiry: token.Expiry}
	}
	return s
}

// Token returns the current token, refreshing it through OAuth2 if necessary.
// Successful calls update the snapshot; failures preserve the previous snapshot.
// As with oauth2.TokenSource, callers must not modify the returned token.
func (s *CachedTokenSource) Token() (*oauth2.Token, error) {
	s.tokenMu.Lock()
	defer s.tokenMu.Unlock()
	token, err := s.source.Token()
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.cached = AccessTokenSnapshot{AccessToken: token.AccessToken, Expiry: token.Expiry}
	s.mu.Unlock()
	// oauth2.NewClient adds its own ReuseTokenSource, which writes expiry
	// metadata into tokens it receives. Give each caller a separate token so
	// clients sharing this source cannot mutate the underlying cache's token.
	copied := *token
	return &copied, nil
}

// Snapshot returns a copy of the last successfully obtained access token and
// expiry, or the initial token, and reports whether an access token is cached.
// It never makes HTTP requests, refreshes, or validates credentials. Expired
// tokens are included, and snapshots remain unchanged after later refreshes.
// Snapshot does not wait for an in-flight refresh to finish.
func (s *CachedTokenSource) Snapshot() (AccessTokenSnapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cached, s.cached.AccessToken != ""
}
