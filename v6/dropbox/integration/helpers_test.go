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

//go:build integration

// Package integration contains tests that talk to the live Dropbox API using
// real credentials. They are excluded from the default build and only compile
// under the "integration" build tag:
//
//	go test -tags integration ./dropbox/integration/...
//
// Credentials are read from the environment, mirroring the CI setup used by the
// Python SDK (SCOPED_USER_* and SCOPED_TEAM_* refresh-token credentials). For
// local runs they can instead be loaded from a dropbox-sdk-python-secrets.json
// file pointed to by DROPBOX_SDK_SECRETS_FILE.
package integration

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/oauth"
	"golang.org/x/oauth2"
)

// Credential environment variable names. These match the names used by the
// dropbox-sdk-python integration tests and CI.
const (
	envUserClientID     = "SCOPED_USER_CLIENT_ID"
	envUserClientSecret = "SCOPED_USER_CLIENT_SECRET"
	envUserRefreshToken = "SCOPED_USER_REFRESH_TOKEN"

	envTeamClientID     = "SCOPED_TEAM_CLIENT_ID"
	envTeamClientSecret = "SCOPED_TEAM_CLIENT_SECRET"
	envTeamRefreshToken = "SCOPED_TEAM_REFRESH_TOKEN"

	// envSecretsFile optionally points to a JSON file holding the credentials
	// above (the dropbox-sdk-python-secrets.json format). Values already set in
	// the environment take precedence.
	envSecretsFile = "DROPBOX_SDK_SECRETS_FILE"
)

// loadSecretsFile loads any credentials from the JSON file referenced by
// DROPBOX_SDK_SECRETS_FILE into the environment, without overwriting values
// that are already set. It is a no-op when the variable is unset.
func loadSecretsFile(t *testing.T) {
	t.Helper()
	path := os.Getenv(envSecretsFile)
	if path == "" {
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s=%q: %v", envSecretsFile, path, err)
	}
	var secrets map[string]string
	if err := json.Unmarshal(data, &secrets); err != nil {
		t.Fatalf("parsing secrets file %q: %v", path, err)
	}
	for k, v := range secrets {
		if _, ok := os.LookupEnv(k); !ok {
			t.Setenv(k, v)
		}
	}
}

// requireEnv returns the value of the named environment variable, skipping the
// test if it is not set. Integration credentials are not always available (for
// example on fork PRs), so a missing credential skips rather than fails.
func requireEnv(t *testing.T, name string) string {
	t.Helper()
	v := os.Getenv(name)
	if v == "" {
		t.Skipf("%s is not set; skipping integration test", name)
	}
	return v
}

// refreshTokenSource builds an OAuth token source that mints short-lived access
// tokens from a long-lived refresh token, matching how the Python SDK
// authenticates its integration fixtures.
func refreshTokenSource(t *testing.T, appKey, appSecret, refreshToken string) oauth2.TokenSource {
	t.Helper()
	return oauth.TokenSource(
		context.Background(),
		appKey,
		&oauth2.Token{RefreshToken: refreshToken},
		oauth.WithAppSecret(appSecret),
	)
}

// userConfig returns a dropbox.Config authenticated as the scoped user via a
// refresh token, skipping the test if the credentials are unavailable.
func userConfig(t *testing.T) dropbox.Config {
	t.Helper()
	loadSecretsFile(t)
	appKey := requireEnv(t, envUserClientID)
	appSecret := requireEnv(t, envUserClientSecret)
	refreshToken := requireEnv(t, envUserRefreshToken)
	return dropbox.Config{
		TokenSource: refreshTokenSource(t, appKey, appSecret, refreshToken),
	}
}

// teamConfig returns a dropbox.Config authenticated as the scoped team via a
// refresh token, skipping the test if the credentials are unavailable.
func teamConfig(t *testing.T) dropbox.Config {
	t.Helper()
	loadSecretsFile(t)
	appKey := requireEnv(t, envTeamClientID)
	appSecret := requireEnv(t, envTeamClientSecret)
	refreshToken := requireEnv(t, envTeamRefreshToken)
	return dropbox.Config{
		TokenSource: refreshTokenSource(t, appKey, appSecret, refreshToken),
	}
}
