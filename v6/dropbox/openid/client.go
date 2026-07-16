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

package openid

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/auth"
)

// Client interface describes all routes in this namespace
type Client interface {
	// Userinfo : This route is used for refreshing the info that is found in
	// the id_token during the OIDC flow. This route doesn't require any
	// arguments and will use the scopes approved for the given access token.
	Userinfo(arg *UserInfoArgs) (res *UserInfoResult, err error)
}

// ContextClient interface describes all routes in this namespace with context support
type ContextClient interface {
	Client
	// UserinfoContext : This route is used for refreshing the info that is
	// found in the id_token during the OIDC flow. This route doesn't require
	// any arguments and will use the scopes approved for the given access
	// token.
	UserinfoContext(ctx context.Context, arg *UserInfoArgs) (res *UserInfoResult, err error)
}

type apiImpl dropbox.Context

// UserinfoAPIError is an error-wrapper for the userinfo route
type UserinfoAPIError struct {
	dropbox.APIError
	EndpointError *UserInfoError `json:"error"`
}

// UserinfoContext : This route is used for refreshing the info that is found in
// the id_token during the OIDC flow. This route doesn't require any arguments
// and will use the scopes approved for the given access token.
func (dbx *apiImpl) UserinfoContext(ctx context.Context, arg *UserInfoArgs) (res *UserInfoResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "openid",
		Route:        "userinfo",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr UserinfoAPIError
		err = auth.ParseError(err, &appErr)
		if errors.Is(err, &appErr) {
			err = appErr
		}
		return
	}

	err = json.Unmarshal(resp, &res)
	if err != nil {
		return
	}

	_ = respBody
	return
}

func (dbx *apiImpl) Userinfo(arg *UserInfoArgs) (res *UserInfoResult, err error) {
	return dbx.UserinfoContext(context.Background(), arg)
}

// NewContext returns a ContextClient implementation for this namespace
func NewContext(c dropbox.Config) ContextClient {
	ctx := apiImpl(dropbox.NewContext(c))
	return &ctx
}

// New returns a Client implementation for this namespace
func New(c dropbox.Config) Client {
	return NewContext(c)
}
