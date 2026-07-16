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

package check

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
	// App : This endpoint performs App Authentication, validating the supplied
	// app key and secret, and returns the supplied string, to allow you to test
	// your code and connection to the Dropbox API. It has no other effect. If
	// you receive an HTTP 200 response with the supplied query, it indicates at
	// least part of the Dropbox API infrastructure is working and that the app
	// key and secret valid.
	App(arg *EchoArg) (res *EchoResult, err error)
	// User : This endpoint performs User Authentication, validating the
	// supplied access token, and returns the supplied string, to allow you to
	// test your code and connection to the Dropbox API. It has no other effect.
	// If you receive an HTTP 200 response with the supplied query, it indicates
	// at least part of the Dropbox API infrastructure is working and that the
	// access token is valid.
	User(arg *EchoArg) (res *EchoResult, err error)
}

// ContextClient interface describes all routes in this namespace with context support
type ContextClient interface {
	Client
	// AppContext : This endpoint performs App Authentication, validating the
	// supplied app key and secret, and returns the supplied string, to allow
	// you to test your code and connection to the Dropbox API. It has no other
	// effect. If you receive an HTTP 200 response with the supplied query, it
	// indicates at least part of the Dropbox API infrastructure is working and
	// that the app key and secret valid.
	AppContext(ctx context.Context, arg *EchoArg) (res *EchoResult, err error)
	// UserContext : This endpoint performs User Authentication, validating the
	// supplied access token, and returns the supplied string, to allow you to
	// test your code and connection to the Dropbox API. It has no other effect.
	// If you receive an HTTP 200 response with the supplied query, it indicates
	// at least part of the Dropbox API infrastructure is working and that the
	// access token is valid.
	UserContext(ctx context.Context, arg *EchoArg) (res *EchoResult, err error)
}

type apiImpl dropbox.Context

// AppAPIError is an error-wrapper for the app route
type AppAPIError struct {
	dropbox.APIError
	EndpointError *EchoError `json:"error"`
}

// AppContext : This endpoint performs App Authentication, validating the
// supplied app key and secret, and returns the supplied string, to allow you to
// test your code and connection to the Dropbox API. It has no other effect. If
// you receive an HTTP 200 response with the supplied query, it indicates at
// least part of the Dropbox API infrastructure is working and that the app key
// and secret valid.
func (dbx *apiImpl) AppContext(ctx context.Context, arg *EchoArg) (res *EchoResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "check",
		Route:        "app",
		Auth:         "app",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr AppAPIError
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

func (dbx *apiImpl) App(arg *EchoArg) (res *EchoResult, err error) {
	return dbx.AppContext(context.Background(), arg)
}

// UserAPIError is an error-wrapper for the user route
type UserAPIError struct {
	dropbox.APIError
	EndpointError *EchoError `json:"error"`
}

// UserContext : This endpoint performs User Authentication, validating the
// supplied access token, and returns the supplied string, to allow you to test
// your code and connection to the Dropbox API. It has no other effect. If you
// receive an HTTP 200 response with the supplied query, it indicates at least
// part of the Dropbox API infrastructure is working and that the access token
// is valid.
func (dbx *apiImpl) UserContext(ctx context.Context, arg *EchoArg) (res *EchoResult, err error) {
	req := dropbox.Request{
		Host:         "api",
		Namespace:    "check",
		Route:        "user",
		Auth:         "user",
		Style:        "rpc",
		Arg:          arg,
		ExtraHeaders: nil,
	}

	var resp []byte
	var respBody io.ReadCloser
	resp, respBody, err = (*dropbox.Context)(dbx).ExecuteContext(ctx, req, nil)
	if err != nil {
		var appErr UserAPIError
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

func (dbx *apiImpl) User(arg *EchoArg) (res *EchoResult, err error) {
	return dbx.UserContext(context.Background(), arg)
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
