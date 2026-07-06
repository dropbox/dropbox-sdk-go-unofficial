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
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	dropboxoauth "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/oauth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/users"
	"golang.org/x/oauth2"
)

func ExampleNewPKCEFlow() {
	ctx := context.Background()

	flow, err := dropboxoauth.NewPKCEFlow(
		"APP_KEY",
		dropboxoauth.WithScopes("files.metadata.read"),
	)
	if err != nil {
		panic(err)
	}

	authURL := flow.AuthCodeURL()
	_ = authURL // Redirect the user to authURL and read the authorization code.

	code := "AUTHORIZATION_CODE"
	_ = code
	_ = ctx
	// token, err := flow.Exchange(ctx, code)
}

func ExampleNewWebPKCEFlow() {
	verifier, err := dropboxoauth.GenerateVerifier()
	if err != nil {
		panic(err)
	}

	flow, err := dropboxoauth.NewWebPKCEFlow(
		"APP_KEY",
		"https://example.com/dropbox/callback",
		dropboxoauth.WithVerifier(verifier),
		dropboxoauth.WithScopes("files.metadata.read"),
	)
	if err != nil {
		panic(err)
	}

	authURL, csrfToken, err := flow.Start("optional-url-state")
	if err != nil {
		panic(err)
	}

	_ = authURL   // Redirect the user to authURL.
	_ = csrfToken // Store csrfToken and verifier in the user's web session.

	// In the callback handler, rebuild the flow with the stored verifier.
	flow, err = dropboxoauth.NewWebPKCEFlow(
		"APP_KEY",
		"https://example.com/dropbox/callback",
		dropboxoauth.WithVerifier(verifier),
	)
	if err != nil {
		panic(err)
	}
	_ = flow
	// result, err := flow.Finish(r.Context(), r.URL.Query(), csrfToken)
}

func ExampleNewWebPKCEFlow_openID() {
	verifier, err := dropboxoauth.GenerateVerifier()
	if err != nil {
		panic(err)
	}

	flow, err := dropboxoauth.NewWebPKCEFlow(
		"APP_KEY",
		"https://example.com/dropbox/callback",
		dropboxoauth.WithVerifier(verifier),
		dropboxoauth.WithScopes(
			dropboxoauth.ScopeOpenID,
			dropboxoauth.ScopeEmail,
			dropboxoauth.ScopeProfile,
		),
	)
	if err != nil {
		panic(err)
	}

	authURL, csrfToken, err := flow.Start("optional-url-state")
	if err != nil {
		panic(err)
	}
	_ = authURL
	_ = csrfToken
	// After Finish succeeds, call UserinfoContext on an openid client.
}

func ExampleNewOAuth2FlowNoRedirect() {
	flow, err := dropboxoauth.NewOAuth2FlowNoRedirect(
		"APP_KEY",
		dropboxoauth.WithAppSecret("APP_SECRET"),
		dropboxoauth.WithScopes("files.metadata.read"),
		dropboxoauth.WithTokenAccessType(dropboxoauth.TokenAccessTypeOffline),
	)
	if err != nil {
		panic(err)
	}

	authURL := flow.Start()
	_ = authURL // Send the user to authURL and read the authorization code.

	code := "AUTHORIZATION_CODE"
	_ = code
	// result, err := flow.Finish(context.Background(), code)
}

func ExampleTokenSource() {
	ctx := context.Background()
	token := &oauth2.Token{
		AccessToken:  "ACCESS_TOKEN",
		RefreshToken: "REFRESH_TOKEN",
		Expiry:       time.Now().Add(time.Hour),
	}

	config := dropbox.Config{
		TokenSource: dropboxoauth.TokenSource(ctx, "APP_KEY", token),
	}

	dbx := users.New(config)
	_ = dbx
	// account, err := dbx.GetCurrentAccount()
}
